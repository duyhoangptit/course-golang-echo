package router

import (
	"github.com/labstack/echo/v4"
	"go-module/api"
	"go-module/middleware"
)

type API struct {
	Echo     *echo.Echo
	UserApi  api.UserApi
	IndexApi api.IndexApi
	RepoApi  api.RepoApi
}

func (api *API) SetupRouter() {
	// login, register
	api.Echo.POST("/user/sign-in", api.UserApi.SignIn)
	api.Echo.POST("/user/sign-up", api.UserApi.SignUp)

	// group user
	user := api.Echo.Group("/user", middleware.JWTMiddleware())
	user.GET("/profile", api.UserApi.Profile)

	// github repo
	github := api.Echo.Group("/github", middleware.JWTMiddleware())
	github.GET("/trending", api.RepoApi.RepoTrending)

	// bookmark
	bookmark := api.Echo.Group("/bookmark", middleware.JWTMiddleware())
	bookmark.GET("/list", api.RepoApi.SelectBookmarks)
	bookmark.POST("/add", api.RepoApi.Bookmark)
	bookmark.DELETE("/delete", api.RepoApi.DeleteBookmark)
}
