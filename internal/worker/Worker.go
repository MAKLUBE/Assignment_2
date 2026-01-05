package worker

import (
	"log"
	"time"

	"github.com/MAKLUBE/Assignment_2/internal/model"
)

func StartWorker(id int, queue <-chan *model.Task) {
	go func() {
		for task := range queue {
			log.Printf("worker %d processing task %s", id, task.Id)

			task.Status = model.InProgress
			time.Sleep(3 * time.Second)
			task.Status = model.Done

			log.Printf("worker %d finished task %s", id, task.Id)
		}
	}()
}
