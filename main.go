package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"crypto/sha256"
	"encoding/hex"
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

	fmt.Printf("Retrieved %d documents from database\n", len(documents))
	for _, doc := range documents {
		fmt.Printf("Document: ID=%d, Title=%s, Path=%s, Content=%s\n", 
			doc.ID, doc.Opts.Title, doc.Opts.Path, doc.Opts.Content)
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

	// Create database connection
	database, err := db.NewDB(db.Config{
		DatabasePath: "cellulose.db",
		DocumentPath: "documents",
	})
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Initialize database if needed
	database.Init()

	// Create documents directory if it doesn't exist
	err = os.MkdirAll("documents", 0755)
	if err != nil {
		http.Error(w, "Failed to create documents directory", http.StatusInternalServerError)
		return
	}

	// Save file to disk and compute hash
	filePath := filepath.Join("documents", handler.Filename)
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	hash := sha256.New()
	writer := io.MultiWriter(dst, hash)

	if _, err = io.Copy(writer, file); err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	hashValue := hex.EncodeToString(hash.Sum(nil))

	// Check if the hash already exists
	exists, err := database.DocumentExistsByHash(hashValue)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "File already exists", http.StatusBadRequest)
		return
	}

	// Add document to database
	fmt.Printf("Attempting to add document to database: %s (path: %s)\n", handler.Filename, filePath)
	doc, err := database.NewDocument(db.DocumentOptions{
		Title:   handler.Filename,
		Path:    filePath,
		Content: fmt.Sprintf("Uploaded document: %s", handler.Filename),
		Hash:    hashValue,
		Tags:    []string{}, // No tags initially
	})
	if err != nil {
		// Clean up file if database insert fails
		os.Remove(filePath)
		fmt.Printf("Failed to add document to database: %v\n", err)
		http.Error(w, fmt.Sprintf("Failed to add document to database: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Printf("Uploaded document: %s (ID: %d)\n", handler.Filename, doc.ID)
	w.Header().Set("HX-Trigger", "{\"documentUploaded\":null}")
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/upload", handleUpload)
	http.HandleFunc("/api/documents", handleDocuments)
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
