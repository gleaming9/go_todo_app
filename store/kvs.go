package store

import (
	"context"
	"fmt"
	"go_todo_app/entity"
	"time"

	"github.com/go-redis/redis/v8"
	"go_todo_app/config"
)

// KVS : Redis 를 사용하는 KVS(Key-Value Store) 구조체 정의
type KVS struct {
	Cli *redis.Client
}

// NewKVS : KVS 구조체 생성자
func NewKVS(ctx context.Context, cfg *config.Config) (*KVS, error) {
	// Redis 클라이언트 생성
	cli := redis.NewClient(&redis.Options{
		// 서버 주소와 포트를 cfg 에서 가져와 설정
		Addr: fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
	})

	// Redis 서버와 연결을 확인. 연결에 실패하면 에러를 반환
	if err := cli.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &KVS{Cli: cli}, nil
}

// Save : 키-값을 저장하는 메서드
func (k *KVS) Save(ctx context.Context, key string, userID entity.UserID) error {
	id := int64(userID)
	// 데이터는 30분 후 만료
	return k.Cli.Set(ctx, key, id, 30*time.Minute).Err()
}

// Load : 키에 해당하는 값을 불러오는 메서드
func (k *KVS) Load(ctx context.Context, key string) (entity.UserID, error) {
	id, err := k.Cli.Get(ctx, key).Int64()
	if err != nil {
		return 0, fmt.Errorf("failed to get by %q: %w", key, ErrNotFound)
	}
	return entity.UserID(id), nil
}
