package main

import (
	"fmt"
	"html/template"
	ascii "mymain/backend/services"
	"net/http"
	"os"
	"slices"
)

var publicUrl = "frontend/public/"

type ResultPageData struct {
	Result string
	Color  string
	Align  string
}

func handleForm(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		handleNotFound(w, r)
		return
	}
	if r.Method == http.MethodGet {
		fmt.Println(3)
		tmpl, err := template.ParseFiles(publicUrl + "index.html")

		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}

}

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	// If the URL is not exactly "/", respond with 404
	tmpl, err := template.ParseFiles("frontend/errors/404.html")
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	tmpl.Execute(w, nil)
}

func handleServerErrors(w http.ResponseWriter, r *http.Request) {
	// If there is an internal server error "/", respond with 500
	tmpl, err := template.ParseFiles("frontend/errors/500.html")
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
	tmpl.Execute(w, nil)
}

func handleBadRequest(w http.ResponseWriter, r *http.Request) {
	// If the request has problems, respond with 400
	tmpl, err := template.ParseFiles("frontend/errors/400.html")
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	tmpl.Execute(w, nil)
}

func handleAsciiWeb(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		// Parse the form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		if len(r.FormValue("banner")) == 0 || len(r.FormValue("text")) == 0 {
			handleBadRequest(w, r)
			return
		}

		inputText := r.FormValue("text")
		banner := r.FormValue("banner")
		color := r.FormValue("color")
		align := r.FormValue("align")

		var banners = []string{"apple", "shadow", "standard", "thinkertoy"}

		if !slices.Contains(banners, banner) {
			handleNotFound(w, r)
			return
		}
		// // Read the banner file if exists
		_, err = os.Stat("backend/banners/" + banner + ".txt")
		if os.IsNotExist(err) {
			handleServerErrors(w, r)
			return
		}

		flags := map[string]string{
			"color":  "",
			"align":  "",
			"output": "",
		}

		// Generate the ASCII art result
		res := ascii.HandleAsciiArt(inputText, inputText, banner, flags)

		// Prepare data for the result page
		resultData := ResultPageData{Result: res, Color: color, Align: align}

		// Parse the HTML template
		tmpl, err := template.ParseFiles(publicUrl + "result.html")
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
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./frontend/public/static/"))))
	http.HandleFunc("/", handleForm)
	http.HandleFunc("/ascii-web", handleAsciiWeb)
	// Start the server on port 8080
	fmt.Println("Starting server on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
