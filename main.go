package main

import (
	"fmt"
	"html/template"
	"main/ascii-art"
	"net/http"
	"strings"
)

type ResultPageData struct {
	Result string
}

func handleForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("index.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func handleAsciiWeb(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		// Parse the form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		//get string and convert it to proper input for ascii-art function
		inputText := r.FormValue("text")
		inputText = strings.ReplaceAll(inputText, "\r\n", "\n")
		convertedText := strings.ReplaceAll(inputText, "\n", "\\n")

		banner := r.FormValue("banner")
		color := r.FormValue("color")

		flags := map[string]string{
			"color":  color,
			"align":  "",
			"output": "",
		}

		// Generate the ASCII art result
		res := ascii.HandleAsciiArt(convertedText, convertedText, banner, flags)

		// Prepare data for the result page
		resultData := ResultPageData{Result: res}

		// Parse the HTML template
		tmpl, err := template.ParseFiles("result.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Render the result page with the ASCII result
		w.Header().Set("Content-Type", "text/html")
		err = tmpl.Execute(w, resultData)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

		//fmt.Fprintf(w, "<h1>Form Submission Result</h1>")
		//fmt.Fprintln(w, "<textarea class='result-box'>"+res+"</textarea>")
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/", handleForm)
	http.HandleFunc("/ascii-web", handleAsciiWeb)
	// Start the server on port 8080
	fmt.Println("Starting server on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
