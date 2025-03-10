package main

import (
	"product-tracker/config"
	"product-tracker/routes"
	"product-tracker/storage"

	_ "github.com/lib/pq"
)

func main() {
    config.LoadConfig()
    storageInstance, err := storage.NewStorage()
    if err != nil {
        panic("Failed to connect to the database")
    }
    defer storageInstance.Close()

    r := routes.SetupRouter()
    r.Run(":8080")
}