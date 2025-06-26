package SQL

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

type input struct {
	data *sql.DB
}

func (db *input) addTask(w http.ResponseWriter, r *http.Request) {
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

	_, _ = db.data.Exec("INSERT INTO TASKS (task,completed) VALUES (?,?);", reqBody.T, false)

}

func (db *input) getByID(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	index, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("%s", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

		return
	}

	ans := db.data.QueryRow("SELECT * FROM TASKS WHERE id=?", index)

	var id int

	var task string

	var completed bool

	err = ans.Scan(&id, &task, &completed)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Task not found")

			return
		}

		http.Error(w, "Failed to read task row", http.StatusInternalServerError)

		log.Printf("While getting by ID -> %v", err)

		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "ID: %d, Task: %s, Completed: %t", id, task, completed)
}

func (db *input) viewTask(w http.ResponseWriter, _ *http.Request) {
	rows, err := db.data.Query("SELECT * FROM TASKS")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, "Failed to query tasks", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var id int

		var task string

		var completed bool

		err := rows.Scan(&id, &task, &completed)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			http.Error(w, "Failed to read task row", http.StatusInternalServerError)
			log.Printf("%s", err.Error())

			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "ID: %d, Task: %s, Completed: %t\n", id, task, completed)
	}
}

func (db *input) completeTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	index, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("%s", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

		return
	}

	res, err := db.data.Exec("UPDATE TASKS SET completed= true WHERE id=?", index)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%s", err.Error())

		return
	}

	check, _ := res.RowsAffected()

	if check == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Task not found")

		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Updated Successfully")
}

func (db *input) deleteTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	index, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("%s", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

		return
	}

	del, err := db.data.Exec("DELETE FROM TASKS WHERE id=?", index)

	if err != nil {
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	check, _ := del.RowsAffected()
	if check == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Task not found")

		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Deleted Successfully")
}

func main() {
	db := &input{}

	var err error

	db.data, err = sql.Open("mysql", "root:root123@tcp(localhost:3306)/test_db")

	if err != nil {
		log.Fatal("Error while running the server -> ", err)
	}

	err = db.data.Ping()
	if err != nil {
		log.Fatal("Error while Checking for start of the server -> ", err)
	}

	_, err = db.data.Exec("CREATE TABLE IF NOT EXISTS TASKS ( id int auto_increment primary key, task text, completed bool );")
	if err != nil {
		log.Fatal("Error while creating database\n", err)
	}

	http.HandleFunc("/", hellohandler)

	http.HandleFunc("POST /task", db.addTask)
	http.HandleFunc("GET /task/{id}", db.getByID)
	http.HandleFunc("GET /task", db.viewTask)
	http.HandleFunc("PUT /task/{id}", db.completeTask)
	http.HandleFunc("DELETE /task/{id}", db.deleteTask)

	srv := http.Server{
		Addr:         ":8080",
		Handler:      nil, // same as default mux
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
