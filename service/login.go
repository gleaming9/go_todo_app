package service

import (
	"context"
	"fmt"
	"go_todo_app/store"
)

type Login struct {
	DB             store.Queryer
	Repo           UserGetter
	TokenGenerator TokenGenerator
}

func (l *Login) Login(ctx context.Context, name, pw string) (string, error) {
	u, err := l.Repo.GetUser(ctx, l.DB, name) // User 정보 조회
	if err != nil {
		return "", fmt.Errorf("failed to list: %w", err)
	}

	if err := u.ComparePassword(pw); err != nil { // 비밀번호 비교
		return "", fmt.Errorf("wrong password: %w", err)
	}

	jwt, err := l.TokenGenerator.GenerateToken(ctx, *u) // JWT 생성
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT: %w", err)
	}

	return string(jwt), nil
}
