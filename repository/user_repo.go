package repository

import (
	"context"
	"go-module/domain/req"
	"go-module/model"
)

type UserRepo interface {
	CheckLogin(ctx context.Context, loginReq req.SignInReq) (model.User, error)

	SaveUser(ctx context.Context, user model.User) (model.User, error)

	SelectUserById(ctx context.Context, userId string) (model.User, error)
}
