package repository

import (
	"context"
	"database/sql"
	"github.com/lib/pq"
	"go-module/config/database"
	"go-module/config/security"
	"go-module/constant"
	"go-module/domain/req"
	"go-module/log"
	"go-module/model"
	"time"
)

type UserRepoImpl struct {
	sql *database.PostgreSql
}

func NewUserRepo(sql *database.PostgreSql) *UserRepoImpl {
	return &UserRepoImpl{
		sql: sql,
	}
}

func (u UserRepoImpl) SaveUser(ctx context.Context, user model.User) (model.User, error) {
	statement := `
		INSERT INTO users(user_id, email, password, role, full_name, created_at, updated_at)
		VALUES(:user_id, :email, :password, :role, :full_name, :created_at, :updated_at)
	`

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := u.sql.Db.NamedExecContext(ctx, statement, user)
	if err != nil {
		log.Error(err.Error())
		if err, ok := err.(*pq.Error); ok {
			if err.Code.Name() == "unique_violation" {
				return user, constant.UserConflict
			}
		}
		return user, constant.SignUpFail
	}

	return user, nil
}

func (u UserRepoImpl) CheckLogin(ctx context.Context, loginReq req.SignInReq) (model.User, error) {
	var user = model.User{}
	err := u.sql.Db.GetContext(ctx, &user, "select * from users where email=$1", loginReq.Email)

	if err != nil {
		log.Error(err.Error())
		if err == sql.ErrNoRows {
			return user, constant.UserNotFound
		}
		return user, err
	}

	// check password
	if !security.ComparePasswords(user.Password, []byte(loginReq.Password)) {
		return user, constant.PasswordInvalid
	}

	return user, nil
}

func (u UserRepoImpl) SelectUserById(ctx context.Context, userId string) (model.User, error) {
	var user = model.User{}
	err := u.sql.Db.GetContext(ctx, &user, "select * from users where user_id=$1", userId)

	if err != nil {
		log.Error(err.Error())
		if err == sql.ErrNoRows {
			return user, constant.UserNotFound
		}
		return user, err
	}

	return user, nil
}
