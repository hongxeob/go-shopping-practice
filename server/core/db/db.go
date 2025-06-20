package db

import (
	"context"
	"fmt"
	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Conn struct {
	Primary *pgxpool.Pool
	Replica *pgxpool.Pool
}

func newConnection(dbCfg Config) (*Conn, error) {
	ctx := context.Background()
	primary, err := connectDB(ctx, dbCfg.Primary)
	if err != nil {
		return nil, err
	}

	replica, err := connectDB(ctx, dbCfg.Replica)
	if err != nil {
		primary.Close()
		return nil, err
	}

	return &Conn{
		Primary: primary,
		Replica: replica,
	}, nil

}

func connectDB(ctx context.Context, cfg Info) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.connString())
	if err != nil {
		return nil, err
	}

	// OpenTelemetry 트레이서 설정
	poolConfig.ConnConfig.Tracer = otelpgx.NewTracer(
		otelpgx.WithTracerProvider(otel.GetTracerProvider()),
		otelpgx.WithTracerAttributes(
			semconv.DBSystemPostgreSQL,
			semconv.DBName(cfg.DbName),
			semconv.NetPeerName(cfg.Host),
		),
	)
	// Logger 설정

	poolConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		zap.S().Infof("Database connected: host=%s, database=%s", cfg.Host, cfg.DbName)
		return nil
	}

	poolConfig.BeforeClose = func(conn *pgx.Conn) {
		zap.S().Infof("Database connection closing: host=%s, database=%s", cfg.Host, cfg.DbName)
	}

	if cfg.Pool.MaxConns != 0 {
		poolConfig.MaxConns = cfg.Pool.MaxConns
	}
	if cfg.Pool.MinConns != 0 {
		poolConfig.MinConns = cfg.Pool.MinConns
	}
	if cfg.Pool.MaxConnLifetime != 0 {
		poolConfig.MaxConnLifetime = cfg.Pool.MaxConnLifetime
	}
	if cfg.Pool.MaxConnIdleTime != 0 {
		poolConfig.MaxConnIdleTime = cfg.Pool.MaxConnIdleTime
	}
	if cfg.Pool.HealthCheckPeriod != 0 {
		poolConfig.HealthCheckPeriod = cfg.Pool.HealthCheckPeriod
	}
	if cfg.Pool.ConnectTimeout != 0 {
		poolConfig.ConnConfig.ConnectTimeout = cfg.Pool.ConnectTimeout
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}

	if err := otelpgx.RecordStats(pool); err != nil {
		return nil, fmt.Errorf("unable to record database stats: %w", err)
	}

	return pool, nil
}

func (q *Conn) close() {
	if q.Primary != nil {
		q.Primary.Close()
	}
	if q.Replica != nil {
		q.Replica.Close()
	}
}

var Module = fx.Module("database",
	fx.Provide(newConnection),
	fx.Invoke(func(lc fx.Lifecycle, q *Conn) {
		lc.Append(fx.Hook{
			OnStop: func(context.Context) error {
				q.close()
				return nil
			},
		})
		zap.S().Info("Loaded Database")
	}),
)
