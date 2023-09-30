package pessoas

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Handlers struct {
	Db *gorm.DB
}

func (h *Handlers) HandlerCreatePerson(c *fiber.Ctx) error {
	pessoa := new(Pessoa)
	if err := c.BodyParser(pessoa); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if err := h.Db.Create(pessoa).Error; err != nil {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	c.Set("Location", fmt.Sprintf("/pessoas/%s", pessoa.ID.String()))
	return c.SendStatus(fiber.StatusCreated)
}

func (h *Handlers) HandlerGetPersonById(c *fiber.Ctx) error {
	pessoa := new(Pessoa)

	if err := h.Db.Where("id = ?", c.Params("id")).First(&pessoa).Error; err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	return c.JSON(pessoa)
}

func (h *Handlers) GetPersonByTerm(c *fiber.Ctx) error {
	term := c.Query("t")
	if term == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	pessoas := []Pessoa{}

	h.Db.Where("LOWER(apelido) LIKE LOWER(?)", fmt.Sprintf("%%%s%%", term)).Or("LOWER(nome) LIKE LOWER(?)", fmt.Sprintf("%%%s%%", term)).Or("EXISTS ( SELECT 1 FROM unnest(stack) AS s WHERE LOWER(s) = LOWER(?))", term).Find(&pessoas).Limit(50).Debug()

	return c.JSON(&pessoas)
}

func (h *Handlers) CoundPeople(c *fiber.Ctx) error {
	var count int64 = 0
	h.Db.Model(&Pessoa{}).Count(&count)
	return c.SendString(strconv.FormatInt(count, 10))
}
