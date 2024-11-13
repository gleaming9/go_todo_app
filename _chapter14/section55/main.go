package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
)

func run(ctx context.Context) error {
	s := &http.Server{
		Addr: ":18080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
		}),
	}
	//errgroup.WithContext(ctx)는 오류 그룹과 새로운 context를 반환합니다.
	//이 context는 서버가 종료되거나 다른 고루틴에서 오류가 발생하면 자동으로 취소됩니다.
	eg, ctx := errgroup.WithContext(ctx)

	//errgroup의 Go 메서드를 사용하여 서버를 실행하는 고루틴을 생성합니다.
	//이 고루틴은 s.ListenAndServe()를 호출하여 서버를 시작합니다.
	eg.Go(func() error {
		if err := s.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			log.Printf("failed to close: %+v", err)
			return err
		}
		return nil
	})
	/*
		ctx.Done() 채널을 대기하여 context가 취소될 때까지 기다립니다.
		ctx.Done()이 반환되면 서버 종료를 시작합니다.
		s.Shutdown(context.Background())를 호출하여 안전하게 서버를 종료하며,
		Shutdown은 기존의 모든 연결을 안전하게 종료합니다.
	*/
	<-ctx.Done()
	if err := s.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %+v", err)
	}
	//모든 고루틴이 종료될 때까지 기다립니다.
	//만약 어떤 고루틴에서 오류가 발생했다면 그 오류를 반환합니다.
	return eg.Wait()
}

func main() {
	if err := run(context.Background()); err != nil {
		log.Printf("failed to terminate server: %v", err)
	}
}
