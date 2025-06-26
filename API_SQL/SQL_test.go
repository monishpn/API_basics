package SQL

import (
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSQL(t *testing.T) {
	db := &input{}
	var test = []byte(`{ "task": "Testing"}`)

	//Opening DataBase
	var err error
	db.data, err = sql.Open("mysql", "root:root123@tcp(localhost:3306)/test_db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.data.Close()
	err = db.data.Ping()
	if err != nil {
		log.Fatal(err)
	}

	//Storing the count of the
	var count int
	err = db.data.QueryRow("Select count(*) from TASKS").Scan(&count)
	if err != nil {
		t.Errorf("Error while checking the count: %v", err)
	}

	//POST request evaluation
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/task", test)

	//GET request evaluation

}
