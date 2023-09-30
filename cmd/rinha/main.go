package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/jeancarloshp/rinha-backend-go/internal/pessoas"
	"github.com/jeancarloshp/rinha-backend-go/pkg/database/gorm/postgresql"
	"gorm.io/gorm"
)

var (
	app = fiber.New()

	db *gorm.DB

	err error
)

func main() {
	db, err = postgresql.New()
	if err != nil {
		panic(err)
	}

	if err := postgresql.Migrate(db); err != nil {
		panic(err)
	}

	handler := pessoas.Handlers{Db: db}

	// POST /pessoas – para criar um recurso pessoa.
	app.Post("/pessoas", handler.HandlerCreatePerson)

	// GET /pessoas/[:id] – para consultar um recurso criado com a requisição anterior.
	app.Get("/pessoas/:id", handler.HandlerGetPersonById)

	// GET /pessoas?t=[:termo da busca] – para fazer uma busca por pessoas.
	app.Get("/pessoas", handler.GetPersonByTerm)

	// GET /contagem-pessoas – endpoint especial para contagem de pessoas cadastradas.
	app.Get("/contagem-pessoas", handler.CoundPeople)

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}

	if err := app.Listen(fmt.Sprintf(":%s", PORT)); err != nil {
		panic(err)
	}
}
