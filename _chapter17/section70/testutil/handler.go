package testutil

import (
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"io"
	"net/http"
	"os"
	"testing"
)

// 두 JSON이 동일한지 확인하는 테스트 헬퍼 함수

func AssertJSON(t *testing.T, want, got []byte) {
	t.Helper()
	var jw, jg any
	if err := json.Unmarshal(want, &jw); err != nil {
		t.Fatalf("failed to unmarshal want %q: %v", want, err)
	}
	if err := json.Unmarshal(got, &jg); err != nil {
		t.Fatalf("failed to unmarshal got %q: %v", got, err)
	}
	if diff := cmp.Diff(jw, jg); diff != "" {
		t.Errorf("got differs: (-got +want)\n%s", diff)
	}
}

// 예상 상태 코드와 JSON 응답 본문이 일치하는지 확인하는 테스트 헬퍼 함수

func AssertResponse(t *testing.T, got *http.Response, status int, body []byte) {
	t.Helper() // 테스트 헬퍼 함수임을 표시

	// 테스트 종료 후 got.Body를 자동으로 닫기 위해 Cleanup 함수를 설정
	t.Cleanup(func() { _ = got.Body.Close() })

	// 응답 바디를 읽어 gb에 저장
	gb, err := io.ReadAll(got.Body)
	if err != nil {
		t.Fatal(err)
	}
	// 실제 상태 코드가 예상 상태 코드와 다른 경우 실패
	if got.StatusCode != status {
		t.Fatalf("want status %d, but got %d, body: %q", status, got.StatusCode, gb)
	}

	// 예상 본문과 실제 본문이 모두 빈 경우
	if len(gb) == 0 && len(body) == 0 {
		// 어느 쪽이든 응답 바디가 없으므로 AssertJSON을 호출하지 않는다.
		return
	}
	// JSON 응답 본문을 비교하여 일치하는지 확인
	AssertJSON(t, body, gb)
}

func LoadFile(t *testing.T, path string) []byte {
	t.Helper()
	bt, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file %q: %v", path, err)
	}
	return bt
}
