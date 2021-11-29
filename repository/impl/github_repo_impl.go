package repository

import (
	"context"
	"course-golang/config/database"
	"course-golang/constant"
	"course-golang/log"
	"course-golang/model"
	"database/sql"
	"github.com/lib/pq"
	"time"
)

type GithubRepoImpl struct {
	sql *database.PostgreSql
}

func NewGithubRepo(sql *database.PostgreSql) *GithubRepoImpl {
	return &GithubRepoImpl{
		sql: sql,
	}
}

func (g GithubRepoImpl) SaveRepo(ctx context.Context, repo model.Repos) (model.Repos, error) {
	statement := `INSERT INTO repos(
					name, description, url, color, lang, fork, stars, 
 			        stars_today, build_by, created_at, updated_at) 
          		  VALUES(
					:name,:description, :url, :color, :lang, :fork, :stars, 
					:stars_today, :build_by, :created_at, :updated_at
				  )`

	repo.CreatedAt = time.Now()
	repo.UpdatedAt = time.Now()

	_, err := g.sql.Db.NamedExecContext(ctx, statement, repo)
	if err != nil {
		log.Error(err.Error())
		if err, ok := err.(*pq.Error); ok {
			if err.Code.Name() == "unique_violation" {
				return repo, constant.UserConflict
			}
		}
		return repo, constant.SignUpFail
	}

	return repo, nil
}
func (g GithubRepoImpl) SelectRepos(ctx context.Context, userId string, limit int) ([]model.Repos, error) {
	var repos []model.Repos
	err := g.sql.Db.SelectContext(ctx, &repos,
		`
			SELECT 
				repos.name, repos.description, repos.url, repos.color, repos.lang, 
				repos.fork, repos.stars, repos.stars_today, repos.build_by, repos.updated_at, 
				COALESCE(repos.name = bookmarks.repo_name, FALSE) as bookmarked
			FROM repos
			FULL OUTER JOIN bookmarks 
			ON repos.name = bookmarks.repo_name AND 
			   bookmarks.user_id=$1  
			WHERE repos.name IS NOT NULL 
			ORDER BY updated_at ASC LIMIT $2
		`, userId, limit)

	if err != nil {
		log.Error(err.Error())
		return repos, err
	}

	return repos, nil
}
func (g GithubRepoImpl) SelectRepoByName(ctx context.Context, name string) (model.Repos, error) {
	var repo = model.Repos{}
	err := g.sql.Db.GetContext(ctx, &repo, `SELECT * FROM repos WHERE name=$1`, name)

	if err != nil {
		if err == sql.ErrNoRows {
			return repo, constant.RepoNotFound
		}
		log.Error(err.Error())
		return repo, err
	}

	return repo, nil
}
func (g GithubRepoImpl) UpdateRepo(ctx context.Context, repo model.Repos) (model.Repos, error) {
	// name, description, url, color, lang, fork, stars, stars_today, build_by, created_at, updated_at
	sqlStatement := `
		UPDATE repos
		SET 
			stars  = :stars,
			fork = :fork,
			stars_today = :stars_today,
			build_by = :build_by,
			updated_at = :updated_at
		WHERE name = :name
	`

	result, err := g.sql.Db.NamedExecContext(ctx, sqlStatement, repo)
	if err != nil {
		log.Error(err.Error())
		return repo, err
	}

	count, err := result.RowsAffected()
	if err != nil {
		log.Error(err.Error())
		return repo, constant.RepoNotUpdated
	}
	if count == 0 {
		return repo, constant.RepoNotUpdated
	}

	return repo, nil
}

func (g GithubRepoImpl) SelectAllBookmarks(context context.Context, userId string) ([]model.Repos, error) {
	repos := []model.Repos{}
	err := g.sql.Db.SelectContext(context, &repos,
		`SELECT 
					repos.name, repos.description, repos.url, 
					repos.color, repos.lang, repos.fork, repos.stars, 
					repos.stars_today, repos.build_by, true as bookmarked
				FROM bookmarks 
				INNER JOIN repos
				ON bookmarks.user_id=$1 AND repos.name = bookmarks.repo_name`, userId)

	if err != nil {
		if err == sql.ErrNoRows {
			return repos, constant.BookmarkNotFound
		}
		log.Error(err.Error())
		return repos, err
	}
	return repos, nil
}
func (g GithubRepoImpl) Bookmark(context context.Context, bid, nameRepo, userId string) error {
	statement := `INSERT INTO bookmarks(
					bid, user_id, repo_name, created_at, updated_at) 
          		  VALUES($1, $2, $3, $4, $5)`

	now := time.Now()
	_, err := g.sql.Db.ExecContext(
		context, statement, bid, userId,
		nameRepo, now, now)

	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code.Name() == "unique_violation" {
				return constant.BookmarkConflict
			}
		}
		log.Error(err.Error())
		return constant.BookmarkFail
	}

	return nil
}
func (g GithubRepoImpl) DelBookmark(context context.Context, nameRepo, userId string) error {
	result := g.sql.Db.MustExecContext(
		context,
		"DELETE FROM bookmarks WHERE repo_name = $1 AND user_id = $2",
		nameRepo, userId)

	_, err := result.RowsAffected()
	if err != nil {
		log.Error(err.Error())
		return constant.DelBookmarkFail
	}

	return nil
}
