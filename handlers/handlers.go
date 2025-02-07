package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/Ardelean-Calin/cellulose/internal/db"
	database "github.com/Ardelean-Calin/cellulose/internal/db"
)

type App struct {
	db *database.DB
}

func NewApp(db *database.DB) *App {
	return &App{db}
}

func (a *App) UploadDocument(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form with 25MB max size
	r.ParseMultipartForm(25 << 20)

	file, handler, err := r.FormFile("pdf")
	if err != nil {
		http.Error(w, "Error retrieving PDF file", http.StatusBadRequest)
		return
	}
	defer file.Close()

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
	exists, err := a.db.DocumentExistsByHash(hashValue)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "File already exists", http.StatusBadRequest)
		return
	}

	// Add document to database
	log.Printf("Attempting to add document to database: %s (path: %s)\n", handler.Filename, filePath)
	doc, err := a.db.NewDocument(db.DocumentOptions{
		Title:   handler.Filename,
		Path:    filePath,
		Content: fmt.Sprintf("Uploaded document: %s", handler.Filename),
		Hash:    hashValue,
		Tags:    []string{}, // No tags initially
	})
	if err != nil {
		// Clean up file if database insert fails
		os.Remove(filePath)
		log.Printf("Failed to add document to database: %v\n", err)
		http.Error(w, fmt.Sprintf("Failed to add document to database: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Uploaded document: %s (ID: %d)\n", handler.Filename, doc.ID)
	w.Header().Set("HX-Trigger", "{\"documentUploaded\":null}")
	w.WriteHeader(http.StatusNoContent)
}

func (app *App) GetDocuments(w http.ResponseWriter, r *http.Request) {
	searchQuery := r.URL.Query().Get("search")

	documents, err := app.db.GetDocumentsByTitle(searchQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(documents)
}

func (app *App) GetTags(w http.ResponseWriter, r *http.Request) {
	tags, err := app.db.GetTags()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tags)
}

func (app *App) CreateTag(w http.ResponseWriter, r *http.Request) {
	var tagData struct {
		Name  string `json:"name"`
		Color string `json:"color"`
	}

	err := json.NewDecoder(r.Body).Decode(&tagData)
	if err != nil {
		log.Printf("Error decoding request body: %v\n", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	log.Printf("Received tag creation request with data: %+v\n", tagData)
	// Validate inputs
	if tagData.Name == "" || tagData.Color == "" {
		log.Printf("Missing required fields: name=%s, color=%s\n", tagData.Name, tagData.Color)
		http.Error(w, "Name and color are required", http.StatusBadRequest)
		return
	}

	// Validate hex color code
	if !regexp.MustCompile(`^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$`).MatchString(tagData.Color) {
		log.Printf("Invalid color code: %s\n", tagData.Color)
		http.Error(w, "Invalid hex color code", http.StatusBadRequest)
		return
	}

	tag, err := app.db.NewTag(tagData.Name, tagData.Color)
	if err != nil {
		log.Printf("Error creating tag: %v\n", err)
		if strings.Contains(err.Error(), "already exists") {
			http.Error(w, "Tag already exists", http.StatusUnprocessableEntity)
		} else {
			http.Error(w, "Failed to create tag", http.StatusInternalServerError)
		}
		return
	}

	log.Printf("Tag created successfully: %+v\n", tag)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tag)
}

// GetTagByID retrieves a tag by its ID.
func (app *App) GetTagByID(w http.ResponseWriter, r *http.Request) {
	log.Printf("GET Tag with ID: %s\n", r.PathValue("id"))

	// Parse the ID from the URL
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Get the tag from the database
	tag, err := app.db.GetTagByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Tag not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to get tag", http.StatusInternalServerError)
		}
		return
	}

	// Return the tag as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tag)
}

func (app *App) DeleteTagByID(w http.ResponseWriter, r *http.Request) {
	log.Printf("DELETE Tag with ID: %s\n", r.PathValue("id"))

	// Parse the ID from the URL
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Delete the tag from the database
	err = app.db.RemoveTag(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Tag not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to delete tag", http.StatusInternalServerError)
		}
		return
	}

	// Return success with no content
	w.WriteHeader(http.StatusNoContent)
}

// Get document by ID
func (app *App) GetDocumentByID(w http.ResponseWriter, r *http.Request) {}
