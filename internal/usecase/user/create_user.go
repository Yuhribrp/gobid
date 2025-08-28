package user

import (
	"context"

	"github.com/Yuhribrp/gobid/internal/validator"
)

type CreateUserReq struct {
	UserName string `json:"user_name" binding:"required"`
	Email    string `json:"email"     binding:"required,email"`
	Password string `json:"password"  binding:"required,min=8"`
	Bio      string `json:"bio"`
}

type CreateUserRes struct {
	ID       string `json:"id"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Bio      string `json:"bio"`
}

func (req CreateUserReq) Valid(ctx context.Context) validator.Evaluator {
	eval := make(validator.Evaluator)

	eval.CheckField(validator.NotBlank(req.UserName), "user_name", "user_name must be provided")
	eval.CheckField(validator.NotBlank(req.Email), "email", "email must be provided")
	eval.CheckField(
		validator.Matches(req.Email, validator.EmailRX),
		"email",
		"email must be a valid email address",
	)
	eval.CheckField(
		validator.MinChars(string(req.Password), 8),
		"password",
		"password must be at least 8 characters long",
	)
	eval.CheckField(
		validator.MinChars(req.Bio, 10) &&
			validator.MaxChars(req.Bio, 500),
		"bio",
		"bio must be between 10 and 500 characters long",
	)
	return eval
}
