package handler

import (
	"encoding/json"
	"github.com/go-playground/validator"
	"go_todo_app/entity"
	"net/http"
)

// RegisterUser 구조체는 사용자 등록 요청을 처리하는 핸들러를 정의
type RegisterUser struct {
	Service   RegisterUserService // 사용자 등록 처리를 위임할 서비스
	Validator *validator.Validate // 입력값 검증을 위한 validator
}

func (ru *RegisterUser) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// 요청 본문에서 사용자 등록 정보를 파싱할 구조체
	var b struct {
		Name     string `json:"name" validate:"required"`
		Password string `json:"password" validate:"required"`
		Role     string `json:"role" validate:"required"`
	}
	// 요청 본문을 JSON으로 파싱해서 구조체에 매핑
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}
	// validator를 사용하여 구조체 필드의 값이 유효한지 검사
	if err := ru.Validator.Struct(b); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}
	// 서비스 레이어에 사용자 등록 요청을 전달
	u, err := ru.Service.RegisterUser(ctx, b.Name, b.Password, b.Role)
	if err != nil {
		// 서비스 호출 중 오류 발생시 500 상태 코드와 오류 메시지를 반환
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	// 사용자 등록 성공시 200 상태 코드와 사용자 ID를 반환
	rsp := struct {
		ID entity.UserID `json:"id"`
	}{ID: u.ID}
	RespondJSON(ctx, w, rsp, http.StatusOK)
}
