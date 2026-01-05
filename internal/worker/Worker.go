package worker

import (
	"log"
	"time"

	"github.com/MAKLUBE/Assignment_2/internal/model"
)

func StartWorker(
	id int,
	queue <-chan *model.Task,
	stop <-chan struct{},
) {
	go func() {
		for {
			select {
			case task, ok := <-queue:
				if !ok {
					return
				}
				task.Status = model.InProgress
				log.Printf("Worker %d task processing %s", id, task.Id)
				time.Sleep(5 * time.Second)
				task.Status = model.Done
				log.Printf("Worker %d finished%s", id, task.Id)

			case <-stop:
				log.Printf("Worker %d stop", id)
				return
			}
		}
	}()
}
