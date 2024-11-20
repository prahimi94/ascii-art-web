package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleForm(t *testing.T) {
	// Test for GET request
	log.Println("Starting TestHandleForm - GET request")
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		log.Fatalf("Failed to create GET request: %v", err)
	}

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleForm)

	// Perform the GET request
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		log.Printf("TestHandleForm - GET request: Expected status %v, got %v", http.StatusOK, status)
		t.Errorf("handleForm returned wrong status code: got %v want %v", status, http.StatusOK)
	} else {
		log.Printf("TestHandleForm - GET request: Received status %v, as expected", status)
	}

	// Check if the response contains the correct content (assuming the HTML file is served)
	expected := "<html>" // Add a unique identifier that would appear in the HTML
	if !strings.Contains(rr.Body.String(), expected) {
		log.Printf("TestHandleForm - GET request: Expected content not found in response body")
		t.Errorf("handleForm returned unexpected body: got %v want %v", rr.Body.String(), expected)
	} else {
		log.Printf("TestHandleForm - GET request: Correct content found in response body")
	}

	// Test for POST request (Should return Method Not Allowed)
	req, err = http.NewRequest("POST", "/", nil)
	if err != nil {
		log.Fatalf("Failed to create POST request: %v", err)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		log.Printf("TestHandleForm - POST request: Expected status %v, got %v", http.StatusMethodNotAllowed, status)
		t.Errorf("handleForm returned wrong status code for POST: got %v want %v", status, http.StatusMethodNotAllowed)
	} else {
		log.Printf("TestHandleForm - POST request: Received status %v, as expected", status)
	}
}

func TestHandleNotFound(t *testing.T) {
	log.Println("Starting TestHandleNotFound")
	req, err := http.NewRequest("GET", "/nonexistent", nil)
	if err != nil {
		log.Fatalf("Failed to create GET request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleNotFound)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		log.Printf("TestHandleNotFound: Expected status %v, got %v", http.StatusNotFound, status)
		t.Errorf("handleNotFound returned wrong status code: got %v want %v", status, http.StatusNotFound)
	} else {
		log.Printf("TestHandleNotFound: Received status %v, as expected", status)
	}
}

func TestHandleServerErrors(t *testing.T) {
	log.Println("Starting TestHandleServerErrors")
	req, err := http.NewRequest("GET", "/server-error", nil)
	if err != nil {
		log.Fatalf("Failed to create GET request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleServerErrors)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		log.Printf("TestHandleServerErrors: Expected status %v, got %v", http.StatusInternalServerError, status)
		t.Errorf("handleServerErrors returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	} else {
		log.Printf("TestHandleServerErrors: Received status %v, as expected", status)
	}
}

func TestHandleBadRequest(t *testing.T) {
	log.Println("Starting TestHandleBadRequest")
	req, err := http.NewRequest("GET", "/bad-request", nil)
	if err != nil {
		log.Fatalf("Failed to create GET request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleBadRequest)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		log.Printf("TestHandleBadRequest: Expected status %v, got %v", http.StatusBadRequest, status)
		t.Errorf("handleBadRequest returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	} else {
		log.Printf("TestHandleBadRequest: Received status %v, as expected", status)
	}
}

func TestHandleAsciiWeb(t *testing.T) {
	// Prepare the handler and create a request
	handler := http.HandlerFunc(handleAsciiWeb)

	// Test valid POST request with valid form data
	formData := "banner=apple&text=hello&color=red&align=center"
	log.Printf("TestHandleAsciiWeb - Valid request: %s", formData)
	req, err := http.NewRequest("POST", "/ascii-web", strings.NewReader(formData))
	if err != nil {
		t.Fatalf("Failed to create POST request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Check if the response code is 200 OK
	if status := rr.Code; status != http.StatusOK {
		log.Printf("TestHandleAsciiWeb - Valid request: Expected status %v, got %v", http.StatusOK, status)
		t.Errorf("handleAsciiWeb returned wrong status code for valid request: got %v want %v", status, http.StatusOK)
	} else {
		log.Printf("TestHandleAsciiWeb - Valid request: Received status %v, as expected", status)
	}

	// Test invalid POST data (missing 'text' field)
	formData = "banner=apple&text=&color=red&align=center"
	log.Printf("TestHandleAsciiWeb - Invalid POST data: %s", formData)
	req, err = http.NewRequest("POST", "/ascii-web", strings.NewReader(formData))
	if err != nil {
		t.Fatalf("Failed to create POST request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Check if the response code is 400 Bad Request
	if status := rr.Code; status != http.StatusBadRequest {
		log.Printf("TestHandleAsciiWeb - Invalid request: Expected status %v, got %v", http.StatusBadRequest, status)
		t.Errorf("handleAsciiWeb returned wrong status code for bad request: got %v want %v", status, http.StatusBadRequest)
	} else {
		log.Printf("TestHandleAsciiWeb - Invalid request: Received status %v, as expected", status)
	}

	// Test invalid banner (banner is not part of allowed banners)
	formData = "banner=invalid-banner&text=hello&color=red&align=center"
	log.Printf("TestHandleAsciiWeb - Invalid banner: %s", formData)
	req, err = http.NewRequest("POST", "/ascii-web", strings.NewReader(formData))
	if err != nil {
		t.Fatalf("Failed to create POST request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Check if the response code is 404 Not Found
	if status := rr.Code; status != http.StatusNotFound {
		log.Printf("TestHandleAsciiWeb - Invalid banner: Expected status %v, got %v", http.StatusNotFound, status)
		t.Errorf("handleAsciiWeb returned wrong status code for invalid banner: got %v want %v", status, http.StatusNotFound)
	} else {
		log.Printf("TestHandleAsciiWeb - Invalid banner: Received status %v, as expected", status)
	}

	// Test missing banner and text (empty fields should result in 400)
	formData = "banner=&text=&color=red&align=center"
	log.Printf("TestHandleAsciiWeb - Missing banner and text: %s", formData)
	req, err = http.NewRequest("POST", "/ascii-web", strings.NewReader(formData))
	if err != nil {
		t.Fatalf("Failed to create POST request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Check if the response code is 400 Bad Request
	if status := rr.Code; status != http.StatusBadRequest {
		log.Printf("TestHandleAsciiWeb - Missing banner and text: Expected status %v, got %v", http.StatusBadRequest, status)
		t.Errorf("handleAsciiWeb returned wrong status code for missing banner/text: got %v want %v", status, http.StatusBadRequest)
	} else {
		log.Printf("TestHandleAsciiWeb - Missing banner and text: Received status %v, as expected", status)
	}
}
