package main

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/Ardelean-Calin/cellulose/internal/data"
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	documents := data.GetSampleDocuments()
	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	tmpl.Execute(w, map[string]interface{}{
		"Documents": documents,
	})
}

func main() {
	http.HandleFunc("/", handleHome)
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
