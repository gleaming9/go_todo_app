package testutil

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
	"testing"
)

/*
레디스 동작 확인용 테스트 헬퍼
환경 변수를 확인해서 깃허브 액션이라고 판단한 경우 다른 설정 정보를 이용해서 레디스에 접속한다.
*/

func OpenRedisForTest(t *testing.T) *redis.Client {
	t.Helper()

	host := "127.0.0.1"
	port := 36379
	if _, defined := os.LookupEnv("CI"); defined {
		port = 6379
	}
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: "",
		DB:       0, // default database number
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		t.Fatalf("failed to connect redis: %s", err)
	}
	return client
}
