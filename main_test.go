package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func TestHelloWorld(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
	defer cancel()

	request, _ := http.NewRequestWithContext(ctx, http.MethodGet, "localhost:8080/", http.NoBody)

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
}

func Test(t *testing.T) {
	re := slices{}
	ID := "1"
	task := "Eating"

	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
	defer cancel()

	// For AddTask
	request, _ := http.NewRequestWithContext(ctx, http.MethodPost, "localhost:8080/task", bytes.NewBufferString(task))

	response := httptest.NewRecorder()

	id, _ := strconv.Atoi(ID)
	re.addTask(response, request, id)

	if response.Code != http.StatusCreated {
		t.Errorf("Expected %d, got %d", http.StatusCreated, response.Code)
	}

	// For GetByID
	request, _ = http.NewRequestWithContext(ctx, http.MethodGet, "localhost:8080/task/{id}", http.NoBody)
	request.SetPathValue("id", ID)

	response = httptest.NewRecorder()

	re.getByID(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("Expected %d, got %d", http.StatusOK, response.Code)
	}

	op, _ := io.ReadAll(response.Body)
	exp, _ := json.Marshal(Record{1, task, false})

	if !bytes.Equal(op, exp) {
		t.Errorf("Expected %s, got %s", exp, string(op))
	}

	// For Completed Task
	request, _ = http.NewRequestWithContext(ctx, http.MethodPut, "localhost:8080/task/{id}", http.NoBody)
	request.SetPathValue("id", ID)

	response = httptest.NewRecorder()

	re.completeTask(response, request)

	if response.Code != http.StatusAccepted {
		t.Errorf("Expected %d, got %d", http.StatusAccepted, response.Code)
	}

	// For Viewing
	request, _ = http.NewRequestWithContext(ctx, http.MethodGet, "localhost:8080/task/", http.NoBody)

	response = httptest.NewRecorder()

	re.viewTask(response, request)

	exp, _ = json.Marshal(Record{1, task, true})
	op, _ = io.ReadAll(response.Body)

	if !bytes.Equal(op, exp) {
		t.Errorf("Expected %s, got %s", exp, string(op))
	}

	// For Deleting
	request, _ = http.NewRequestWithContext(ctx, http.MethodDelete, "localhost:8080/task/{id}", http.NoBody)
	request.SetPathValue("id", ID)

	response = httptest.NewRecorder()

	re.deleteTask(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("Expected %d, got %d", http.StatusOK, response.Code)
	}
}
