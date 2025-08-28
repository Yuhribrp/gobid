package user

import (
	"context"

	"github.com/Yuhribrp/gobid/internal/validator"
)

type LoginUserReq struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginUserRes struct {
	ID       string `json:"id"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Bio      string `json:"bio"`
}

func (req LoginUserReq) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator

	eval.CheckField(validator.NotBlank(req.Email), "email", "email must be provided")
	eval.CheckField(
		validator.Matches(req.Email, validator.EmailRX),
		"email",
		"email must be a valid email address",
	)
	eval.CheckField(validator.NotBlank(req.Password), "password", "password must be provided")

	return eval
}
