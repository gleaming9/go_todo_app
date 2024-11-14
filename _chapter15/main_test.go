package _chapter15

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io"
	"net"
	"net/http"
	"testing"
)

func TestRun(t *testing.T) {
	// net/http 에서는 포트 번호에 0을 지정하면 사용 가능한 포트 번호를 동적으로 선택합니다.
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to listen port %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return run(ctx, l)
	})

	in := "message"
	url := fmt.Sprintf("http://%s/%s", l.Addr().String(), in)
	// 어떤 포트 번호로 리슨중인지 확인합니다.
	t.Logf("try request to %q", url)
	rsp, err := http.Get(url)

	// 이후 코드 동일
	if err != nil {
		t.Errorf("failed to get: %+v", err)
	}
	defer rsp.Body.Close()

	got, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	// HTTP 서버의 반환값을 검증합니다.
	want := fmt.Sprintf("Hello, %s!", in)
	if string(got) != want {
		t.Errorf("want %q, but got %q", want, got)
	}

	// run 함수를 종료합니다.
	cancel()
	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}
}
