package main

import (
	"context"
	"log"
	"os"

	"discord/internal/database"
	"discord/internal/server"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки файла .env")
	}
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintf(w, "Hello, World!")
	// })
	// http.ListenAndServe(":8080", nil)

	ctx := context.Background()

	databaseURL := os.Getenv("DATABASE_URL")

	db, err := database.NewPostgresPool(ctx, databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Println("Connected to PostgreSQL")

	server := server.NewServer(db)

	log.Println("Starting server on :8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

	// hub := websocket.NewHub()

	// go hub.Run()
}
