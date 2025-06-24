package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
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
	data *sql.DB
}

func idGen() func() int {
	id := 0

	return func() int {
		id++
		return id
	}
}

func (db *input) addTask(w http.ResponseWriter, r *http.Request, i int) {
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

	w.WriteHeader(http.StatusCreated)
	_, err = db.data.Exec("INSERT INTO tasks (id,task,completed) VALUES (?,?,?)", i, reqBody.T, false)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%s", err.Error())
		return
	}

}

func (db *input) getByID(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	index, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("%s", err.Error())
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)

		return
	}

	ans := db.data.QueryRow("SELECT * FROM tasks WHERE id=?", index)

	var op string
	err = ans.Scan(op)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte(op))

	w.WriteHeader(http.StatusNotFound)
}

func (db *input) viewTask(w http.ResponseWriter, _ *http.Request) {

	rows, err := db.data.Query("SELECT task FROM tasks")
	if err != nil {
		http.Error(w, "Failed to query tasks", http.StatusInternalServerError)
		return
	}
	//defer rows.Close()

	for rows.Next() {
		var id int
		var task string
		var completed bool

		err := rows.Scan(&id, &task, &completed)
		if err != nil {
			http.Error(w, "Failed to read task row", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "ID: %d, Task: %s, Completed: %t\n", id, task, completed)
	}

}

//
//func (db *input) completeTask(w http.ResponseWriter, r *http.Request) {
//	defer r.Body.Close()
//
//	index, err := strconv.Atoi(r.PathValue("id"))
//	if err != nil {
//		w.WriteHeader(http.StatusNotFound)
//		log.Printf("%s", err.Error())
//		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
//
//		return
//	}
//
//	for i, item := range re.slice {
//		if item.idx != index {
//			continue
//		}
//
//		w.WriteHeader(http.StatusAccepted)
//
//		re.slice[i].Rec.Completed = true
//
//		return
//	}
//
//	w.WriteHeader(http.StatusNotFound)
//}
//
//func (db *input) deleteTask(w http.ResponseWriter, r *http.Request) {
//	defer r.Body.Close()
//
//	index, err := strconv.Atoi(r.PathValue("id"))
//	if err != nil {
//		w.WriteHeader(http.StatusNotFound)
//		log.Printf("%s", err.Error())
//		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
//
//		return
//	}
//
//	for i, item := range re.slice {
//		if item.idx == index {
//			w.WriteHeader(http.StatusOK)
//
//			re.slice = append(re.slice[:i], re.slice[i+1:]...)
//
//			return
//		}
//	}
//
//	w.WriteHeader(http.StatusNotFound)
//}

func main() {
	db := &input{}
	var err error
	db.data, err = sql.Open("mysql", "root:root123@tcp(localhost:3306)/test_db")
	if err != nil {
		log.Fatal(err)
	}

	//var q string = "CREATE TABLE TASKS ( id int, task text, completed bool )"
	//_, err = db.data.Exec(q)
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}

	getID := idGen()

	http.HandleFunc("/", hellohandler)

	http.HandleFunc("POST /task", func(w http.ResponseWriter, r *http.Request) {
		i := getID()
		db.addTask(w, r, i)
	})
	http.HandleFunc("GET /task/{id}", db.getByID)
	http.HandleFunc("GET /task", db.viewTask)
	//http.HandleFunc("PUT /task/{id}", db.completeTask)
	//http.HandleFunc("DELETE /task/{id}", db.deleteTask)

	srv := http.Server{
		Addr:         ":8080",
		Handler:      nil, // same as default mux
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
