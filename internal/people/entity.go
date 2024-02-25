package people

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type People struct {
	ID         uuid.UUID            `gorm:"type:uuid;primary_key" json:"id"`
	Nickname   string               `gorm:"type:varchar(32);unique;not null" json:"apelido"`
	Name       string               `gorm:"type:varchar(100);not null" json:"nome"`
	Birth      string               `gorm:"type:date;not null" json:"nascimento"`
	Stack      pgtype.Array[string] `gorm:"type:varchar(32)[]" json:"stack"`
	Searchable string               `gorm:"-" json:"-"`
}

type PeopleDTO struct {
	Nickname string   `json:"apelido" validate:"required,max=32"`
	Name     string   `json:"nome" validate:"required,max=100"`
	Birth    string   `json:"nascimento" validate:"required,datetime=2006-01-02"`
	Stack    []string `json:"stack" validate:"dive,max=32"`
}

func (p *PeopleDTO) Validate() error {
	if p.Nickname == "" || p.Name == "" || len(p.Nickname) > 32 || len(p.Name) > 100 {
		return errors.New("Unprocessable")
	}

	if !isValidDate(p.Birth) {
		return errors.New("Unprocessable")
	}

	for _, s := range p.Stack {
		if len(s) > 32 {
			return errors.New("Unprocessable")
		}
	}

	return nil
}

func isValidDate(dateStr string) bool {
	layout := "2006-01-02"

	_, err := time.Parse(layout, dateStr)
	return err == nil
}
