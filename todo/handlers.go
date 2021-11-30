// todo/handlers.go
package todo

import (
    "strconv"

    "github.com/gofiber/fiber/v2"
    "github.com/jinzhu/gorm"
)

type TodoHandler struct {
    repository *TodoRepository
}

func (handler *TodoHandler) GetAllTodos(c *fiber.Ctx) error {
    var todos []Todo = handler.repository.FindAllTodos()
    return c.JSON(todos)
}

func (handler *TodoHandler) GetTodo(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        print(err)
    }

    todo, err := handler.repository.FindTodo(id)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{
            "status": 404,
            "error":  err,
        })
    }

    return c.JSON(todo)
}

func (handler *TodoHandler) CreateNewTodo(c *fiber.Ctx) error {
    data := new(Todo)

    if err := c.BodyParser(data); err != nil {
        return c.Status(500).JSON(fiber.Map{
            "status": "error",
            "message": "Review your input",
            "error": err,
        })
    }

    item, err := handler.repository.CreateTodo(*data)

    if err != nil {
        return c.Status(400).JSON(fiber.Map{
            "status":  400,
            "message": "Failed creating item",
            "error":   err,
        })
    }

    return c.JSON(item)
}

func (handler *TodoHandler) UpdateTodo(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))

    if err != nil {
        return c.Status(400).JSON(fiber.Map{
            "status":  400,
            "message": "Item not found",
            "error":   err,
        })
    }

    todo, err := handler.repository.FindTodo(id)

    if err != nil {
        return c.Status(404).JSON(fiber.Map{
            "message": "Item not found",
        })
    }

    todoData := new(Todo)

    if err := c.BodyParser(todoData); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "status": "error",
            "message": "Review your input",
            "data": err,
        })
    }

    todo.Name = todoData.Name
    todo.Description = todoData.Description
    todo.Status = todoData.Status

    item, err := handler.repository.SaveTodo(todo)

    if err != nil {
        return c.Status(400).JSON(fiber.Map{
            "message": "Error updating todo",
            "error":   err,
        })
    }

    return c.JSON(item)
}

func (handler *TodoHandler) DeleteTodo(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(400).JSON(fiber.Map{
            "status":  400,
            "message": "Failed deleting todo",
            "err":     err,
        })
    }

    RowsAffected := handler.repository.DeleteTodo(id)
    statusCode := 204
    if RowsAffected == 0 {
        statusCode = 400
    }

    return c.Status(statusCode).JSON(nil)
}

func NewTodoHandler(repository *TodoRepository) *TodoHandler {
    return &TodoHandler{
        repository: repository,
    }
}

func Register(router fiber.Router, database *gorm.DB) {
    database.AutoMigrate(&Todo{})
    todoRepository := NewTodoRepository(database)
    todoHandler := NewTodoHandler(todoRepository)

    todoRouter := router.Group("/todo")
    todoRouter.Get("/", todoHandler.GetAllTodos)
    todoRouter.Get("/:id", todoHandler.GetTodo)
    todoRouter.Put("/:id", todoHandler.UpdateTodo)
    todoRouter.Post("/", todoHandler.CreateNewTodo)
    todoRouter.Delete("/:id", todoHandler.DeleteTodo)
}
