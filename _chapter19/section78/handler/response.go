package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrResponse struct {
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

//HTTP 핸들러 내에서 귀찮은 JSON 응답 작성을 간략화한다.

func RespondJSON(ctx context.Context, w http.ResponseWriter, body any, status int) {
	// 응답 헤더에 JSON 형식임을 지정
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("encode response error: %v", err)
		// 서버 오류(500) 상태 코드를 반환합니다.
		w.WriteHeader(http.StatusInternalServerError)

		// ErrResponse 구조체를 사용해서 응답을 반환
		rsp := ErrResponse{
			Message: http.StatusText(http.StatusInternalServerError),
		}
		// 오류 응답을 json 형태로 인코딩하여 클라이언트한테 전송
		if err := json.NewEncoder(w).Encode(rsp); err != nil {
			// 인코딩마저 실패하면 로그를 출력
			fmt.Printf("write error response error: %v", err)
		}
		return
	}

	w.WriteHeader(status)
	if _, err := fmt.Fprintf(w, "%s", bodyBytes); err != nil {
		fmt.Printf("write response error: %v", err)
	}
}
