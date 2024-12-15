package logm

var (
	_ EmitOption = (EmitOptionFunc)(nil)
	_ EmitOption = (EmitArgs)(nil)
	_ EmitOption = (EmitMetadata)(nil)
)

// EmitOption describes an option.
type EmitOption interface {
	Apply(o *emitOptions)
}

type emitOptions struct {
	args     []any
	metadata EmitMetadata
}

func newEmitOptions(options ...EmitOption) *emitOptions {
	o := &emitOptions{
		metadata: EmitMetadata{},
	}

	for _, option := range options {
		option.Apply(o)
	}
	return o
}

// EmitOptionFunc is a shorthand for EmitOption.
type EmitOptionFunc func(o *emitOptions)

// Apply implements the EmitOption interface.
func (f EmitOptionFunc) Apply(o *emitOptions) {
	f(o)
}

// EmitArgs describes a list of args.
type EmitArgs []any

// Apply implements the EmitOption interface.
func (a EmitArgs) Apply(o *emitOptions) {
	o.args = append(o.args, a...)
}

// EmitA is a shorthand for EmitArgs.
func EmitA(a ...any) EmitArgs {
	return a
}

// EmitMetadata describes metadata.
type EmitMetadata map[string]any

// Apply implements the EmitOption interface.
func (m EmitMetadata) Apply(o *emitOptions) {
	for k, v := range m {
		o.metadata[k] = v
	}
}

// EmitM is a shorthand for EmitMetadata.
func EmitM(k string, v any) EmitOptionFunc {
	return func(o *emitOptions) {
		o.metadata[k] = v
	}
}

var (
	_ BeginOption = (BeginOptionFunc)(nil)
	_ BeginOption = (BeginMetadata)(nil)
)

// BeginOption describes an option.
type BeginOption interface {
	Apply(o *beginOptions)
}

type beginOptions struct {
	metadata    BeginMetadata
	errMetadata BeginErrMetadata
}

func newBeginOptions(options ...BeginOption) *beginOptions {
	o := &beginOptions{
		metadata:    BeginMetadata{},
		errMetadata: BeginErrMetadata{},
	}

	for _, option := range options {
		option.Apply(o)
	}
	return o
}

// BeginOptionFunc is a shorthand for BeginOption.
type BeginOptionFunc func(o *beginOptions)

// Apply implements the BeginOption interface.
func (f BeginOptionFunc) Apply(o *beginOptions) {
	f(o)
}

// BeginMetadata describes metadata.
type BeginMetadata map[string]any

// Apply implements the BeginOption interface.
func (m BeginMetadata) Apply(o *beginOptions) {
	for k, v := range m {
		o.metadata[k] = v
	}
}

// BeginM is a shorthand for BeginMetadata.
func BeginM(k string, v any) BeginOptionFunc {
	return func(o *beginOptions) {
		o.metadata[k] = v
	}
}

// BeginErrMetadata describes metadata that is applied only in error cases.
type BeginErrMetadata map[string]any

// Apply implements the BeginOption interface.
func (m BeginErrMetadata) Apply(o *beginOptions) {
	for k, v := range m {
		o.errMetadata[k] = v
	}
}

// BeginErrM is a shorthand for BeginErrMetadata.
func BeginErrM(k string, v any) BeginOptionFunc {
	return func(o *beginOptions) {
		o.errMetadata[k] = v
	}
}
