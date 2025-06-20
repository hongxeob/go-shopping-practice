package db

import (
	"fmt"
	"time"
)

type Config struct {
	Primary Info `yaml:"primary"`
	Replica Info `yaml:"replica"`
}
type Info struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DbName   string `yaml:"dbname"`
	SslMode  string `yaml:"ssl-mode"`
	Pool     Pool   `yaml:"pool"`
}

// Info 구조체는 필드를 기반으로 PostgreSQL 연결 문자열을 형식화하여 반환합니다.
func (c Info) connString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.DbName,
		c.SslMode,
	)
}

type Pool struct {
	MaxConns          int32         `yaml:"max-conns"`           // 최대 연결 수 (기본값: 4 or CPU)
	MinConns          int32         `yaml:"min-conns"`           // 최소 유지할 연결 수 (기본값: 0)
	MaxConnLifetime   time.Duration `yaml:"max-conn-lifetime"`   // 연결 최대 수명 (기본값: 1시간)
	MaxConnIdleTime   time.Duration `yaml:"max-conn-idle-time"`  // 최대 유휴 시간 (기본값: 30분)
	HealthCheckPeriod time.Duration `yaml:"health-check-period"` // 헬스체크 주기 (기본값: 1분)
	ConnectTimeout    time.Duration `yaml:"connect-timeout"`     // 연결 타임아웃
}
