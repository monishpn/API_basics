package SQL

import (
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestHelloWorld(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	hellohandler(w, req)

	if w.Body.String() != "Hello, World!" {
		t.Error("Hello world failed")
	}

}

func TestSQL(t *testing.T) {
	db := &input{}
	test := `{ "task": "Testing"}`

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
	err = db.data.QueryRow("Select count(*) from TASKS;").Scan(&count)
	if err != nil {
		t.Errorf("Error while checking the count: %v", err)
	}

	//POST request evaluation
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/task", strings.NewReader(test))
	response := httptest.NewRecorder()

	db.addTask(response, request)

	var checkPost int
	err = db.data.QueryRow("Select count(*) from TASKS").Scan(&checkPost)
	if err != nil {
		t.Errorf("Error while checking the count after POSTing: %v", err)
	}
	if checkPost != count+1 {
		t.Errorf("Wrong count after POSTing: expected %d, got %d", count+1, checkPost)
	}

	//GET BY ID request evaluation
	var lastID int
	err = db.data.QueryRow("Select id from TASKS order by id DESC Limit 1;").Scan(&lastID)
	if err != nil {
		t.Errorf("Error while checking the Last Index: %v", err)
	}

	request = httptest.NewRequest(http.MethodGet, "http://localhost:8080/task/{id}", http.NoBody)
	response = httptest.NewRecorder()
	request.SetPathValue("id", strconv.Itoa(lastID))

	db.getByID(response, request)

	exp := "ID: " + strconv.Itoa(lastID) + ", Task: Testing, Completed: false"

	if response.Body.String() != exp {
		t.Errorf("Wrong result:\n expected %s,\n got %s", exp, response.Body.String())
	}

	//PUT COMPLETED request evaluation
	request = httptest.NewRequest(http.MethodPut, "http://localhost:8080/task/{id}", http.NoBody)
	response = httptest.NewRecorder()
	request.SetPathValue("id", strconv.Itoa(lastID))

	db.completeTask(response, request)

	var completedCheck bool

	err = db.data.QueryRow("select completed from TASKS where id=?", lastID).Scan(&completedCheck)
	if err != nil {
		t.Errorf("Error while checking the Copleted Status: %v", err)
	}
	if completedCheck != true {
		t.Errorf("Wrong result: expected %v, got %v", true, completedCheck)
	}

	//DELETE request evaluation
	request = httptest.NewRequest(http.MethodDelete, "http://localhost:8080/task/{id}", http.NoBody)
	response = httptest.NewRecorder()
	request.SetPathValue("id", strconv.Itoa(lastID))

	db.deleteTask(response, request)

	var lastID_afterDelete int
	err = db.data.QueryRow("Select id from TASKS order by id DESC Limit 1;").Scan(&lastID_afterDelete)
	if err != nil {
		t.Errorf("Error while checking the Last Index: %v", err)
	}

	if lastID_afterDelete == lastID {
		t.Errorf("Wrong result: expected %d, got %d", lastID_afterDelete, lastID)
	}

}


func TestSQLWithError(t *testing.T) {
	db := &input{}
	test := `{"Testing"}`

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

	//
	//
	//
	//ADD TASK when the JSON is corrupted
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/task", strings.NewReader(test))
	response := httptest.NewRecorder()

	db.addTask(response, request)

	if response.Header().


}
