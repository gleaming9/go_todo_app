package service

import (
	"context"
	"go_todo_app/entity"
	"go_todo_app/store"
)

// TaskAdder 인터페이스는 Task를 추가하는 메서드를 정의합니다.

//go:generate go run github.com/matryer/moq -out moq_test.go . TaskAdder TaskLister UserRegister UserGetter TokenGenerator
type TaskAdder interface {
	// AddTask는 주어진 Task를 데이터베이스에 추가합니다.
	// 컨텍스트와 데이터베이스 실행기, Task를 인자로 받습니다.
	AddTask(ctx context.Context, db store.Execer, t *entity.Task) error
}

// TaskLister 인터페이스는 Task 목록을 가져오는 메서드를 정의합니다.
type TaskLister interface {
	// ListTasks는 데이터베이스에서 Task 목록을 가져옵니다.
	// 컨텍스트와 데이터베이스 Queryer를 인자로 받습니다.
	ListTasks(ctx context.Context, db store.Queryer) (entity.Tasks, error)
}

// UserRegister 인터페이스는 사용자 등록 메서드를 정의합니다.
type UserRegister interface {
	// RegisterUser 주어진 사용자 정보를 데이터베이스에 등록합니다.
	RegisterUser(ctx context.Context, db store.Execer, u *entity.User) error
}

type UserGetter interface {
	GetUser(ctx context.Context, db store.Queryer, name string) (*entity.User, error)
}

type TokenGenerator interface {
	GenerateToken(ctx context.Context, u entity.User) ([]byte, error)
}
