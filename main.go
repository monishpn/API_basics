package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

func hellohandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Hello, World!"))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}
}

type Record struct {
	ID        int
	Task      string
	Completed bool
}

type input struct {
	idx int
	Rec []byte
}

type slices struct {
	slice []input
}

func idGen() func() int {
	id := 0

	return func() int {
		id++
		return id
	}
}

func (re *slices) addTask(w http.ResponseWriter, r *http.Request, i int) {
	defer r.Body.Close()

	msg, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	rec := Record{i, string(msg), false}

	jsonRec, err := json.Marshal(rec)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusCreated)

	re.slice = append(re.slice, input{rec.ID, jsonRec})
}

func (re *slices) getByID(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	index, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("%s", err.Error())
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)

		return
	}

	for _, item := range re.slice {
		if item.idx == index {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write(item.Rec)

			if err != nil {
				log.Printf("%s", err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

				return
			}

			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func (re *slices) viewTask(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)

	for _, task := range re.slice {
		_, err := w.Write(task.Rec)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("%s", err.Error())

			return
		}
	}

	log.Printf("\n")
}

func (re *slices) completeTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	index, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("%s", err.Error())
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)

		return
	}

	var jsonSlice Record

	for i, item := range re.slice {
		if item.idx != index {
			continue
		}

		err := json.Unmarshal(item.Rec, &jsonSlice)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("%s", err.Error())

			return
		}

		jsonSlice.Completed = true

		updatedJSON, err := json.Marshal(jsonSlice)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("%s", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusAccepted)

		re.slice[i].Rec = updatedJSON

		return
	}

	w.WriteHeader(http.StatusNotFound)
}

func (re *slices) deleteTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	index, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("%s", err.Error())
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)

		return
	}

	for i, item := range re.slice {
		if item.idx == index {
			w.WriteHeader(http.StatusOK)

			re.slice = append(re.slice[:i], re.slice[i+1:]...)

			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func main() {
	data := &slices{}
	getID := idGen()

	http.HandleFunc("/", hellohandler)

	http.HandleFunc("POST /task", func(w http.ResponseWriter, r *http.Request) {
		i := getID()
		data.addTask(w, r, i)
	})
	http.HandleFunc("GET /task/{id}", data.getByID)
	http.HandleFunc("GET /task", data.viewTask)
	http.HandleFunc("PUT /task/{id}", data.completeTask)
	http.HandleFunc("DELETE /task/{id}", data.deleteTask)

	srv := http.Server{
		Addr:         ":8080",
		Handler:      nil, // same as default mux
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
