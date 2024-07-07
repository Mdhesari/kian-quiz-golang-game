package matchingvalidator

import (
	"context"
	"errors"
	"fmt"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/param"
	"mdhesari/kian-quiz-golang-game/pkg/errmsg"
	"mdhesari/kian-quiz-golang-game/pkg/richerror"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ValidationBag map[string]string

type CategoryRepo interface {
	FindById(ctx context.Context, id primitive.ObjectID) (*entity.Category, error)
}

type Validator struct {
	categoryRepo CategoryRepo
}

func New(categoryRepo CategoryRepo) Validator {
	return Validator{
		categoryRepo: categoryRepo,
	}
}

func (v Validator) getValidationBag(err error) ValidationBag {
	errFields := make(ValidationBag)

	errV, ok := err.(validation.Errors)
	if ok {
		for k, v := range errV {
			errFields[k] = v.Error()
		}
	}

	return errFields
}

func (v Validator) ValidateAddToWaitingListRequest(req param.MatchingAddToWaitingListRequest) (ValidationBag, error) {
	const op = "matchingvalidator.ValidateAddToWaitingListRequest"

	// TODO - Category is not mapped correctly with koanf
	err := validation.ValidateStruct(
		&req,
		validation.Field(&req.CategoryID, validation.Required, validation.By(v.isCategoryValid)),
		validation.Field(&req.UserID, validation.Required),
	)
	if err != nil {
		errFields := v.getValidationBag(err)

		fmt.Println("fields: ", errFields)
		return errFields, richerror.New(op, errmsg.ErrInvalidInput).
			WithKind(richerror.KindInvalid).
			WithMeta(map[string]interface{}{
				// TODO: password is exposed
				"req": req,
			}).WithErr(err)
	}

	return nil, nil
}

func (v Validator) isCategoryValid(value interface{}) error {
	categoryId := value.(primitive.ObjectID)

	cat, err := v.categoryRepo.FindById(context.Background(), categoryId)
	if err != nil || cat == nil {

		return errors.New(errmsg.ErrCategoryNotFound)
	}

	return nil
}
