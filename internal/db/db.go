package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"
)

type DB struct {
	db *sql.DB
}

type Config struct {
	// DatabasePath is the path to the database file
	DatabasePath string
	// DocumentPath is the path to the directory where documents are stored
	DocumentPath string
}

// NewDB creates a new DB instance
func NewDB(config Config) (*DB, error) {
	db, err := sql.Open("sqlite", config.DatabasePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Create database directory if it doesn't exist
	err = os.MkdirAll(filepath.Dir(config.DatabasePath), 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	return &DB{db: db}, nil
}

func (db *DB) Init() {
	// Create documents table
	db.db.Exec(`
		CREATE TABLE IF NOT EXISTS documents (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			path TEXT NOT NULL,
			content TEXT NOT NULL,
			tags TEXT DEFAULT ''
		)
	`)

	// Create tags table
	db.db.Exec(`
		CREATE TABLE IF NOT EXISTS tags (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			color TEXT NOT NULL -- hex color code
		)
	`)
}

// AddDocument adds a document to the database
func (db *DB) AddDocument(document Document) (id int, err error) {
	err = db.db.QueryRow(`
		INSERT INTO documents (title, path, content, tags) VALUES (?, ?, ?, ?)
	`, document.Title, document.Path, document.Content, strings.Join(document.Tags, ",")).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to add document: %w", err)
	}

	return id, nil
}

// RemoveDocument removes a document from the database
func (db *DB) RemoveDocument(id int) error {
	var path string

	// Get document path
	err := db.db.QueryRow(`
		SELECT path FROM documents WHERE id = ?
	`, id).Scan(&path)
	if err != nil {
		return fmt.Errorf("failed to get document path: %w", err)
	}

	// Remove document from database
	_, err = db.db.Exec(`
		DELETE FROM documents WHERE id = ?
	`, id)
	if err != nil {
		return fmt.Errorf("failed to remove document from database: %w", err)
	}

	// Remove document from filesystem
	err = os.Remove(path)
	if err != nil {
		return fmt.Errorf("failed to remove document from filesystem: %w", err)
	}

	return nil
}

// Document represents a document in the database
type Document struct {
	Title   string   // title of the document
	Path    string   // path to the document
	Content string   // content of the document extracted from the file
	Tags    []string // tags associated with the document
}

// Tag represents a tag in the database
type Tag struct {
	Name  string // name of the tag
	Color string // hex color code
}

// AddTag adds a tag to the database
func (db *DB) AddTag(tag Tag) (id int, err error) {
	err = db.db.QueryRow(`
		INSERT INTO tags (name, color) VALUES (?, ?)
	`, tag.Name, tag.Color).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to add tag: %w", err)
	}

	return id, nil
}

// RemoveTag removes a tag from the database
func (db *DB) RemoveTag(id int) error {
	_, err := db.db.Exec(`
		DELETE FROM tags WHERE id = ?
	`, id)
	if err != nil {
		return fmt.Errorf("failed to remove tag from database: %w", err)
	}

	return nil
}
