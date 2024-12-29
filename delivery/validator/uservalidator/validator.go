package uservalidator

import (
	"errors"
	"fmt"
	"mdhesari/kian-quiz-golang-game/param"
	"mdhesari/kian-quiz-golang-game/pkg/errmsg"
	"mdhesari/kian-quiz-golang-game/pkg/richerror"
	"regexp"

	"github.com/go-ozzo/ozzo-validation/v4"
)

type ValidationBag map[string]string

type Repository interface {
	IsMobileUnique(mobile string) (bool, error)
	IsEmailUnique(email string) (bool, error)
}

type Validator struct {
	repo Repository
}

func New(repo Repository) Validator {
	return Validator{
		repo: repo,
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

func (v Validator) checkPhoneNumberUniqueness(value interface{}) error {
	mobile := value.(string)

	isUnique, err := v.repo.IsMobileUnique(mobile)
	if err != nil {

		return err
	}

	if isUnique {

		return nil
	}

	return errors.New(errmsg.ErrMobileUnique)
}

func (v Validator) checkEmailUniqueness(value interface{}) error {
	email := value.(string)

	isUnique, err := v.repo.IsEmailUnique(email)
	if err != nil {

		return err
	}

	if isUnique {

		return nil
	}

	return errors.New(errmsg.ErrEmailUnique)
}

func (v Validator) ValidateRegisterRequest(req param.RegisterRequest) (ValidationBag, error) {
	const op = "uservalidator.ValidateRegisterRequest"

	err := validation.ValidateStruct(&req,
		validation.Field(&req.Name, validation.Required),
		validation.Field(&req.Email, validation.Required, validation.By(v.checkEmailUniqueness)),
		validation.Field(&req.Mobile, validation.By(v.checkPhoneNumberUniqueness)),
		validation.Field(&req.Password, validation.Match(regexp.MustCompile(`^[A-Za-z0-9!@#%^&*]{8,}$`)).Error(errmsg.ErrWeakPassword)),
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

func (v Validator) ValidateLoginRequest(req param.LoginRequest) (ValidationBag, error) {
	const op = "uservalidator.ValidateLoginRequest"

	err := validation.ValidateStruct(&req,
		validation.Field(&req.Email, validation.Required),
	)

	if err != nil {
		errFields := v.getValidationBag(err)

		return errFields, richerror.New(op, errmsg.ErrInvalidInput).
			WithKind(richerror.KindInvalid).
			WithMeta(map[string]interface{}{
				// TODO: password is exposed
				"req": req,
			}).WithErr(err)
	}

	return nil, nil
}
