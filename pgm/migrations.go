package pgm

import (
	"context"
	"embed"

	"github.com/ibrt/golang-utils/errorz"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/tern/v2/migrate"
)

const (
	// MigrationsTableName is the name of the migrations table.
	MigrationsTableName = "migrations"
)

func maybeMustApplyMigrationsInternal(ctx context.Context, pool *pgxpool.Pool, migrationsFS *embed.FS) {
	if migrationsFS == nil {
		return
	}

	conn, err := pool.Acquire(ctx)
	errorz.MaybeMustWrap(err)
	defer conn.Release()

	m, err := migrate.NewMigrator(ctx, conn.Conn(), MigrationsTableName)
	errorz.MaybeMustWrap(err)
	errorz.MaybeMustWrap(m.LoadMigrations(migrationsFS))
	errorz.MaybeMustWrap(m.Migrate(ctx))
}
