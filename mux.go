package main

import (
	"context"
	"go_todo_app/service"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
	"go_todo_app/clock"
	"go_todo_app/config"
	"go_todo_app/handler"
	"go_todo_app/store"
)

// 파라미터와 리턴값 둘다 변경
func NewMux(ctx context.Context, cfg *config.Config) (http.Handler, func(), error) {
	mux := chi.NewRouter() // 새로운 라우터 인스턴스 생성

	// /health 경로에 대한 핸들러를 추가하여 애플리케이션 상태를 확인
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write([]byte(`{"status": "ok"}`)) // JSON 형식의 상태 응답을 반환
	})
	v := validator.New()
	db, cleanup, err := store.New(ctx, cfg) // 데이터베이스 연결을 위한 New 함수 호출
	if err != nil {
		return nil, cleanup, err
	}
	r := store.Repository{Clocker: clock.RealClocker{}}
	at := &handler.AddTask{
		Service:   &service.AddTask{DB: db, Repo: &r},
		Validator: v,
	}
	mux.Post("/tasks", at.ServeHTTP) // /tasks 경로에 대한 POST 요청을 AddTask 핸들러로 라우팅
	lt := &handler.ListTask{
		Service: &service.ListTask{DB: db, Repo: &r},
	}
	mux.Get("/tasks", lt.ServeHTTP) // /tasks 경로에 대한 GET 요청을 ListTask 핸들러로 라우팅

	ru := &handler.RegisterUser{
		Service:   &service.RegisterUser{DB: db, Repo: &r},
		Validator: v,
	}
	// /register 경로에 대한 POST 요청을 RegisterUser 핸들러로 라우팅
	mux.Post("/register", ru.ServeHTTP)

	return mux, cleanup, nil
}
