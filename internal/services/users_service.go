package services

import (
	"context"
	"errors"

	"github.com/Yuhribrp/gobid/internal/store/pgstore"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)


var ErrDuplicatedEmailOrPassword = errors.New("user with this email or username already exists")

type UserService struct {
	pool *pgxpool.Pool
	queries *pgstore.Queries
}


func NewUserService(pool *pgxpool.Pool) UserService {
	return UserService{
		pool: pool,
		queries: pgstore.New(pool),
	}
}

func (us *UserService) CreateUser(ctx context.Context, userName, email, password, bio string) (uuid.UUID, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.Nil, err
	}

	args := pgstore.CreateUserParams{
		UserName:     userName,
		Email:        email,
		PasswordHash: hash,
		Bio:          bio,
	}

	id, err := us.queries.CreateUser(ctx, args)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return uuid.Nil, ErrDuplicatedEmailOrPassword
		}
		return uuid.Nil, err
	}

	return id.ID, nil
}

func (us *UserService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	return us.queries.DeleteUser(ctx, userID)
}

func (us *UserService) AuthenticateUser(ctx context.Context, username string, password string) (bool, error) {
	// Implement user authentication logic here
	return true, nil
}
