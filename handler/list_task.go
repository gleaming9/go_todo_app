package handler

import (
	"go_todo_app/entity"
	"net/http"
)

// 모든 Task를 가져오는 핸들러
type ListTask struct {
	Service ListTasksService
}

// task 구조체는 JSON 데이터 형식을 정의
type task struct {
	ID     entity.TaskID     `json:"id"`
	Title  string            `json:"title"`
	Status entity.TaskStatus `json:"status"`
}

func (lt *ListTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//요청 컨텍스트를 가져옴
	ctx := r.Context()

	tasks, err := lt.Service.ListTasks(ctx)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	// Task를 JSON으로 변환
	rsp := []task{}

	// 저장된 각 태스크를 순회하며 rsp에 추가
	for _, t := range tasks {
		rsp = append(rsp, task{
			ID:     t.ID,
			Title:  t.Title,
			Status: t.Status,
		})
	}

	// JSON 형식으로 응답 반환
	RespondJSON(ctx, w, rsp, http.StatusOK)
}
