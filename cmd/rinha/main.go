package main

import (
	"context"
	"fmt"
	"os"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jeancarloshp/rinha-backend-go/internal/people/handler"
	"github.com/jeancarloshp/rinha-backend-go/internal/people/repository"

	"github.com/jeancarloshp/rinha-backend-go/pkg/workerPool"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var (
	app         *fiber.App
	databaseUrl string
	port        string
)

func init() {
	env := os.Getenv("RINHA_ENV")
	if env == "" {
		err := godotenv.Load(".env.dev")
		if err != nil {
			panic(fmt.Errorf("error loading .env file: %s", err))
		}
	}

	databaseUrl = os.Getenv("DATABASE_URL")
	if port = os.Getenv("PORT"); port == "" {
		port = "8080"
	}
}

func main() {
	conn, err := pgxpool.New(context.Background(), databaseUrl)
	if err != nil {
		panic(fmt.Sprintf("Unable to connect to database: %s", err))
	}
	defer conn.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	taskQueue := workerPool.NewTaskQueue()

	dispather := workerPool.NewDispatcher(conn, taskQueue)
	go dispather.Run()

	peopleRepository := repository.NewPeopleRepository(conn, rdb, taskQueue)

	handler := handler.Handlers{PeopleRepository: peopleRepository}

	fiberCfg := fiber.Config{
		AppName:       "rinha-go by @jeancarloshp",
		CaseSensitive: true,
	}

	fiberCfg.JSONEncoder = sonic.Marshal
	fiberCfg.JSONDecoder = sonic.Unmarshal

	app = fiber.New()

	// POST /pessoas – para criar um recurso pessoa.
	app.Post("/pessoas", handler.HandlerCreatePeople)

	// // GET /pessoas/[:id] – para consultar um recurso criado com a requisição anterior.
	app.Get("/pessoas/:id", handler.HandlerGetPeopleById)

	// // GET /pessoas?t=[:termo da busca] – para fazer uma busca por pessoas.
	app.Get("/pessoas", handler.GetPeopleByTerm)

	// // GET /contagem-pessoas – endpoint especial para contagem de pessoas cadastradas.
	app.Get("/contagem-pessoas", handler.CountPeoples)

	if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
		panic(err)
	}
}
