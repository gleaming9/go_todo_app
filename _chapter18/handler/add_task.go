package handler

import (
	"encoding/json"
	"github.com/go-playground/validator"
	"go_todo_app/entity"
	"go_todo_app/store"
	"net/http"
	"time"
)

type AddTask struct {
	Store     *store.TaskStore
	Validator *validator.Validate
}

func (at *AddTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	//요청 본문을 JSON으로 파싱하여 Title 필드를 가져옵니다.
	var b struct {
		// Title 필드는 필수입니다.
		Title string `json:"title" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		// JSON 파싱 중 오류가 발생하면 500 상태 코드와 오류 메시지를 반환
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}
	// Validator를 사용하여 Title 필드가 있는지 검증
	// Unmarshal한 타입에 대한 검증을 위해 Validator를 사용, JSON 구조가 방대하거나 복잡한 경우 자주 사용
	if err := at.Validator.Struct(b); err != nil {
		// 필수 필드가 없으면 400 상태 코드와 오류 메시지를 반환
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}
	// Task 구조체를 생성하여, Title, Status, Created 필드를 설정
	t := &entity.Task{
		Title:   b.Title,
		Status:  "todo",
		Created: time.Now(),
	}
	// TaskStore에 Task를 추가하고 ID를 반환
	id, err := store.Tasks.Add(t)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	// Task의 ID를 JSON 응답으로 반환
	rsp := struct {
		ID entity.TaskID `json:"id"` // JSON 응답에서 ID 필드의 이름을 "id"로 설정
	}{ID: id}
	RespondJSON(ctx, w, rsp, http.StatusOK)
}
