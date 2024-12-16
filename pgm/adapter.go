package pgm

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	_ PG      = (*adapterPGImpl)(nil)
	_ pgx.Row = (*errRow)(nil)
)

type errRow struct {
	err error
}

// Scan implements the [pgx.Row] interface.
func (r *errRow) Scan(_ ...any) error {
	return r.err
}

type adapterPGImpl struct {
	ctx context.Context
	pg  RawPG
}

// Exec implements the [PG] interface.
func (a *adapterPGImpl) Exec(name, query string, args ...any) (pgconn.CommandTag, error) {
	return a.pg.Exec(a.ctx, name, query, args...)
}

// Query implements the [PG] interface.
func (a *adapterPGImpl) Query(name, query string, args ...any) (pgx.Rows, error) {
	return a.pg.Query(a.ctx, name, query, args...)
}

// QueryRow implements the [PG] interface.
func (a *adapterPGImpl) QueryRow(name, query string, args ...any) pgx.Row {
	return a.pg.QueryRow(a.ctx, name, query, args...)
}

// Begin implements the [PG] interface.
func (a *adapterPGImpl) Begin(name string, options ...BeginOption) (context.Context, func(), func() error, error) {
	return a.pg.Begin(a.ctx, name, options...)
}
