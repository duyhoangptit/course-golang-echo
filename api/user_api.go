package api

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go-module/config/security"
	"go-module/constant"
	req2 "go-module/domain/req"
	"go-module/domain/res"
	"go-module/log"
	"go-module/model"
	"go-module/repository"
	"net/http"
)

type UserApi struct {
	UserRepo repository.UserRepo
}

// SignIn godoc
// @Summary Sign in to access your account
// @Tags users
// @Accept json
// @Produce json
// @Param data body req.SignInReq true "Thông tin đăng nhập"
// @Success 200 {object} res.Response
// @Failure 400 {object} res.Response
// @Failure 500 {object} res.Response
// @Router /user/sign-in [post]
func (u *UserApi) SignIn(c echo.Context) error {
	req := req2.SignInReq{}

	// binding data
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, res.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	// validate
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, res.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	// check login
	user, err := u.UserRepo.CheckLogin(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, res.Response{
			StatusCode: http.StatusUnauthorized,
			Message:    "Đăng nhập thất bại",
			Data:       nil,
		})
	}

	// gen token
	token, err := security.GenToken(user)
	if err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, res.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Server error",
			Data:       nil,
		})
	}
	user.Token = token

	return c.JSON(http.StatusOK, res.Response{
		StatusCode: http.StatusOK,
		Message:    "Đăng nhập thành công",
		Data:       user,
	})
}

// Profile godoc
// @Summary View profile of user
// @Tags users
// @Accept json
// @Produce json
// @Security jwt
// @Success 200 {object} res.Response
// @Failure 400 {object} res.Response
// @Failure 500 {object} res.Response
// @Router /user/profile [get]
func (u *UserApi) Profile(c echo.Context) error {
	tokenData := c.Get("user").(*jwt.Token)
	claims := tokenData.Claims.(*model.JwtCustomClaims)

	// get user by user id
	user, err := u.UserRepo.SelectUserById(c.Request().Context(), claims.UserId)
	if err != nil {
		log.Error(err.Error())
		if err == constant.UserNotFound {
			return c.JSON(http.StatusNotFound, res.Response{
				StatusCode: http.StatusNotFound,
				Message:    err.Error(),
				Data:       nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, res.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	return c.JSON(http.StatusOK, res.Response{
		StatusCode: http.StatusOK,
		Message:    "Đăng nhập thành công",
		Data:       user,
	})
}

// SignUp godoc
// @Summary Create new account
// @Tags users
// @Accept json
// @Produce json
// @Param data body req.SignUpReq true "Thông tin đăng ký"
// @Success 200 {object} res.Response
// @Failure 400 {object} res.Response
// @Failure 500 {object} res.Response
// @Router /user/sign-up [post]
func (u *UserApi) SignUp(c echo.Context) error {
	req := req2.SignUpReq{}

	// binding data
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, res.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	// validate
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, res.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	hash := security.HashAndSalt([]byte(req.Password))
	role := model.MEMBER.String()
	userId, err := uuid.NewUUID()
	if err != nil {
		return c.JSON(http.StatusBadRequest, res.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	user := model.User{
		UserId:   userId.String(),
		FullName: req.FullName,
		Email:    req.Email,
		Password: hash,
		Role:     role,
	}

	user, err = u.UserRepo.SaveUser(c.Request().Context(), user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, res.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	return c.JSON(http.StatusOK, res.Response{
		StatusCode: http.StatusOK,
		Message:    "Success",
		Data:       user,
	})
}
