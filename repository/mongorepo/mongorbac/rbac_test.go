package mongorbac

import (
	"context"
	"mdhesari/kian-quiz-golang-game/config"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/repository/mongorepo"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCanCreateAndGetRoles(t *testing.T) {
	repo := getTestRepo()

	repo.CreateRole(context.TODO(), entity.Role{
		Name: "test",
	})

	roles, err := repo.GetRoles(context.TODO())
	if err != nil {

		t.Error(err.Error())
	}
	
	if len(roles) < 1 {
		t.Error("No roles.")
	}

	if roles[0].Name != "test" {
		t.Error("Invalid role name.")
	}
}

func TestCanGetPermissionIds(t *testing.T) {
	repo := getTestRepo()

	repo.HasPermissions(context.Background(), primitive.NewObjectID())
}

func getTestRepo() *DB {
	cfg := config.Load("../../../config.test.yml")

	mongoRepo, err := mongorepo.New(cfg.Database.MongoDB)
	if err != nil {
		panic(err)
	}

	return New(mongoRepo)
}