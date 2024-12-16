package pgm

import (
	"context"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/ibrt/golang-utils/errorz"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/ibrt/golang-modules/logm"
)

const (
	txMaxRetries    = 10
	randJitterMinMs = 100
	randJitterMaxMs = 500
)

// Wrap0 wraps a function that returns (error) in a transaction, retrying it if needed.
func Wrap0(
	ctx context.Context,
	name string,
	f func(ctx context.Context) error,
	options ...BeginOption) error {

	return logm.Wrap0(
		ctx,
		fmt.Sprintf("pgm.Wrap.[%v]", name),
		func(ctx context.Context) error {
			var err error

			for i := 0; i < txMaxRetries; i++ {
				if err = wrap0(ctx, name, i, f, options...); err != nil {
					if pgErr := errAsPGError(err); pgErr != nil && pgErr.Code == pgerrcode.SerializationFailure && i < txMaxRetries-1 {
						time.Sleep(randJitter())
						continue
					}
				}

				break
			}

			return errorz.MaybeWrap(err)
		})
}

func wrap0(
	ctx context.Context,
	name string,
	i int,
	f func(ctx context.Context) error,
	options ...BeginOption) error {
	return logm.Wrap0(
		ctx,
		fmt.Sprintf("pgm.Wrap.%v.[%v]", i, name),
		func(ctx context.Context) (err error) {
			ctx, end, commit, err := MustGet(ctx).Begin(name, options...)
			if err != nil {
				return errorz.Wrap(err)
			}
			defer end()

			if err := f(ctx); err != nil {
				return errorz.Wrap(err)
			}

			return errorz.MaybeWrap(commit())
		})
}

// Wrap1 wraps a function that returns (T, error) in a transaction, retrying it if needed.
func Wrap1[T any](
	ctx context.Context,
	name string,
	f func(ctx context.Context) (T, error),
	options ...BeginOption) (T, error) {

	return logm.Wrap1(
		ctx,
		fmt.Sprintf("pgm.Wrap.[%v]", name),
		func(ctx context.Context) (T, error) {
			var t T
			var err error

			for i := 0; i < txMaxRetries; i++ {
				if t, err = wrap1(ctx, name, i, f, options...); err != nil {
					if pgErr := errAsPGError(err); pgErr != nil && pgErr.Code == pgerrcode.SerializationFailure && i < txMaxRetries-1 {
						time.Sleep(randJitter())
						continue
					}
				}

				break
			}

			return t, errorz.MaybeWrap(err)
		})
}

func wrap1[T any](
	ctx context.Context,
	name string,
	i int,
	f func(ctx context.Context) (T, error),
	options ...BeginOption) (T, error) {
	return logm.Wrap1(
		ctx,
		fmt.Sprintf("pgm.Wrap.%v.[%v]", i, name),
		func(ctx context.Context) (T, error) {
			var t T

			ctx, end, commit, err := MustGet(ctx).Begin(name, options...)
			if err != nil {
				return t, errorz.Wrap(err)
			}
			defer end()

			t, err = f(ctx)
			if err != nil {
				return t, errorz.Wrap(err)
			}

			return t, errorz.MaybeWrap(commit())
		})
}

// Wrap2 wraps a function that returns (T1, T2, error) in a transaction, retrying it if needed.
func Wrap2[T1 any, T2 any](
	ctx context.Context,
	name string,
	f func(ctx context.Context) (T1, T2, error),
	options ...BeginOption) (T1, T2, error) {

	return logm.Wrap2(
		ctx,
		fmt.Sprintf("pgm.Wrap.[%v]", name),
		func(ctx context.Context) (T1, T2, error) {
			var t1 T1
			var t2 T2
			var err error

			for i := 0; i < txMaxRetries; i++ {
				if t1, t2, err = wrap2(ctx, name, i, f, options...); err != nil {
					if pgErr := errAsPGError(err); pgErr != nil && pgErr.Code == pgerrcode.SerializationFailure && i < txMaxRetries-1 {
						time.Sleep(randJitter())
						continue
					}
				}

				break
			}

			return t1, t2, errorz.MaybeWrap(err)
		})
}

func wrap2[T1 any, T2 any](
	ctx context.Context,
	name string,
	i int,
	f func(ctx context.Context) (T1, T2, error),
	options ...BeginOption) (T1, T2, error) {
	return logm.Wrap2(
		ctx,
		fmt.Sprintf("pgm.Wrap.%v.[%v]", i, name),
		func(ctx context.Context) (T1, T2, error) {
			var t1 T1
			var t2 T2

			ctx, end, commit, err := MustGet(ctx).Begin(name, options...)
			if err != nil {
				return t1, t2, errorz.Wrap(err)
			}
			defer end()

			t1, t2, err = f(ctx)
			if err != nil {
				return t1, t2, errorz.Wrap(err)
			}

			return t1, t2, errorz.MaybeWrap(commit())
		})
}

// Wrap3 wraps a function that returns (T1, T2, T3, error) in a transaction, retrying it if needed.
func Wrap3[T1 any, T2 any, T3 any](
	ctx context.Context,
	name string,
	f func(ctx context.Context) (T1, T2, T3, error),
	options ...BeginOption) (T1, T2, T3, error) {

	return logm.Wrap3(
		ctx,
		fmt.Sprintf("pgm.Wrap.[%v]", name),
		func(ctx context.Context) (T1, T2, T3, error) {
			var t1 T1
			var t2 T2
			var t3 T3
			var err error

			for i := 0; i < txMaxRetries; i++ {
				if t1, t2, t3, err = wrap3(ctx, name, i, f, options...); err != nil {
					if pgErr := errAsPGError(err); pgErr != nil && pgErr.Code == pgerrcode.SerializationFailure && i < txMaxRetries-1 {
						time.Sleep(randJitter())
						continue
					}
				}

				break
			}

			return t1, t2, t3, errorz.MaybeWrap(err)
		})
}

func wrap3[T1 any, T2 any, T3 any](
	ctx context.Context,
	name string,
	i int,
	f func(ctx context.Context) (T1, T2, T3, error),
	options ...BeginOption) (T1, T2, T3, error) {
	return logm.Wrap3(
		ctx,
		fmt.Sprintf("pgm.Wrap.%v.[%v]", i, name),
		func(ctx context.Context) (T1, T2, T3, error) {
			var t1 T1
			var t2 T2
			var t3 T3

			ctx, end, commit, err := MustGet(ctx).Begin(name, options...)
			if err != nil {
				return t1, t2, t3, errorz.Wrap(err)
			}
			defer end()

			t1, t2, t3, err = f(ctx)
			if err != nil {
				return t1, t2, t3, errorz.Wrap(err)
			}

			return t1, t2, t3, errorz.MaybeWrap(commit())
		})
}

func errAsPGError(err error) *pgconn.PgError {
	if pgErr, ok := errorz.As[*pgconn.PgError](err); ok && pgErr != nil {
		return pgErr
	}
	return nil
}

func randJitter() time.Duration {
	return time.Duration(randJitterMinMs+((randJitterMaxMs-randJitterMinMs)*rand.Float64())) * time.Millisecond
}
