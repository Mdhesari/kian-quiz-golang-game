package uservalidator

import (
	"errors"
	"mdhesari/kian-quiz-golang-game/param"
	"mdhesari/kian-quiz-golang-game/pkg/errmsg"
	"regexp"

	"github.com/go-ozzo/ozzo-validation/v4"
)

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

func (v Validator) ValidateRegisterRequest(req param.RegisterRequest) error {
	err := validation.ValidateStruct(&req,
		validation.Field(&req.Name, validation.Required),
		validation.Field(&req.Email, validation.Required, validation.By(v.checkEmailUniqueness)),
		validation.Field(&req.Mobile,
			validation.Required,
			validation.Match(regexp.MustCompile("^09[0-9]{9}$")),
			validation.By(v.checkPhoneNumberUniqueness),
		),
		validation.Field(&req.Password, validation.Match(regexp.MustCompile(`^[A-Za-z0-9!@#%^&*]{8,}$`))),
	)

	return err
}

func (v Validator) ValidateLoginRequest(req param.LoginRequest) error {
	err := validation.ValidateStruct(&req,
		validation.Field(&req.Email, validation.Required),
	)

	return err
}
