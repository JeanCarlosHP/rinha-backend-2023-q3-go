package workerPool

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jeancarloshp/rinha-backend-go/internal/people"
)

var (
	MaxWorker = 1
	MaxQueue  = 1
)

type Task struct {
	Payload *people.People
}

type TaskQueue chan Task

func NewTaskQueue() TaskQueue {
	return make(TaskQueue, MaxQueue)
}

type Worker struct {
	WorkerPool  chan chan Task
	TaskChannel chan Task
	db          *pgxpool.Pool
}

func NewWorker(workerPool chan chan Task, db *pgxpool.Pool) Worker {
	return Worker{
		WorkerPool:  workerPool,
		TaskChannel: make(chan Task),
		db:          db,
	}
}

func (wp *Worker) Run() {
	dataCh := make(chan Task)
	insertCh := make(chan []Task)

	go wp.bootstrap(dataCh)

	go wp.processData(dataCh, insertCh)

	go wp.processInsert(insertCh)
}

func (w Worker) bootstrap(dataCh chan Task) {
	for {
		w.WorkerPool <- w.TaskChannel

		job := <-w.TaskChannel

		dataCh <- job
	}
}

func (w Worker) processData(dataCh chan Task, insertCh chan []Task) {
	tickInsertRate := time.Duration(7 * time.Second)
	tickInsert := time.NewTicker(tickInsertRate).C

	batchMaxSize := 10000
	batch := make([]Task, 0, batchMaxSize)

	for {
		select {
		case data := <-dataCh:
			batch = append(batch, data)

		case <-tickInsert:
			batchLen := len(batch)
			if batchLen > 0 {
				insertCh <- batch

				batch = make([]Task, 0, batchMaxSize)
			}
		}
	}
}

func (wp Worker) processInsert(insertCh chan []Task) {
	for payload := range insertCh {
		bb := bytes.Buffer{}
		bb.WriteString("INSERT INTO peoples (id, nickname, name, birth, stack) VALUES ")

		for i := 0; i < len(payload); i++ {
			task := payload[i]
			p := task.Payload

			bb.WriteString("('" + p.ID.String() + "', '" + p.Nickname + "', '" + p.Name + "', '" + p.Birth + "', '{" + strings.ToLower(strings.Join(p.Stack.Elements, ", ")) + "}'),")
		}

		bb.Truncate(bb.Len() - 1)

		_, err := wp.db.Exec(context.Background(), bb.String())
		if err != nil {
			fmt.Print(err)
		}
	}
}

type Dispatcher struct {
	maxWorkers int
	WorkerPool chan chan Task
	jobQueue   chan Task
	db         *pgxpool.Pool
}

func NewDispatcher(db *pgxpool.Pool, jobQueue TaskQueue) *Dispatcher {
	maxWorkers := MaxWorker

	pool := make(chan chan Task, maxWorkers)

	return &Dispatcher{
		WorkerPool: pool,
		maxWorkers: maxWorkers,
		jobQueue:   jobQueue,
		db:         db,
	}
}

func (d *Dispatcher) Run() {
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(d.WorkerPool, d.db)
		worker.Run()
	}

	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for job := range d.jobQueue {
		go func(job Task) {
			jobChannel := <-d.WorkerPool

			jobChannel <- job
		}(job)
	}
}
