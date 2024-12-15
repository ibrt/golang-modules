package logm

import (
	"context"

	"github.com/ibrt/golang-utils/errorz"
)

// Wrap0 traces a function that returns (error).
func Wrap0(
	ctx context.Context,
	name string,
	f func(ctx context.Context) error,
	options ...BeginOption) error {

	ctx, end := MustGet(ctx).Begin(name, options...)
	defer end()

	err := errorz.Catch0Ctx(ctx, func(ctx context.Context) error {
		return errorz.MaybeWrap(f(ctx))
	})
	maybeHandleError(ctx, err)
	return errorz.MaybeWrap(err)
}

// Wrap0Panic traces a function that returns (error). It panics if the returned error is non-nil.
func Wrap0Panic(
	ctx context.Context,
	name string,
	f func(ctx context.Context) error,
	options ...BeginOption) {

	ctx, end := MustGet(ctx).Begin(name, options...)
	defer end()

	err := errorz.Catch0Ctx(ctx, func(ctx context.Context) error {
		return errorz.MaybeWrap(f(ctx))
	})
	maybeHandleError(ctx, err)
	errorz.MaybeMustWrap(err)
}

// Wrap1 traces a function that returns (T, error).
func Wrap1[T any](
	ctx context.Context,
	name string,
	f func(ctx context.Context) (T, error),
	options ...BeginOption) (T, error) {

	ctx, end := MustGet(ctx).Begin(name, options...)
	defer end()

	out, err := errorz.Catch1Ctx(ctx, func(ctx context.Context) (T, error) {
		out, err := f(ctx)
		return out, errorz.MaybeWrap(err)
	})
	maybeHandleError(ctx, err)
	return out, errorz.MaybeWrap(err)
}

// Wrap1Panic traces a function that returns (T, error). It panics if the returned error is non-nil.
func Wrap1Panic[T any](
	ctx context.Context,
	name string,
	f func(ctx context.Context) (T, error),
	options ...BeginOption) T {

	ctx, end := MustGet(ctx).Begin(name, options...)
	defer end()

	out, err := errorz.Catch1Ctx(ctx, func(ctx context.Context) (T, error) {
		out, err := f(ctx)
		return out, errorz.MaybeWrap(err)
	})
	maybeHandleError(ctx, err)
	errorz.MaybeMustWrap(err)
	return out
}

// Wrap2 traces a function that returns (T1, T2, error).
func Wrap2[T1 any, T2 any](
	ctx context.Context,
	name string,
	f func(ctx context.Context) (T1, T2, error),
	options ...BeginOption) (T1, T2, error) {

	ctx, end := MustGet(ctx).Begin(name, options...)
	defer end()

	out1, out2, err := errorz.Catch2Ctx(ctx, func(ctx context.Context) (T1, T2, error) {
		out1, out2, err := f(ctx)
		return out1, out2, errorz.MaybeWrap(err)
	})
	maybeHandleError(ctx, err)
	return out1, out2, errorz.MaybeWrap(err)
}

// Wrap2Panic traces a function that returns (T1, T2, error). It panics if the returned error is non-nil.
func Wrap2Panic[T1 any, T2 any](
	ctx context.Context,
	name string,
	f func(ctx context.Context) (T1, T2, error),
	options ...BeginOption) (T1, T2) {

	ctx, end := MustGet(ctx).Begin(name, options...)
	defer end()

	out1, out2, err := errorz.Catch2Ctx(ctx, func(ctx context.Context) (T1, T2, error) {
		out1, out2, err := f(ctx)
		return out1, out2, errorz.MaybeWrap(err)
	})
	maybeHandleError(ctx, err)
	errorz.MaybeMustWrap(err)
	return out1, out2
}

// Wrap3 traces a function that returns (T1, T2, T2, error).
func Wrap3[T1 any, T2 any, T3 any](
	ctx context.Context,
	name string,
	f func(ctx context.Context) (T1, T2, T3, error),
	options ...BeginOption) (T1, T2, T3, error) {

	ctx, end := MustGet(ctx).Begin(name, options...)
	defer end()

	out1, out2, out3, err := errorz.Catch3Ctx(ctx, func(ctx context.Context) (T1, T2, T3, error) {
		out1, out2, out3, err := f(ctx)
		return out1, out2, out3, errorz.MaybeWrap(err)
	})
	maybeHandleError(ctx, err)
	return out1, out2, out3, errorz.MaybeWrap(err)
}

// Wrap3Panic traces a function that returns (T1, T2, T2, error). It panics if the returned error is non-nil.
func Wrap3Panic[T1 any, T2 any, T3 any](
	ctx context.Context,
	name string,
	f func(ctx context.Context) (T1, T2, T3, error),
	options ...BeginOption) (T1, T2, T3) {

	ctx, end := MustGet(ctx).Begin(name, options...)
	defer end()

	out1, out2, out3, err := errorz.Catch3Ctx(ctx, func(ctx context.Context) (T1, T2, T3, error) {
		out1, out2, out3, err := f(ctx)
		return out1, out2, out3, errorz.MaybeWrap(err)
	})
	maybeHandleError(ctx, err)
	errorz.MaybeMustWrap(err)
	return out1, out2, out3
}

func maybeHandleError(ctx context.Context, err error) {
	if err != nil {
		if getIsEmitted(err) {
			MustGet(ctx).SetErrorFlag()
		} else {
			MustGet(ctx).EmitError(err)
		}
	}
}
