package handler

import (
	"encoding/json"
	"github.com/go-playground/validator"
	"log"
	"net/http"
)

type Login struct {
	Service   LoginService
	Validator *validator.Validate
}

// 입출력 JSON 을 구성하는 역할만 수행

func (l *Login) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Login ServeHTTP\n")
	ctx := r.Context()
	var body struct {
		UserName string `json:"user_name" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	err := l.Validator.Struct(body)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	jwt, err := l.Service.Login(ctx, body.UserName, body.Password)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	rsp := struct {
		AccessToken string `json:"access_token"`
	}{
		AccessToken: jwt,
	}

	RespondJSON(ctx, w, rsp, http.StatusOK)
}
