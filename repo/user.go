package repo

import (
	"context"
	"randomusergen/domain"
)

type UserRepo interface {
	SaveAll(ctx context.Context, users []*domain.UserData) (int, error)
	Close() error
}
