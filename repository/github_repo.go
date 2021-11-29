package repository

import (
	"context"
	"go-module/model"
)

type GithubRepo interface {
	//Github
	SaveRepo(ctx context.Context, repo model.Repos) (model.Repos, error)
	SelectRepos(ctx context.Context, userId string, limit int) ([]model.Repos, error)
	SelectRepoByName(ctx context.Context, name string) (model.Repos, error)
	UpdateRepo(ctx context.Context, repo model.Repos) (model.Repos, error)

	//Bookmark
	SelectAllBookmarks(context context.Context, userId string) ([]model.Repos, error)
	Bookmark(context context.Context, bid, nameRepo, userId string) error
	DelBookmark(context context.Context, nameRepo, userId string) error
}
