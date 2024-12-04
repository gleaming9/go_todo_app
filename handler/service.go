package handler

import (
	"context"
	"go_todo_app/entity"
)

//go:generate go run github.com/matryer/moq -out moq_test.go . ListTasksService AddTaskService RegisterUserService LoginService
type ListTasksService interface {
	// ListTasks는 컨텍스트를 사용하여 Task 목록을 반환합니다.
	// 호출 시 entity.Tasks 타입의 모든 Task 목록과 오류를 반환할 수 있습니다.
	ListTasks(ctx context.Context) (entity.Tasks, error)
}

type AddTaskService interface {
	// AddTask는 주어진 제목(title)을 가진 Task를 추가합니다.
	// 컨텍스트를 사용하여 호출하며, 추가된 Task와 오류를 반환할 수 있습니다.
	AddTask(ctx context.Context, title string) (*entity.Task, error)
}

// RegisterUserService 사용자 등록 서비스를 의미합니다.
type RegisterUserService interface {
	// RegisterUser 주어진 이름(name), 비밀번호(password), 권한(role)을 사용하여 사용자를 등록합니다.
	RegisterUser(ctx context.Context, name, password, role string) (*entity.User, error)
}

type LoginService interface {
	Login(ctx context.Context, name, pw string) (string, error)
}
