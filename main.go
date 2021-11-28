// main.go
package main

import (
    "log"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/shashoto-nur/go-todo/database"
    "github.com/shashoto-nur/go-todo/todo"
)

func main() {
    app := fiber.New()
    app.Use(cors.New())
    database.ConnectDB()
    defer database.DB.Close()

    api := app.Group("/api")
    todo.Register(api, database.DB)

    log.Fatal(app.Listen(":5000"))
}