.PHONY: help build build-local up down logs ps test
.DEFAULT_GOAL := help

DOCKER_TAG := latest
build: ## Build the docker image
	docker build -t gleaming9/go_todo_app:${DOCKER_TAG} \
 		--target deploy ./

build-local: ## Build the docker image for local development
	docker compose build --no-cache

up: ## 자동 새로고침을 사용한 도커 컴포즈 실행
	docker compose up -d

down: ## 도커 컴포즈 종료
	docker compose down

logs: ## 도커 컴포즈 로그 출력
	docker compose logs -f

ps: ## 실행중인 컨테이너 확인
	docker compose ps

test: ## 테스트 실행
	go test -race -shuffle=on ./...

help: ## 옵션 보기
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
    		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'