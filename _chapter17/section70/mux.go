package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
	"go_todo_app/handler"
	"go_todo_app/store"
	"net/http"
)

// 어떤 처리를 어떤 URL 패스로 공개할지 라우팅하는 NewMux 함수 구현
func NewMux() http.Handler {
	mux := chi.NewRouter() // 새로운 라우터 인스턴스 생성

	// /health 경로에 대한 핸들러를 추가하여 애플리케이션 상태를 확인
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write([]byte(`{"status": "ok"}`)) // JSON 형식의 상태 응답을 반환
	})

	// 입력값 검증을 위한 validator 인스턴스 생성
	v := validator.New()
	// AddTask 핸들러를 /tasks 경로에 추가
	mux.Handle("/tasks", &handler.AddTask{Store: store.Tasks, Validator: v})
	// AddTask 핸들러 인스턴스를 생성하여 /tasks 경로에 POST 요청 처리
	at := &handler.AddTask{Store: store.Tasks, Validator: v}
	// /tasks 경로에 POST 요청을 처리하기 위한 AddTask 핸들러 추가
	mux.Post("/tasks", at.ServeHTTP)
	// ListTask 핸들러 인스턴스를 생성하여 /tasks 경로에 GET 요청 처리
	lt := &handler.ListTask{Store: store.Tasks}
	// /tasks 경로에 GET 요청이 들어오면 lt.ServeHTTP를 호출하여 모든 Task 목록을 반환
	mux.Get("/tasks", lt.ServeHTTP)
	// 설정이 완료된 라우터 반환
	return mux
}
