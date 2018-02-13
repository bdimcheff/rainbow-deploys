.PHONY: build
build: rainbow-deploys

rainbow-deploys: main.go
	go build

image:
	docker build .