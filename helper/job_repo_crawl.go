package helper

import (
	"context"
	"course-golang/constant"
	"course-golang/log"
	"course-golang/model"
	"course-golang/repository"
	"fmt"
	"github.com/gocolly/colly"
	"regexp"
	"runtime"
	"strings"
	"time"
)

type RepoProcess struct {
	repo       model.Repos
	githubRepo repository.GithubRepo
}

func CrawlRepo(repository repository.GithubRepo) {
	c := colly.NewCollector()
	repos := make([]model.Repos, 0, 30)
	c.OnHTML(`article[class=Box-row]`, func(e *colly.HTMLElement) {
		var githubRepo model.Repos

		githubRepo.Name = e.ChildText("h1.h3 > a")

		n := strings.Replace(e.ChildText("h1.h3 > a"), "\n", "", -1)
		githubRepo.Name = strings.Replace(n, " ", "", -1)

		githubRepo.Description = e.ChildText("p.col-9")

		bgColor := e.ChildAttr(".repo-language-color", "style")
		re := regexp.MustCompile("#[a-zA-Z0-9_]+")
		match := re.FindStringSubmatch(bgColor)

		if len(match) > 0 {
			githubRepo.Color = match[0]
		}

		githubRepo.Url = e.ChildAttr("h1.h3 > a", "href")
		githubRepo.Lang = e.ChildText("span[itemprop=programmingLanguage]")

		e.ForEach(".mt-2 a", func(i int, el *colly.HTMLElement) {
			if i == 0 {
				githubRepo.Stars = strings.TrimSpace(el.Text)
			} else if i == 1 {
				githubRepo.Fork = strings.TrimSpace(el.Text)
			}
		})

		e.ForEach(".mt-2 .float-sm-right", func(i int, el *colly.HTMLElement) {
			githubRepo.StarsToday = strings.TrimSpace(el.Text)
		})

		var buildBy []string
		e.ForEach(".mt-2 span a img", func(i int, el *colly.HTMLElement) {
			avatarContributor := el.Attr("src")
			buildBy = append(buildBy, avatarContributor)
		})

		githubRepo.BuildBy = strings.Join(buildBy, ",")

		repos = append(repos, githubRepo)
	})

	c.OnScraped(func(r *colly.Response) {
		queue := NewJobQueue(runtime.NumCPU())
		queue.Start()
		defer queue.Stop()

		for _, repo := range repos {
			queue.Submit(&RepoProcess{
				repo:       repo,
				githubRepo: repository,
			})
		}
	})

	c.Visit("https://github.com/trending")
}

func (r RepoProcess) Process() {
	// select repo by update
	cacheRepo, err := r.githubRepo.SelectRepoByName(context.Background(), r.repo.Name)
	if err == constant.RepoNotFound {
		// không tìm thấy repo - insert repo to database
		fmt.Println("Add: ", r.repo.Name)
		_, err := r.githubRepo.SaveRepo(context.Background(), r.repo)
		if err != nil {
			log.Error(err)
		}
		return
	}

	// neu ton tai thi update
	if r.repo.Stars != cacheRepo.Stars ||
		r.repo.StarsToday != cacheRepo.StarsToday ||
		r.repo.Fork != cacheRepo.Fork {
		fmt.Println("Update: ", r.repo.Name)
		r.repo.UpdatedAt = time.Now()
		_, err := r.githubRepo.UpdateRepo(context.Background(), r.repo)
		if err != nil {
			log.Error(err)
		}
	}

}
