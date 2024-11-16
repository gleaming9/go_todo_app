package config

import "github.com/caarlos0/env/v6"

// Config 구조체는 애플리케이션 설정을 위한 구조체이다.
// 환경 변수에서 값을 가져와서 필드를 초기화한다.
type Config struct {
	Env  string `env:"TODO_ENV" envDefault:"dev"`
	Port int    `env:"PORT" envDefault:"80"`
}

// New 함수는 Config 구조체를 생성하고 환경 변수를 파싱하여 필드를 초기화
func New() (*Config, error) {
	cfg := &Config{} // config 구조체 초기화
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
