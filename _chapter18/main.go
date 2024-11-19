package main

import (
	"context"
	"fmt"
	"go_todo_app/config"
	"log"
	"net"
	"os"
)

func run(ctx context.Context) error {
	cfg, err := config.New()
	if err != nil {
		return err
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen port %d: %v", cfg.Port, err)
	}
	url := fmt.Sprintf("http://%s", l.Addr().String())
	log.Printf("start with: %v", url)

	mux := NewMux()        // NewMux 함수를 사용하여 HTTP 핸들러를 생성한다.
	s := NewServer(l, mux) // NewServer 함수를 사용하여 서버를 생성한다.
	return s.Run(ctx)      // Run 메서드를 사용하여 서버를 실행한다.
}

func main() {
	if err := run(context.Background()); err != nil {
		log.Printf("failed to terminate server: %v", err)
		os.Exit(1)
	}
}
