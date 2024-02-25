package handler

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jeancarloshp/rinha-backend-go/internal"
	"github.com/jeancarloshp/rinha-backend-go/internal/people"
	"github.com/jeancarloshp/rinha-backend-go/internal/people/repository"
)

type Handlers struct {
	PeopleRepository *repository.PeopleRepository
}

func (h *Handlers) HandlerCreatePeople(c *fiber.Ctx) error {
	p := new(people.PeopleDTO)

	if err := c.BodyParser(&p); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if p.Validate() != nil {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	validPeople := people.People{
		ID:       uuid.New(),
		Nickname: p.Nickname,
		Name:     p.Name,
		Birth:    p.Birth,
		Stack:    pgtype.Array[string]{Elements: p.Stack},
	}

	err := h.PeopleRepository.Create(&validPeople)
	if err != nil {
		if err.Error() == internal.ErrPeopleExists.Error() {
			return c.SendStatus(fiber.StatusUnprocessableEntity)
		}
	}

	c.Set("Location", fmt.Sprintf("/pessoas/%s", validPeople.ID))

	return c.SendStatus(fiber.StatusCreated)
}

func (h *Handlers) HandlerGetPeopleById(c *fiber.Ctx) error {
	id := c.Params("id")
	p, err := h.PeopleRepository.GetPeopleById(id)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	return c.JSON(
		people.PeopleDTO{
			Nickname: p.Nickname,
			Name:     p.Name,
			Birth:    p.Birth,
			Stack:    p.Stack.Elements,
		},
	)
}

func (h *Handlers) GetPeopleByTerm(c *fiber.Ctx) error {
	term := c.Query("t")
	if term == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	peoples, err := h.PeopleRepository.GetPeopleByTerm(term)
	if err != nil {
		return err
	}

	peoplesDTO := make([]people.PeopleDTO, len(*peoples))
	var wg sync.WaitGroup

	for i, p := range *peoples {
		wg.Add(1)
		go func(i int, p people.People) {
			defer wg.Done()
			peoplesDTO[i] = people.PeopleDTO{
				Nickname: p.Nickname,
				Name:     p.Name,
				Birth:    p.Birth,
				Stack:    p.Stack.Elements,
			}
		}(i, p)
	}
	wg.Wait()

	return c.JSON(&peoplesDTO)
}

func (h *Handlers) CountPeoples(c *fiber.Ctx) error {
	count, err := h.PeopleRepository.CountPeoples()
	if err != nil {
		return err
	}

	return c.SendString(strconv.FormatInt(count, 10))
}
