package main

import (
	"fmt"
	"html/template"
	ascii "mymain/backend/services"
	"net/http"
	"os"
	"slices"
	"strconv"
)

var publicUrl = "frontend/public/"

type ResultPageData struct {
	Result string
	Color  string
	Align  string
}

type ErrorPageData struct {
	Name       string
	Code       string
	CodeNumber int
	Info       string
}

var PredefinedErrors = map[string]ErrorPageData{
	"BadRequestError": {
		Name:       "BadRequestError",
		Code:       strconv.Itoa(http.StatusBadRequest),
		CodeNumber: http.StatusBadRequest,
		Info:       "Bad request",
	},
	"NotFoundError": {
		Name:       "NotFoundError",
		Code:       strconv.Itoa(http.StatusNotFound),
		CodeNumber: http.StatusNotFound,
		Info:       "Page not found",
	},
	"MethodNotAllowedError": {
		Name:       "MethodNotAllowedError",
		Code:       strconv.Itoa(http.StatusMethodNotAllowed),
		CodeNumber: http.StatusMethodNotAllowed,
		Info:       "Method not allowed",
	},
	"InternalServerError": {
		Name:       "InternalServerError",
		Code:       strconv.Itoa(http.StatusInternalServerError),
		CodeNumber: http.StatusInternalServerError,
		Info:       "Internal server error",
	},
}

var (
	BadRequestError       = PredefinedErrors["BadRequestError"]
	NotFoundError         = PredefinedErrors["NotFoundError"]
	MethodNotAllowedError = PredefinedErrors["MethodNotAllowedError"]
	InternalServerError   = PredefinedErrors["InternalServerError"]
)

func handleForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if r.URL.Path != "/" {
			// If the URL is not exactly "/", respond with 404
			// handleNotFound(w, r)
			handleErrorPage(w, r, NotFoundError)
			return
		}
		tmpl, err := template.ParseFiles(publicUrl + "index.html")

		if err != nil {
			// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			handleErrorPage(w, r, InternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	} else {
		// http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		handleErrorPage(w, r, MethodNotAllowedError)
	}
}

func handleAsciiWeb(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Parse the form data
		err := r.ParseForm()
		if err != nil {
			// http.Error(w, "Bad Request", http.StatusBadRequest)
			handleErrorPage(w, r, BadRequestError)
			return
		}

		if len(r.FormValue("banner")) == 0 || len(r.FormValue("text")) == 0 {
			// handleBadRequest(w, r)
			handleErrorPage(w, r, BadRequestError)
			return
		}

		inputText := r.FormValue("text")
		banner := r.FormValue("banner")
		color := r.FormValue("color")
		align := r.FormValue("align")

		var banners = []string{"apple", "shadow", "standard", "thinkertoy"}

		if !slices.Contains(banners, banner) {
			// handleNotFound(w, r)
			handleErrorPage(w, r, NotFoundError)
			return
		}
		// // Read the banner file if exists
		_, err = os.Stat("backend/banners/" + banner + ".txt")
		if os.IsNotExist(err) {
			// handleServerErrors(w, r)
			handleErrorPage(w, r, InternalServerError)
			return
		}
		flags := map[string]string{
			"output": "",
		}

		// Generate the ASCII art result
		res := ascii.HandleAsciiArt(inputText, banner, flags)

		// Prepare data for the result page
		resultData := ResultPageData{Result: res, Color: color, Align: align}

		// Parse the HTML template
		tmpl, err := template.ParseFiles(publicUrl + "result.html")
		if err != nil {
			// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			handleErrorPage(w, r, InternalServerError)
			return
		}
		// Render the result page with the ASCII result
		w.Header().Set("Content-Type", "text/html")
		err = tmpl.Execute(w, resultData)
		if err != nil {
			// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			handleErrorPage(w, r, InternalServerError)
		}
	} else {
		// http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		handleErrorPage(w, r, MethodNotAllowedError)
	}
}

func handleErrorPage(w http.ResponseWriter, r *http.Request, errorType ErrorPageData) {
	tmpl, err := template.ParseFiles("frontend/errors/error.html")
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(errorType.CodeNumber)
	tmpl.Execute(w, errorType)
}

func main() {
	http.Handle("/static/", http.FileServer(http.Dir("./frontend/public/")))
	http.HandleFunc("/", handleForm)
	http.HandleFunc("/ascii-web", handleAsciiWeb)
	// http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./frontend/public/static/"))))
	// Start the server on port 8080
	fmt.Println("Starting server on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
