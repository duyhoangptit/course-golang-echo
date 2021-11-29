package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
)

type PostgreSql struct {
	Db       *sqlx.DB
	Host     string
	Port     int32
	UserName string
	Password string
	DBName   string
}

func (s *PostgreSql) Connect() {
	dataSource := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		s.Host, s.Port, s.UserName, s.Password, s.DBName)

	s.Db = sqlx.MustConnect("postgres", dataSource)
	if err := s.Db.Ping(); err != nil {
		log.Error(err.Error())
		return
	}

	fmt.Println("Connect db successfully")
}

func (s *PostgreSql) Close() {
	s.Db.Close()
}
