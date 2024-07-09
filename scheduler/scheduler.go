package scheduler

import (
	"fmt"
	"os"
	"os/signal"
	"time"
)

type Scheduler struct {
	//
}

func New() Scheduler {
	return Scheduler{}
}

func (s Scheduler) Start() {
	fmt.Println("started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)


	for {
		select {
		case <-quit:
			fmt.Println("exiting scheduller...")
			return
		default:
			fmt.Println(time.Now())

			time.Sleep(1 * time.Second)
		}
	}
}
