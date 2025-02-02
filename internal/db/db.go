package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

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
			tags INTEGER[] DEFAULT '{}' -- comma separated list of tag ids
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

// NewDocument adds a document to the database
func (db *DB) NewDocument(opts DocumentOptions) (Document, error) {
	// Verify that the tags exist
	for _, tag := range opts.Tags {
		var tagID int
		err := db.db.QueryRow(`
			SELECT id FROM tags WHERE name = ?
		`, tag).Scan(&tagID)
		if err != nil {
			return Document{}, fmt.Errorf("failed to verify tag: %w", err)
		}
	}

	var id int
	err := db.db.QueryRow(`
		INSERT INTO documents (title, path, content, tags) VALUES (?, ?, ?, ?)
	`, opts.Title, opts.Path, opts.Content, opts.Tags).Scan(&id)
	if err != nil {
		return Document{}, fmt.Errorf("failed to add document: %w", err)
	}

	return Document{ID: id, Opts: opts}, nil
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
	ID   int // id of the document
	Opts DocumentOptions
}

type DocumentOptions struct {
	Title   string
	Path    string
	Content string
	Tags    []string
}

// Tag represents a tag in the database
type Tag struct {
	ID    int    // id of the tag
	Name  string // name of the tag
	Color string // hex color code
}

// NewTag creates a new tag inside the database
func (db *DB) NewTag(name string, color string) (Tag, error) {
	var tag Tag
	err := db.db.QueryRow(`
		INSERT INTO tags (name, color) VALUES (?, ?)
	`, name, color).Scan(&tag.ID, &tag.Name, &tag.Color)
	if err != nil {
		return Tag{}, fmt.Errorf("failed to add tag: %w", err)
	}
	return tag, nil
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

// GetDocuments returns all documents in the database
func (db *DB) GetDocuments() ([]Document, error) {
	rows, err := db.db.Query(`
		SELECT id, title, path, content FROM documents
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get documents: %w", err)
	}
	defer rows.Close()

	var documents []Document
	for rows.Next() {
		var doc Document
		err = rows.Scan(&doc.ID, &doc.Opts.Title, &doc.Opts.Path, &doc.Opts.Content)
		if err != nil {
			return nil, fmt.Errorf("failed to scan document: %w", err)
		}
		documents = append(documents, doc)
	}

	return documents, nil
}

// GetTags returns all tags in the database
func (db *DB) GetTags() ([]Tag, error) {
	rows, err := db.db.Query(`
		SELECT id, name, color FROM tags
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}
	defer rows.Close()

	tags := []Tag{}
	for rows.Next() {
		var tag Tag
		err = rows.Scan(&tag.Name, &tag.Color)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tag: %w", err)
		}
		tags = append(tags, tag)
	}

	return tags, nil
}
