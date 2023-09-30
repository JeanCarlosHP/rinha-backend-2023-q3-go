package pessoas

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Pessoa struct {
	ID         uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Apelido    string         `gorm:"type:varchar(32);unique;not null;default:null" json:"apelido"`
	Nome       string         `gorm:"type:varchar(100);not null;default:null" json:"nome"`
	Nascimento string         `gorm:"type:date;not null;default:null" json:"nascimento"`
	Stack      pq.StringArray `gorm:"type:varchar(32)[]" json:"stack"`
}
