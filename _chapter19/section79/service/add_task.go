package service

import (
	"context"
	"fmt"
	"go_todo_app/entity"
	"go_todo_app/store"
)

// AddTask 구조체는 Task를 추가하는 기능을 정의합니다.
type AddTask struct {
	DB   store.Execer
	Repo TaskAdder
}

// AddTask 메서드는 주어진 제목을 가진 Task를 추가합니다.
func (a *AddTask) AddTask(ctx context.Context, title string) (*entity.Task, error) {
	// 새로운 Task를 생성합니다.
	t := &entity.Task{
		Title:  title,
		Status: entity.TaskStatusTodo,
	}
	// Task를 데이터베이스에 추가합니다.
	err := a.Repo.AddTask(ctx, a.DB, t)
	if err != nil {
		return nil, fmt.Errorf("failed to register: %w", err)
	}
	// 추가된 Task를 반환합니다.
	return t, nil
}
