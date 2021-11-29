pro:
	docker rmi -f course-golang:1.0
	docker-compose up -d # background
local:
	go run main.go
dev:
	# workdir in dockerfile: WORKDIR cmd/pro
	cd cmd/dev;go run main.go
