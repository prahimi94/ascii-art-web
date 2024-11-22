package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"
)

func TestHandleForm(t *testing.T) {
	// Change to the root directory of the project
	err := os.Chdir("../../")
	if err != nil {
		t.Fatalf("Error changing directory: %v", err)
	}
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
	expected := `<p>ASCII-ART-WEB by PR & MR</p>` // Add a unique identifier that would appear in the HTML

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

func TestHandleErrorPage(t *testing.T) {
	testCases := []struct {
		name         string
		errorType    ErrorPageData
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Test 404 Not Found",
			errorType:    NotFoundError,
			expectedCode: http.StatusNotFound,
			expectedBody: "Page not found",
		},
		{
			name:         "Test 400 Bad Request",
			errorType:    BadRequestError,
			expectedCode: http.StatusBadRequest,
			expectedBody: "Bad request",
		},
		{
			name:         "Test 500 Internal Server Error",
			errorType:    InternalServerError,
			expectedCode: http.StatusInternalServerError,
			expectedBody: "Internal server error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()

			handleErrorPage(w, req, tc.errorType)

			// get result
			res := w.Result()
			defer res.Body.Close()

			// check http status code
			if res.StatusCode != tc.expectedCode {
				t.Errorf("Expected status code %d, but got %d", tc.expectedCode, res.StatusCode)
			}

			// check response content
			body := w.Body.String()
			if !strings.Contains(body, tc.expectedBody) {
				t.Errorf("Expected body to contain '%s', but got '%s'", tc.expectedBody, body)
			}
		})
	}
}

func TestHandleAsciiWeb(t *testing.T) {

	// Prepare the handler and create a request
	handler := http.HandlerFunc(handleAsciiWeb)

	// Test valid POST request with valid form data
	formData := "banner=standard&text=hello&color=red&align=center"
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
	formData = "banner=standard&text=&color=red&align=center"
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

	// Audit tests
	testCases := []struct {
		name         string
		formData     string
		expectedFile string
		expectedCode int
	}{
		{
			"Case 1: 123 and <Hello> (World)!",
			`banner=standard&text={123}\n<Hello> (World)!`,
			"backend/api/tests/1.txt",
			http.StatusOK,
		},

		{
			"Case 2: 123??",
			`banner=standard&text=123??`,
			"backend/api/tests/2.txt",
			http.StatusOK,
		},

		{
			"Case 3: $% \"= (shadow)",
			`banner=shadow&text=%24%25%20%22%3D`,
			"backend/api/tests/3.txt",
			http.StatusOK,
		},

		{
			"Case 4: 123 T/fs#R (thinkertoy)",
			`banner=thinkertoy&text=123%20T%2Ffs%23R`,
			"backend/api/tests/4.txt",
			http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/ascii-web", strings.NewReader(tc.formData))
			if err != nil {
				t.Fatalf("Failed to create POST request: %v", err)
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			// Verify HTTP status code
			if rr.Code != tc.expectedCode {
				t.Errorf("Expected status %v, got %v", tc.expectedCode, rr.Code)
			}
			// Load expected output from file and compare
			expectedContent, err := ioutil.ReadFile(tc.expectedFile)
			if err != nil {
				t.Fatalf("Failed to read expected file %v: %v", tc.expectedFile, err)
			}
			if extractDivValueByID(rr.Body.String()) != string(expectedContent) {
				t.Errorf("Expected output:\n%v\nGot:\n%v", string(expectedContent), extractDivValueByID(rr.Body.String()))
			}
		})
	}
}
func extractTextareaText(htmlContent string) string {
	// Define a regular expression to match the content inside the <textarea> with id="result"
	re := regexp.MustCompile(`<textarea[^>]*>(.*?)</textarea>`)
	match := re.FindStringSubmatch(htmlContent)
	if len(match) < 2 {
		return "er"
	}

	// The inner content will be the second element in the match slice
	return match[1]
}

func extractDivValueByID(html string) string {
	startTag := `<textarea class="result-box" id="result" style="color:; text-align:;">`
	endTag := `</textarea>`

	// Find the start index of the desired <div>
	startIndex := strings.Index(html, startTag)
	if startIndex == -1 {
		return "Error - Call Support"
	}
	startIndex += len(startTag)

	// Find the end index of the </div>
	endIndex := strings.Index(html[startIndex:], endTag)
	if endIndex == -1 {
		return "Error - Call Support"
	}
	// Extract and return the content
	return decodeHTMLEntities(html[startIndex : startIndex+endIndex])
}

// Function to replace HTML entities manually
func decodeHTMLEntities(input string) string {
	// Replacing the common HTML entities with their corresponding characters
	replacements := map[string]string{
		"&lt;":   "<",
		"&gt;":   ">",
		"&amp;":  "&",
		"&quot;": "\"",
		"&#39;":  "'",
	}

	// Loop through the map and replace the entities in the input string
	for entity, char := range replacements {
		input = strings.ReplaceAll(input, entity, char)
	}

	return input
}
