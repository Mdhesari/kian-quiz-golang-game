package scheduler

import (
	"context"
	"fmt"
	"log"
	"mdhesari/kian-quiz-golang-game/logger"
	"mdhesari/kian-quiz-golang-game/param"
	"mdhesari/kian-quiz-golang-game/service/matchingservice"
	"os"
	"os/signal"
	"time"

	"github.com/go-co-op/gocron/v2"
	"go.uber.org/zap"
)

type Config struct {
	MatchWaitedUsersIntervalSeconds int `koanf:"match_waited_users_interval_seconds"`
}

type Scheduler struct {
	config      Config
	sch         gocron.Scheduler
	matchingSrv *matchingservice.Service
}

func New(config Config, matchingSrv *matchingservice.Service) Scheduler {
	sch, err := gocron.NewScheduler()
	if err != nil {
		log.Println("Schedule error")
	}

	return Scheduler{
		config:      config,
		sch:         sch,
		matchingSrv: matchingSrv,
	}
}

func (s Scheduler) Start() {
	fmt.Println("started")

	_, err := s.sch.NewJob(
		gocron.DurationJob(time.Duration(s.config.MatchWaitedUsersIntervalSeconds)*time.Second),
		gocron.NewTask(s.matchWaitedUsers),
	)
	if err != nil {
		log.Fatalf("Schedule job failed: %v\n", err)

		return
	}

	s.sch.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	for {
		select {
		case <-quit:
			err := s.sch.Shutdown()
			if err != nil {
				logger.L().Error("Could not shutdown scheduller.", zap.Error(err))

				return
			}

			logger.L().Info("Scheduler shutdown gracefully.")

			return
		default:
			time.Sleep(1 * time.Second)
		}
	}
}

func (s Scheduler) matchWaitedUsers() {
	log.Println("matching waited users...")

	s.matchingSrv.MatchWaitedUsers(context.Background(), param.MatchingWaitedUsersRequest{})
}
