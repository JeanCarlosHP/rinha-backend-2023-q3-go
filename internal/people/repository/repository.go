package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jeancarloshp/rinha-backend-go/internal"
	"github.com/jeancarloshp/rinha-backend-go/internal/people"
	"github.com/jeancarloshp/rinha-backend-go/pkg/workerPool"
	"github.com/redis/go-redis/v9"
)

type PeopleRepository struct {
	conn *pgxpool.Pool
	tq   workerPool.TaskQueue
	rdb  *redis.Client
}

func NewPeopleRepository(conn *pgxpool.Pool, rdb *redis.Client, tq workerPool.TaskQueue) *PeopleRepository {
	return &PeopleRepository{
		conn: conn,
		rdb:  rdb,
		tq:   tq,
	}
}

func (ur *PeopleRepository) Create(p *people.People) error {
	people, err := ur.rdb.Get(context.Background(), p.Nickname).Result()
	if err != nil {
		if err != redis.Nil {
			return err
		}
	}

	if people == "1" {
		return internal.ErrPeopleExists
	}

	item, err := sonic.MarshalString(p)
	if err != nil {
		return err
	}

	ur.rdb.Set(context.Background(), p.ID.String(), item, 0)
	ur.rdb.Set(context.Background(), p.Nickname, 1, 0)

	ur.tq <- workerPool.Task{Payload: p}

	return err
}

func (ur *PeopleRepository) GetPeopleById(id string) (p *people.People, err error) {
	val, err := ur.rdb.Get(context.Background(), id).Result()

	if err != nil {
		if err == redis.Nil {
			return nil, internal.ErrPeopleNotFound
		}

		return nil, err
	}

	if val != "" {
		err = sonic.UnmarshalString(val, &p)
		if err != nil {
			return nil, err
		}

		return p, err
	}

	row := ur.conn.QueryRow(context.Background(), "SELECT id, nickname, name, birth, stack FROM peoples WHERE id = $1", id)

	p = &people.People{}
	birth := time.Time{}

	err = row.Scan(&p.ID, &p.Nickname, &p.Name, &birth, &p.Stack)
	if err != nil {
		return nil, internal.ErrPeopleNotFound
	}

	p.Birth = birth.Format("2006-01-02")

	return
}

func (ur *PeopleRepository) GetPeopleByTerm(term string) (p *[]people.People, err error) {
	rows, err := ur.conn.Query(context.Background(), "SELECT id, nickname, name, birth, stack FROM peoples WHERE searchable LIKE $1 LIMIT 50", strings.ToLower(fmt.Sprintf("%%%s%%", term)))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	p = &[]people.People{}

	for rows.Next() {
		var _p people.People
		birth := time.Time{}
		err = rows.Scan(&_p.ID, &_p.Nickname, &_p.Name, &birth, &_p.Stack)
		if err != nil {
			return nil, err
		}

		_p.Birth = birth.Format("2006-01-02")

		*p = append(*p, _p)
	}

	return
}

func (ur *PeopleRepository) CountPeoples() (count int64, err error) {
	row := ur.conn.QueryRow(context.Background(), "SELECT COUNT(*) FROM peoples")
	err = row.Scan(&count)

	return
}
