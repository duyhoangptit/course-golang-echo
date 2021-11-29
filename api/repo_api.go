package api

import (
	req2 "course-golang/domain/req"
	"course-golang/domain/res"
	"course-golang/log"
	"course-golang/model"
	"course-golang/repository"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

type RepoApi struct {
	GithubRepo repository.GithubRepo
}

// RepoTrending godoc
// @Summary Get list repo trending
// @Tags github
// @Accept json
// @Produce json
// @Security jwt
// @Success 200 {object} res.Response
// @Failure 400 {object} res.Response
// @Failure 500 {object} res.Response
// @Router /github/trending [get]
func (r RepoApi) RepoTrending(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*model.JwtCustomClaims)

	repos, _ := r.GithubRepo.SelectRepos(c.Request().Context(), claims.UserId, 25)
	for i, repo := range repos {
		repos[i].Contributors = strings.Split(repo.BuildBy, ",")
	}

	return c.JSON(http.StatusOK, res.Response{
		StatusCode: http.StatusOK,
		Message:    "Xu ly thanh cong",
		Data:       repos,
	})
}

// SelectBookmarks godoc
// @Summary Get list bookmark
// @Tags github
// @Accept json
// @Produce json
// @Security jwt
// @Success 200 {object} res.Response
// @Failure 400 {object} res.Response
// @Failure 500 {object} res.Response
// @Router /bookmark/list [get]
func (r RepoApi) SelectBookmarks(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*model.JwtCustomClaims)

	repos, _ := r.GithubRepo.SelectAllBookmarks(c.Request().Context(), claims.UserId)
	for i, repo := range repos {
		repos[i].Contributors = strings.Split(repo.BuildBy, ",")
	}

	return c.JSON(http.StatusOK, res.Response{
		StatusCode: http.StatusOK,
		Message:    "Xử lý thành công",
		Data:       repos,
	})
}

// Bookmark godoc
// @Summary Bookmark repository
// @Tags github
// @Accept json
// @Produce json
// @Security jwt
// @Param data body req.BookmarkReq true "Thông tin repo name cần bookmark"
// @Success 200 {object} res.Response
// @Failure 400 {object} res.Response
// @Failure 500 {object} res.Response
// @Router /bookmark/add [post]
func (r RepoApi) Bookmark(c echo.Context) error {
	req := req2.BookmarkReq{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	// validate thong tin gui len
	err := c.Validate(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, res.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
	}

	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*model.JwtCustomClaims)

	bid, err := uuid.NewUUID()
	if err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusForbidden, res.Response{
			StatusCode: http.StatusForbidden,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	err = r.GithubRepo.Bookmark(
		c.Request().Context(),
		bid.String(),
		req.RepoName,
		claims.UserId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, res.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	return c.JSON(http.StatusOK, res.Response{
		StatusCode: http.StatusOK,
		Message:    "Bookmark thành công",
		Data:       nil,
	})
}

// DeleteBookmark godoc
// @Summary Delete bookmark repository
// @Tags github
// @Accept json
// @Produce json
// @Security jwt
// @Param data body req.BookmarkReq true "Thông tin bookmark cần unbookmark"
// @Success 200 {object} res.Response
// @Failure 400 {object} res.Response
// @Failure 500 {object} res.Response
// @Router /bookmark/delete [delete]
func (r RepoApi) DeleteBookmark(c echo.Context) error {
	req := req2.BookmarkReq{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	// validate thong tin gui len
	err := c.Validate(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, res.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Yeu cau khong dung dinh dang",
		})
	}

	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*model.JwtCustomClaims)

	err = r.GithubRepo.DelBookmark(c.Request().Context(), req.RepoName, claims.UserId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, res.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	return c.JSON(http.StatusOK, res.Response{
		StatusCode: http.StatusOK,
		Message:    "Xoá bookmark thành công",
		Data:       nil,
	})
}
