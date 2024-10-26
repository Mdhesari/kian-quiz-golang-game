package scheduler

import (
	"fmt"
	"log"
	"mdhesari/kian-quiz-golang-game/param"
	"mdhesari/kian-quiz-golang-game/service/matchingservice"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"
)

type Scheduler struct {
	sch         gocron.Scheduler
	matchingSrv *matchingservice.Service
}

func New(matchingSrv *matchingservice.Service) Scheduler {
	sch, err := gocron.NewScheduler()
	if err != nil {
		log.Println("Schedule error")
	}

	return Scheduler{
		sch:         sch,
		matchingSrv: matchingSrv,
	}
}

func (s Scheduler) Start(wg *sync.WaitGroup) {
	fmt.Println("started")

	j, err := s.sch.NewJob(
		gocron.DurationJob(
			5*time.Second,
		),
		gocron.NewTask(s.matchWaitedUsers),
	)
	if err != nil {
		log.Println("Schedule job failed: ", err)
	}

	log.Println(j.ID())

	s.sch.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	for {
		select {
		case <-quit:
			s.sch.Shutdown()
			fmt.Println("exiting scheduller...")

			wg.Done()
			return
		default:
			fmt.Println(time.Now())

			time.Sleep(1 * time.Second)
		}
	}
}

func (s Scheduler) matchWaitedUsers() {
	log.Println("matching waited users...")

	s.matchingSrv.MatchWaitedUsers(param.MatchingWaitedUsersRequest{})
}
