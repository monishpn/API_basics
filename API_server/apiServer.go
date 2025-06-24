package API_server

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
	Rec Record
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

	var reqBody struct {
		T string `json:"task"`
	}

	err = json.Unmarshal(msg, &reqBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	rec := Record{i, reqBody.T, false}

	w.WriteHeader(http.StatusCreated)

	re.slice = append(re.slice, input{rec.ID, rec})
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
		if item.idx != index {
			continue
		}

		w.WriteHeader(http.StatusOK)

		msg, _ := json.Marshal(item.Rec)
		_, err := w.Write(msg)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

			return
		}

		return
	}

	w.WriteHeader(http.StatusNotFound)
}

func (re *slices) viewTask(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)

	for _, task := range re.slice {
		msg, err := json.Marshal(task.Rec)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

			return
		}

		_, err = w.Write(msg)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("%s", err.Error())

			return
		}
	}
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

	for i, item := range re.slice {
		if item.idx != index {
			continue
		}

		w.WriteHeader(http.StatusAccepted)

		re.slice[i].Rec.Completed = true

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
