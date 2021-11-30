// todo/repositories.go
package todo

import (
    "errors"

    "github.com/jinzhu/gorm"
)

type TodoRepository struct {
    database *gorm.DB
}

func (repository *TodoRepository) FindAllTodos() []Todo {
    var todos []Todo
    repository.database.Find(&todos)

    return todos
}

func (repository *TodoRepository) FindTodo(id int) (Todo, error) {
    var todo Todo
    err := repository.database.Find(&todo, id).Error
    if todo.Name == "" {
        err = errors.New("Todo not found")
    }
    return todo, err
}

func (repository *TodoRepository) CreateTodo(todo Todo) (Todo, error) {
    err := repository.database.Create(&todo).Error
    if err != nil {
        return todo, err
    }

    return todo, nil
}

func (repository *TodoRepository) SaveTodo(user Todo) (Todo, error) {
    err := repository.database.Save(user).Error
    return user, err
}

func (repository *TodoRepository) DeleteTodo(id int) int64 {
    count := repository.database.Delete(&Todo{}, id).RowsAffected
    return count
}

func NewTodoRepository(database *gorm.DB) *TodoRepository {
    return &TodoRepository{
        database: database,
    }
}