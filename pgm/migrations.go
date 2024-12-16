package pgm

import (
	"context"
	"embed"

	"github.com/ibrt/golang-utils/errorz"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/tern/v2/migrate"
)

// MigrationsConfig describes the configuration for migrations.
type MigrationsConfig struct {
	TableName string
	FS        embed.FS
}

func maybeMustApplyMigrationsInternal(ctx context.Context, pool *pgxpool.Pool, cfg *MigrationsConfig) {
	if cfg == nil {
		return
	}

	conn, err := pool.Acquire(ctx)
	errorz.MaybeMustWrap(err)
	defer conn.Release()

	m, err := migrate.NewMigrator(ctx, conn.Conn(), cfg.TableName)
	errorz.MaybeMustWrap(err)
	errorz.MaybeMustWrap(m.LoadMigrations(cfg.FS))
	errorz.MaybeMustWrap(m.Migrate(ctx))
}
