package service

import (
	"context"
	"fmt"
	"go_todo_app/entity"
	"go_todo_app/store"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUser 구조체는 사용자 등록 서비스를 정의
type RegisterUser struct {
	DB   store.Execer // 데이터베이스 실행 객체
	Repo UserRegister // 사용자 등록을 처리하는 인터페이스
}

// RegisterUser 사용자 정보를 등록
func (r *RegisterUser) RegisterUser(ctx context.Context, name, password, role string,
) (*entity.User, error) {
	// 1. 비밀번호를 해싱
	pw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("cannot hash password: %w", err)
	}
	// 2. 사용자 정보 생성
	u := &entity.User{
		Name:     name,
		Password: string(pw),
		Role:     role,
	}
	// 3. 사용자 정보를 데이터베이스에 저장
	if err := r.Repo.RegisterUser(ctx, r.DB, u); err != nil {
		return nil, fmt.Errorf("failed to register: %w", err)
	}
	// 4. 성공적으로 등록된 사용자 정보를 반환
	return u, nil
}
