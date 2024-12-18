package main

import (
	// "encoding/json"
	"context"
	"encoding/json"
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

	"go.mongodb.org/mongo-driver/bson/primitive"
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

	var wg sync.WaitGroup

	for catId, openTtdbId := range categoriesOpenTDB {
		wg.Add(1)

		go func() {
			defer wg.Done()

			findAndAddQuestions(catId, openTtdbId, questions)
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
		fmt.Println("Could not create questions", err)
	}

	<-ctx.Done()

	fmt.Printf("\n%d questions created from %d generated questions count.", len(rows.InsertedIDs), len(qus))
}

func findAndAddQuestions(categoryId primitive.ObjectID, opentdbId int, ques chan<- entity.Question) {
	res, err := http.Get(fmt.Sprintf("https://opentdb.com/api.php?amount=%d&category=%d&type=multiple", amountOfQuestions, opentdbId))
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		logger.L().Error(err.Error())

		return
	}

	// TODO - should be refactored
	var results OpenTDB
	json.Unmarshal(b, &results)
	if len(results.Results) < 1 {
		fmt.Println("results are empty let's try one more time!")
		time.Sleep(2 * time.Second)

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()

			findAndAddQuestions(categoryId, opentdbId, ques)
		}()

		wg.Wait()
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
	mongoCli, err := mongorepo.New(cfg.Database.MongoDB)
	if err != nil {

		panic("could not connect to mongodb.")
	}

	return mongoCli
}
