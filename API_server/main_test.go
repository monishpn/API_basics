package API_server

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

// CREATING A MIRROR COPY OF RESPONSE_WRITER.
type errorWriter struct {
	status int
}

func (e *errorWriter) Header() http.Header {
	return http.Header{}
}

func (e *errorWriter) Write([]byte) (int, error) {
	return 0, errors.New("Forced error")
}

func (e *errorWriter) WriteHeader(statusCode int) {
	e.status = statusCode
}

func TestHelloWorld(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "localhost:8080/", http.NoBody)

	response := httptest.NewRecorder()

	hellohandler(response, request)

	exp := "Hello, World!"

	op, _ := io.ReadAll(response.Body)

	if string(op) != exp {
		t.Errorf("Expected %s, got %s", exp, string(op))
	}

	if response.Code != http.StatusOK {
		t.Errorf("Expected %d, got %d", http.StatusOK, response.Code)
	}

	// W.Write err test
	errResponse := &errorWriter{}
	hellohandler(errResponse, request)

	if errResponse.status != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, errResponse.status)
	}
}

func Test_Passing(t *testing.T) {
	re := slices{}
	ID := "1"
	task := `{
		"task":"Eating"
	}`
	work := "Eating"

	//
	//
	//
	// For AddTask
	request := httptest.NewRequest(http.MethodPost, "localhost:8080/task", strings.NewReader(task))

	response := httptest.NewRecorder()

	id, _ := strconv.Atoi(ID)
	re.addTask(response, request, id)

	if response.Code != http.StatusCreated {
		t.Errorf("Expected %d, got %d", http.StatusCreated, response.Code)
	}

	//
	//
	//
	// For GetByID
	request = httptest.NewRequest(http.MethodGet, "localhost:8080/task/{id}", http.NoBody)
	request.SetPathValue("id", ID)

	response = httptest.NewRecorder()

	re.getByID(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("Expected %d, got %d", http.StatusOK, response.Code)
	}

	op, _ := io.ReadAll(response.Body)
	exp, _ := json.Marshal(Record{1, work, false})

	if !bytes.Equal(op, exp) {
		t.Errorf("Expected %s, got %s", exp, string(op))
	}

	//
	//
	//
	// For Completed Task
	request = httptest.NewRequest(http.MethodPut, "localhost:8080/task/{id}", http.NoBody)
	request.SetPathValue("id", ID)

	response = httptest.NewRecorder()

	re.completeTask(response, request)

	if response.Code != http.StatusAccepted {
		t.Errorf("Expected %d, got %d", http.StatusAccepted, response.Code)
	}

	//
	//
	//
	// For Viewing
	request = httptest.NewRequest(http.MethodGet, "localhost:8080/task/", http.NoBody)

	response = httptest.NewRecorder()

	re.viewTask(response, request)

	exp, _ = json.Marshal(Record{1, work, true})
	op, _ = io.ReadAll(response.Body)

	if !bytes.Equal(op, exp) {
		t.Errorf("Expected %s, got %s", exp, string(op))
	}

	//
	//
	//
	// For Deleting
	request = httptest.NewRequest(http.MethodDelete, "localhost:8080/task/{id}", http.NoBody)
	request.SetPathValue("id", ID)

	response = httptest.NewRecorder()

	re.deleteTask(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("Expected %d, got %d", http.StatusOK, response.Code)
	}
}

func Test_Failing(t *testing.T) {
	re := slices{}
	task := `{
		"task":"Eating"
	}`
	work := "Eating"

	//
	//
	//
	// ADD TASK error check - when the input is not JSON
	request := httptest.NewRequest(http.MethodPost, "localhost:8080/task", strings.NewReader(work))
	response := httptest.NewRecorder()

	re.addTask(response, request, 1)

	if response.Code != http.StatusInternalServerError {
		t.Errorf("Expected %d, got %d", http.StatusInternalServerError, response.Code)
	}

	//
	//
	//
	// Add a sample data into re to check for completed and delete
	request = httptest.NewRequest(http.MethodPut, "/task", bytes.NewBufferString(task))
	re.addTask(response, request, 1)

	//
	//
	//
	// GetByID PRINT ONLY REQUESTED ID
	request = httptest.NewRequest(http.MethodGet, "/task/{id}", http.NoBody)
	response = httptest.NewRecorder()

	// GetByID error check - when the ID given is not an integer
	request.SetPathValue("id", "r")

	re.getByID(response, request)

	if response.Code != http.StatusNotFound {
		t.Errorf("Expected %d, got %d", http.StatusNotFound, response.Code)
	}

	// GetByID error check - when the ID not present
	request.SetPathValue("id", "3")

	response = httptest.NewRecorder()

	re.getByID(response, request)

	if response.Code != http.StatusNotFound {
		t.Errorf("Expected %d, got %d", http.StatusNotFound, response.Code)
	}

	// GetByID error check - w.Write err test
	request.SetPathValue("id", "1")

	errResponse := &errorWriter{}

	re.getByID(errResponse, request)

	if errResponse.status != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, errResponse.status)
	}

	//
	//
	//
	// VIEWING TASK
	// viewTask error check - w.Write err test
	errResponse = &errorWriter{}
	re.viewTask(errResponse, request)

	if errResponse.status != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, errResponse.status)
	}

	//
	//
	//
	// COMPLETE TASK - MARKING TRUE
	request = httptest.NewRequest(http.MethodPut, "/task/{id}", http.NoBody)
	response = httptest.NewRecorder()

	// CompleteTask error check - when the ID given is not a integer
	request.SetPathValue("id", "r")

	re.completeTask(response, request)

	if response.Code != http.StatusNotFound {
		t.Errorf("Expected %d, got %d", http.StatusNotFound, response.Code)
	}

	// CompleteTask error check - when the ID not present
	request.SetPathValue("id", "3")

	response = httptest.NewRecorder()

	re.completeTask(response, request)

	if response.Code != http.StatusNotFound {
		t.Errorf("Expected %d, got %d", http.StatusNotFound, response.Code)
	}

	//
	//
	//
	// DELETE TASK - REMOVING
	request = httptest.NewRequest(http.MethodDelete, "/task/{id}", http.NoBody)
	response = httptest.NewRecorder()

	// DeleteTask  error check - when the ID given is not a integer
	request.SetPathValue("id", "r")

	re.deleteTask(response, request)

	if response.Code != http.StatusNotFound {
		t.Errorf("Expected %d, got %d", http.StatusNotFound, response.Code)
	}

	// DeleteTask error check - when the ID not present
	request.SetPathValue("id", "3")

	response = httptest.NewRecorder()

	re.deleteTask(response, request)

	if response.Code != http.StatusNotFound {
		t.Errorf("Expected %d, got %d", http.StatusNotFound, response.Code)
	}
}
