package main

import (
	"course-golang/api"
	"course-golang/config/database"
	_ "course-golang/docs"
	"course-golang/helper"
	"course-golang/log"
	repository "course-golang/repository/impl"
	"course-golang/router"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/swaggo/echo-swagger"
	"os"
	"time"
)

func setupLog() {
	fmt.Println("Development enviroment")
	// log
	fmt.Println(">>>>", os.Getenv("APP_NAME"))
	//os.Setenv("APP_NAME", "tiger")
	// log ra file
	log.InitLogger(false)
}

func setupDB() *database.PostgreSql {
	fmt.Println("Setup DB")

	// setup db
	sql := &database.PostgreSql{
		Host:     "host.docker.internal",
		Port:     5432,
		UserName: "postgres",
		Password: "root",
		DBName:   "postgres",
	}

	// connect db
	sql.Connect()

	return sql
}

// @title Github Trending API
// @version 1.0
// @description More
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey jwt
// @in header
// @name Authorization

// @host localhost:8080
// @BasePath /
func main() {
	setupLog()
	sql := setupDB()

	// Khi ham main ket thuc se thuc thi close
	defer sql.Close()

	e := echo.New()

	// swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	structValidator := helper.NewStructValidator()
	structValidator.RegisterValidate()

	// register validator
	e.Validator = structValidator

	// register api
	userApi := api.UserApi{
		UserRepo: repository.NewUserRepo(sql),
	}

	repoApi := api.RepoApi{
		GithubRepo: repository.NewGithubRepo(sql),
	}

	api := router.API{
		Echo:    e,
		UserApi: userApi,
		RepoApi: repoApi,
	}

	// setup router
	api.SetupRouter()

	// crawl data from github trending
	go scheduleUpdateTrending(15*time.Second, repoApi)

	// start server license port
	e.Logger.Fatal(e.Start(":8080"))
}

func scheduleUpdateTrending(timeSchedule time.Duration, repoApi api.RepoApi) {
	ticker := time.NewTicker(timeSchedule)
	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println("Checking from github...")
				helper.CrawlRepo(repoApi.GithubRepo)
			}
		}
	}()

}
