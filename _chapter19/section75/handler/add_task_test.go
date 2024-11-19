package handler

import (
	"bytes"
	"github.com/go-playground/validator"
	"go_todo_app/entity"
	"go_todo_app/store"
	"go_todo_app/testutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddTask(t *testing.T) {
	type want struct {
		status  int    // 예상되는 HTTP 상태 코드
		rspFile string // 예상되는 응답 JSON 파일 경로
	}
	//테스트 케이스 정의
	tests := map[string]struct {
		reqFile string // 요청 JSON 파일 경로
		want    want   // 예상되는 응답
	}{ // 익명 구조체 정의 부분
		"ok": { // 정상적인 요청에 대한 테스트 케이스
			reqFile: "testdata/add_task/ok_req.json.golden",
			want: want{
				status:  http.StatusOK,
				rspFile: "testdata/add_task/ok_rsp.json.golden",
			},
		},
		"badRequest": { // 잘못된 요청에 대한 테스트 케이스
			reqFile: "testdata/add_task/bad_req.json.golden",
			want: want{
				status:  http.StatusBadRequest,
				rspFile: "testdata/add_task/bad_rsp.json.golden",
			},
		},
	}
	// 각 테스트 케이스에 대한 반복 실행
	for n, tt := range tests {
		tt := tt // tt 변수를 로컬로 복사하여 병렬 실행 시 데이터 경합 방지
		t.Run(n, func(t *testing.T) {
			t.Parallel() // 각 테스트 케이스를 병렬로 실행

			//가짜 응답 기록기
			w := httptest.NewRecorder()
			//가짜 요청 생성
			r := httptest.NewRequest(
				http.MethodPost,
				"/tasks",
				bytes.NewReader(testutil.LoadFile(t, tt.reqFile)),
			)

			sut := AddTask{
				Store: &store.TaskStore{
					Tasks: map[entity.TaskID]*entity.Task{}, // 빈 TaskStore를 사용하여 저장
				}, Validator: validator.New()} // 입력 검증을 위한 validator 생성

			// AddTask 핸들러 실행
			sut.ServeHTTP(w, r)

			// 예상 응답과 비교
			resp := w.Result()
			testutil.AssertResponse(t,
				resp, tt.want.status, testutil.LoadFile(t, tt.want.rspFile))
		})
	}
}
