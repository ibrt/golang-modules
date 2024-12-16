package pgm

import (
	"context"
	"fmt"

	"github.com/ibrt/golang-utils/errorz"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ibrt/golang-modules/logm"
)

var (
	_ RawPG = (*backgroundPGImpl)(nil)
)

type backgroundPGImpl struct {
	pool *pgxpool.Pool
}

// Exec implements the RawPG interface.
func (b *backgroundPGImpl) Exec(ctx context.Context, name, query string, args ...any) (pgconn.CommandTag, error) {
	return logm.Wrap1(
		ctx,
		fmt.Sprintf("pgm.Exec.[%v]", name),
		func(ctx context.Context) (pgconn.CommandTag, error) {
			tag, err := b.pool.Exec(ctx, query, args...)
			return tag, errorz.MaybeWrap(err)
		})
}

// Query implements the RawPG interface.
func (b *backgroundPGImpl) Query(ctx context.Context, name, query string, args ...any) (pgx.Rows, error) {
	return logm.Wrap1(
		ctx,
		fmt.Sprintf("pgm.Query.[%v]", name),
		func(ctx context.Context) (pgx.Rows, error) {
			rows, err := b.pool.Query(ctx, query, args...)
			return rows, errorz.MaybeWrap(err)
		})
}

// QueryRow implements the RawPG interface.
func (b *backgroundPGImpl) QueryRow(ctx context.Context, name, query string, args ...any) pgx.Row {
	row, err := logm.Wrap1(
		ctx,
		fmt.Sprintf("pgm.QueryRow.[%v]", name),
		func(ctx context.Context) (pgx.Row, error) {
			return b.pool.QueryRow(ctx, query, args...), nil
		})
	if err != nil {
		return &errRow{
			err: err,
		}
	}

	return row
}

// Begin implements the RawPG interface.
func (b *backgroundPGImpl) Begin(ctx context.Context, name string, options ...BeginOption) (context.Context, func(), func() error, error) {
	return logm.Wrap3(
		ctx,
		fmt.Sprintf("pgm.Begin.[%v]", name),
		func(ctx context.Context) (context.Context, func(), func() error, error) {
			tx, err := b.pool.BeginTx(ctx, newBeginOptions(options...).ToTxOptions())
			if err != nil {
				return nil, nil, nil, errorz.Wrap(err)
			}

			tP := &transactionPGImpl{
				tx:   tx,
				name: name,
			}

			ctx = NewSingletonInjector(tP)(ctx)
			return ctx, func() { tP.end(ctx) }, func() error { return tP.commit(ctx) }, nil
		})
}
