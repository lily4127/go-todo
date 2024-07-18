package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

type Todo struct {
	ID   int    `json:"id"`
	Task string `json:"task"`
	Done bool   `json:"done"`
}

func main() {
	e := echo.New()

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s",
		getEnv("DB_USER", "root"),
		getEnv("DB_PASSWORD", "password"),
		getEnv("DB_HOST", "db"),
		getEnv("DB_NAME", "todo_app"),
	))
	if err != nil {
		e.Logger.Fatal(err)
	}

	defer db.Close()

	e.GET("/todos", func(c echo.Context) error {
		rows, err := db.Query("SELECT id, task, done FROM todos")
		if err != nil {
			return err
		}
		defer rows.Close()

		var todos []Todo
		for rows.Next() {
			var todo Todo
			if err := rows.Scan(&todo.ID, &todo.Task, &todo.Done); err != nil {
				return err
			}
			todos = append(todos, todo)
		}

		return c.JSON(http.StatusOK, todos)
	})

	e.POST("/todos", func(c echo.Context) error {
		todo := new(Todo)
		if err := c.Bind(todo); err != nil {
			return err
		}

		_, err := db.Exec("INSERT INTO todos (task, done) VALUES (?, ?)", todo.Task, todo.Done)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusCreated, todo)
	})

	e.Logger.Fatal(e.Start(":8080"))
}

func getEnv(key, defaultVal string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultVal
	}
	return value
}
