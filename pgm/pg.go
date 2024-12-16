// Package pgm implements a Postgres module.
package pgm

import (
	"context"
	"embed"

	"github.com/ibrt/golang-utils/errorz"
	"github.com/ibrt/golang-utils/injectz"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ibrt/golang-modules/cfgm"
	"github.com/ibrt/golang-modules/clkm"
	"github.com/ibrt/golang-modules/logm"
)

type contextKey int

const (
	pgContextKey contextKey = iota
)

// PG describes the module (with cached context).
type PG interface {
	Exec(name, query string, args ...any) (pgconn.CommandTag, error)
	Query(name, query string, args ...any) (pgx.Rows, error)
	QueryRow(name, query string, args ...any) pgx.Row
	Begin(name string, options ...BeginOption) (context.Context, func(), func() error, error)
}

// RawPG describes the module.
type RawPG interface {
	Exec(ctx context.Context, name, query string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, name, query string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, name, query string, args ...any) pgx.Row
	Begin(ctx context.Context, name string, options ...BeginOption) (context.Context, func(), func() error, error)
}

// NewInitializer returns a new [injectz.Initializer] that applies the given migrations (if any).
func NewInitializer(migrationsFS *embed.FS) injectz.Initializer {
	return func(outCtx context.Context) (injectz.Injector, injectz.Releaser) {
		return logm.Wrap2Panic(outCtx, "pgm.Initializer", func(ctx context.Context) (injectz.Injector, injectz.Releaser, error) {
			clkm.MustGet(ctx)
			logm.MustGet(ctx)
			pgCfg := cfgm.MustGet[PGConfigMixin](ctx).GetPGConfig()

			poolCfg, err := pgxpool.ParseConfig(pgCfg.PostgresURL)
			errorz.MaybeMustWrap(err)

			pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
			errorz.MaybeMustWrap(err)

			maybeMustApplyMigrationsInternal(ctx, pool, migrationsFS)
			return NewSingletonInjector(NewPGFromPool(pool)), getReleaser(outCtx, pool), nil
		})
	}
}

func getReleaser(ctx context.Context, pool *pgxpool.Pool) injectz.Releaser {
	return func() {
		_ = logm.Wrap0(ctx, "pgm.Releaser", func(_ context.Context) error {
			pool.Close()
			return nil
		})
	}
}

// NewPGFromPool initializes a new [RawPG] using the given [*pgxpool.Pool].
func NewPGFromPool(pool *pgxpool.Pool) RawPG {
	return &backgroundPGImpl{
		pool: pool,
	}
}

// NewSingletonInjector injects.
func NewSingletonInjector(pg RawPG) injectz.Injector {
	return injectz.NewSingletonInjector(pgContextKey, pg)
}

// MustGet extracts, panics if not found.
func MustGet(ctx context.Context) PG {
	return &adapterPGImpl{
		ctx: ctx,
		pg:  ctx.Value(pgContextKey).(RawPG),
	}
}
