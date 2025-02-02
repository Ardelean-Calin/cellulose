package main

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"text/template"

	"github.com/Ardelean-Calin/cellulose/internal/data"
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	documents := data.GetSampleDocuments()
	allTags := getAllUniqueTags(documents)

	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	tmpl.Execute(w, map[string]interface{}{
		"Documents": documents,
		"AllTags":   allTags,
	})
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
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
