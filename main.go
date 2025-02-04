package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/Ardelean-Calin/cellulose/internal/data"
	"github.com/Ardelean-Calin/cellulose/internal/db"
)

func handleDocuments(w http.ResponseWriter, r *http.Request) {
	database, err := db.NewDB(db.Config{
		DatabasePath: "cellulose.db",
		DocumentPath: "documents",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	searchQuery := r.URL.Query().Get("search")

	documents, err := database.GetDocumentsByTitle(searchQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(documents)
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

func handleTags(w http.ResponseWriter, r *http.Request) {
	database, err := db.NewDB(db.Config{
		DatabasePath: "cellulose.db",
		DocumentPath: "documents",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tags, err := database.GetTags()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tags)
}

func handleCreateTag(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		return
	}

	var tagData struct {
		Name  string `json:"name"`
		Color string `json:"color"`
	}

	fmt.Printf("Received tag creation request with data: %+v\n", tagData)

	err := json.NewDecoder(r.Body).Decode(&tagData)
	if err != nil {
		fmt.Printf("Error decoding request body: %v\n", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate inputs
	if tagData.Name == "" || tagData.Color == "" {
		fmt.Printf("Missing required fields: name=%s, color=%s\n", tagData.Name, tagData.Color)
		http.Error(w, "Name and color are required", http.StatusBadRequest)
		return
	}

	// Validate hex color code
	if !regexp.MustCompile(`^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$`).MatchString(tagData.Color) {
		fmt.Printf("Invalid color code: %s\n", tagData.Color)
		http.Error(w, "Invalid hex color code", http.StatusBadRequest)
		return
	}

	database, err := db.NewDB(db.Config{
		DatabasePath: "cellulose.db",
		DocumentPath: "documents",
	})
	if err != nil {
		fmt.Printf("Database connection failed: %v\n", err)
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}

	tag, err := database.NewTag(tagData.Name, tagData.Color)
	if err != nil {
		fmt.Printf("Error creating tag: %v\n", err)
		if strings.Contains(err.Error(), "already exists") {
			http.Error(w, "Tag already exists", http.StatusUnprocessableEntity)
		} else {
			http.Error(w, "Failed to create tag", http.StatusInternalServerError)
		}
		return
	}

	fmt.Printf("Tag created successfully: %+v\n", tag)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tag)
}

func main() {
	http.HandleFunc("/upload", handleUpload)

	// POST-ing to /api/documents will upload a new document using handleUpload while GET will return the list of available documents (as JSON). AI!
	http.HandleFunc("/api/documents", handleDocuments)

	http.HandleFunc("/api/tags", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handleTags(w, r)
		} else if r.Method == http.MethodPost {
			handleCreateTag(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
