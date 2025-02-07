package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Ardelean-Calin/cellulose/handlers"
	"github.com/Ardelean-Calin/cellulose/internal/db"
	"github.com/Ardelean-Calin/cellulose/middleware"
)

func main() {
	database, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	// Create app with dependencies
	app := handlers.NewApp(database)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/documents", app.UploadDocument)
	mux.HandleFunc("GET /api/documents", app.GetDocuments)
	// mux.HandleFunc("PUT /api/documents/{id}", handler.UpdateByID)
	mux.HandleFunc("GET /api/documents/{id}", app.GetDocumentByID)
	// mux.HandleFunc("DELETE /api/documents/{id}", handler.DeleteByID)

	mux.HandleFunc("POST /api/tags", app.CreateTag)
	mux.HandleFunc("GET /api/tags", app.GetTags)
	mux.HandleFunc("GET /api/tags/{id}", app.GetTagByID)
	mux.HandleFunc("DELETE /api/tags/{id}", app.DeleteTagByID)

	fmt.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", middleware.Logging(mux)); err != nil {
		log.Panic(err)
	}
}
