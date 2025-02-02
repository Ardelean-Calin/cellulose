package main

import (
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strings"

	"github.com/Ardelean-Calin/cellulose/internal/data"
	"github.com/Ardelean-Calin/cellulose/internal/db"
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	tmpl.Execute(w, nil)
}

func handleDocuments(w http.ResponseWriter, r *http.Request) {
	database, err := db.NewDB(db.Config{
		DatabasePath: "cellulose.db",
		DocumentPath: "documents",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	documents, err := database.GetDocuments()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/partials/document_cards.html"))
	tmpl.Execute(w, documents)
}

// Helper function to get unique tags from documents
func getAllUniqueTags(documents []data.Document) []data.Tag {
	tagMap := make(map[string]data.Tag)
	for _, doc := range documents {
		for _, tag := range doc.Tags {
			tagMap[tag.Name] = tag
		}
	}

	uniqueTags := make([]data.Tag, 0, len(tagMap))
	for _, tag := range tagMap {
		uniqueTags = append(uniqueTags, tag)
	}

	// Sort tags by name (case-insensitive)
	sort.Slice(uniqueTags, func(i, j int) bool {
		return strings.ToLower(uniqueTags[i].Name) < strings.ToLower(uniqueTags[j].Name)
	})

	return uniqueTags
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form with 25MB max size
	r.ParseMultipartForm(25 << 20)
	
	file, handler, err := r.FormFile("pdf")
	if err != nil {
		http.Error(w, "Error retrieving PDF file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fmt.Printf("Received PDF upload: %s (Size: %d bytes)\n", handler.Filename, handler.Size)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Upload successful"))
}

func main() {
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/upload", handleUpload)
	http.HandleFunc("/api/documents", handleDocuments)
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
