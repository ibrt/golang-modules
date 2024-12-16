package pgm

import (
	"context"
	"errors"
	"fmt"

	"github.com/ibrt/golang-utils/errorz"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/ibrt/golang-modules/logm"
)

var (
	_ RawPG = (*transactionPGImpl)(nil)
)

type transactionPGImpl struct {
	tx   pgx.Tx
	name string
}

// Exec implements the [RawPG] interface.
func (t *transactionPGImpl) Exec(ctx context.Context, name, query string, args ...any) (pgconn.CommandTag, error) {
	return logm.Wrap1(
		ctx,
		fmt.Sprintf("pgm.Exec.[%v]", name),
		func(ctx context.Context) (pgconn.CommandTag, error) {
			tag, err := t.tx.Exec(ctx, query, args...)
			return tag, errorz.MaybeWrap(err)
		})
}

// Query implements the [RawPG] interface.
func (t *transactionPGImpl) Query(ctx context.Context, name, query string, args ...any) (pgx.Rows, error) {
	return logm.Wrap1(
		ctx,
		fmt.Sprintf("pgm.Query.[%v]", name),
		func(ctx context.Context) (pgx.Rows, error) {
			rows, err := t.tx.Query(ctx, query, args...)
			return rows, errorz.MaybeWrap(err)
		})
}

// QueryRow implements the [RawPG] interface.
func (t *transactionPGImpl) QueryRow(ctx context.Context, name, query string, args ...any) pgx.Row {
	row, err := logm.Wrap1(
		ctx,
		fmt.Sprintf("pgm.QueryRow.[%v]", name),
		func(ctx context.Context) (pgx.Row, error) {
			return t.tx.QueryRow(ctx, query, args...), nil
		})
	if err != nil {
		return &errRow{
			err: err,
		}
	}

	return row
}

// Begin implements the [RawPG] interface.
func (t *transactionPGImpl) Begin(_ context.Context, _ string, _ ...BeginOption) (context.Context, func(), func() error, error) {
	return nil, nil, nil, errorz.Errorf("nested transaction")
}

func (t *transactionPGImpl) end(ctx context.Context) {
	_ = logm.Wrap0(
		ctx,
		fmt.Sprintf("pgm.End.[%v]", t.name),
		func(ctx context.Context) error {
			if err := t.tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
				return errorz.Wrap(err)
			}
			return nil
		})
}

func (t *transactionPGImpl) commit(ctx context.Context) error {
	return logm.Wrap0(
		ctx,
		fmt.Sprintf("pgm.Commit.[%v]", t.name),
		func(ctx context.Context) error {
			return errorz.MaybeWrap(t.tx.Commit(ctx))
		})
}
