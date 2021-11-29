# BE golang learning

1. go mod init course-golang

2. go build

3. go mod vendor

4. migration db
   go get -v github.com/rubenv/sql-migrate/...
   
   `sql-migrate up`
   `sql-migrate down`
5. Thu vien luu log file
   go get github.com/rifflock/lfshook
   go get github.com/lestrrat/go-file-rotatelogs
   
6. Lib validator
   `go get github.com/go-playground/validator/v10`

7. Crawl data
   `go get -u github.com/gocolly/colly/...`
8. Swagger
   https://github.com/swaggo/echo-swagger
   `go get github.com/swaggo/swag/cmd/swag`
   `swag init`
   
9. Makefile
   `make pro`
# Lưu ý
- Các method khác package muốn gọi được ở package khác thì cần phải để viết hoa