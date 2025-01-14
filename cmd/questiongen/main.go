package main

import (
	// "encoding/json"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"mdhesari/kian-quiz-golang-game/config"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/logger"
	"mdhesari/kian-quiz-golang-game/param"
	"mdhesari/kian-quiz-golang-game/repository/mongorepo"
	"mdhesari/kian-quiz-golang-game/repository/mongorepo/mongocategory"
	"mdhesari/kian-quiz-golang-game/service/categoryservice"
	"net/http"
	"sync"
	"time"

	"github.com/avast/retry-go/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

var cfg config.Config

var amountOfQuestions = 40
var questions chan entity.Question
var difficulties map[string]entity.QuestionDifficulty = map[string]entity.QuestionDifficulty{
	"easy":   entity.QuestionDifficultyEasy,
	"medium": entity.QuestionDifficultyMedium,
	"hard":   entity.QuestionDifficultyHard,
}

type Question struct {
	Title         string   `json:"question"`
	CorrectAnswer string   `json:"correct_answer"`
	OtherAnswers  []string `json:"incorrect_answers"`
	Category      string   `json:"category"`
	Difficulty    string   `json:"difficulty"`
}

type OpenTDB struct {
	Results []Question `json:"results"`
}

func init() {
	cfg = config.Load("config.yml")
	flag.Parse()
}

func main() {
	mongoCli := getMongoCli()

	categoriesOpenTDB := getCategoriesOpenTdbMap(mongoCli)

	questions = make(chan entity.Question, len(categoriesOpenTDB)*amountOfQuestions)

	logger.L().Info("Started generating questions from opentdb.", zap.Int("amount", amountOfQuestions), zap.Int("categories_count", len(categoriesOpenTDB)))

	var wg sync.WaitGroup

	for catId, openTtdbId := range categoriesOpenTDB {
		wg.Add(1)

		go func() {
			defer wg.Done()

			err := retry.Do(func() error {
				err := findAndAddQuestions(catId, openTtdbId, questions)
				if err != nil {
					return err
				}

				return nil
			}, retry.Attempts(uint(len(categoriesOpenTDB))), retry.Delay(5*time.Second))

			if err != nil {
				logger.L().Error("Finally Could not find and add questions", zap.Error(err), zap.Any("category", catId), zap.Any("opentdbId", openTtdbId))
			}
		}()
	}

	wg.Wait()

	// in order to range
	close(questions)

	qus := []interface{}{}
	for q := range questions {
		qus = append(qus, q)
	}

	ctx, cancel := context.WithTimeout(context.Background(), mongoCli.QueryTimeout)
	defer cancel()
	rows, err := mongoCli.Conn().Collection("questions").InsertMany(ctx, qus)
	if err != nil {
		logger.L().Error("Could not create questions", zap.Error(err))
	}

	<-ctx.Done()

	msg := fmt.Sprintf("\n%d questions created from %d generated questions count.", len(rows.InsertedIDs), len(qus))
	logger.L().Info(msg)
}

func findAndAddQuestions(categoryId primitive.ObjectID, opentdbId int, ques chan<- entity.Question) error {
	res, err := http.Get(fmt.Sprintf("https://opentdb.com/api.php?amount=%d&category=%d&type=multiple", amountOfQuestions, opentdbId))
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		logger.L().Error(err.Error())

		return err
	}

	// TODO - should be refactored
	var results OpenTDB
	json.Unmarshal(b, &results)
	if len(results.Results) < 1 {
		return errors.New("Ratelimit.")
	}

	for _, q := range results.Results {
		answers := make([]entity.Answer, 0)
		answers = append(answers, entity.Answer{
			Title:     q.CorrectAnswer,
			IsCorrect: true,
		})
		for _, oa := range q.OtherAnswers {
			answers = append(answers, entity.Answer{
				Title:     oa,
				IsCorrect: false,
			})
		}

		ques <- entity.Question{
			Title:      q.Title,
			CategoryID: categoryId,
			Answers:    answers,
			Difficulty: difficulties[q.Difficulty],
		}
	}

	return nil
}

func getCategoriesOpenTdbMap(mongoCli *mongorepo.MongoDB) map[primitive.ObjectID]int {
	categoryRepo := mongocategory.New(mongoCli)
	categorySrv := categoryservice.New(categoryRepo)
	categoriesRes, err := categorySrv.GetAll(context.Background(), param.CategoryParam{})
	if err != nil {
		panic(err)
	}

	// map our category titles to api's ids
	openTdb := map[interface{}]int{
		"General Knowledge": 9,
		"Sports & Fitness":  21,
		"Science":           17,
		"Music":             12,
		"History":           23,
		"Technology":        18,
	}
	categoriesOpenTDB := map[primitive.ObjectID]int{}

	for _, cat := range categoriesRes.Items {
		if tdbId, ok := openTdb[cat.Title]; ok {

			categoriesOpenTDB[cat.ID] = tdbId
		}
	}

	return categoriesOpenTDB
}

func getMongoCli() *mongorepo.MongoDB {
	mongoCli := mongorepo.New(cfg.Database.MongoDB)

	return mongoCli
}
