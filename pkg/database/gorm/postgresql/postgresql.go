package postgresql

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	entity "github.com/jeancarloshp/rinha-backend-go/internal/pessoas"
)

var (
	dsn = "host=db user=postgres password=postgres dbname=rinha port=5432 sslmode=disable TimeZone=America/Sao_Paulo"
)

func New() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db, nil
}

func Migrate(db *gorm.DB) (err error) {
	if !db.Migrator().HasTable(&entity.Pessoa{}) {
		err = db.Migrator().CreateTable(&entity.Pessoa{})
	}

	return
}
