package gamevalidator

import (
	"errors"
	"mdhesari/kian-quiz-golang-game/param"
	"mdhesari/kian-quiz-golang-game/pkg/errmsg"
	"mdhesari/kian-quiz-golang-game/pkg/richerror"

	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ValidationBag map[string]string

type Validator struct {
	//
}

func New() Validator {
	return Validator{}
}

func (v *Validator) getValidationBag(err error) ValidationBag {
	errFields := make(ValidationBag)

	errV, ok := err.(validation.Errors)
	if ok {
		for k, v := range errV {
			errFields[k] = v.Error()
		}
	}

	return errFields
}

func (v *Validator) ValidateAnswerQuestion(req param.GameAnswerQuestionRequest) (ValidationBag, error) {
	const op = "gamevalidator.ValidateAnswerQuestion"

	err := validation.ValidateStruct(&req,
		validation.Field(&req.GameId, validation.Required, validation.By(v.isPrimitiveValid)),
		validation.Field(&req.QuestionId, validation.Required, validation.By(v.isPrimitiveValid)),
		validation.Field(&req.Answer, validation.Required),
	)

	if err != nil {
		errFields := v.getValidationBag(err)

		return errFields, richerror.New(op, errmsg.ErrInvalidInput).
			WithKind(richerror.KindInvalid).
			WithMeta(map[string]interface{}{
				"req": req,
			}).WithErr(err)
	}

	return nil, nil
}

func (v *Validator) isPrimitiveValid(value interface{}) error {
	id, _ := value.(primitive.ObjectID)
	if id.IsZero() {

		return errors.New(errmsg.ErrInvalidInput)
	}

	return nil
}
