package pgm

import (
	"github.com/jackc/pgx/v5"
)

var (
	_ BeginOption = BeginIsoLevel("")
	_ BeginOption = BeginAccessMode("")
)

// BeginOption describes an option.
type BeginOption interface {
	Apply(o *beginOptions)
}

type beginOptions struct {
	isoLevel   pgx.TxIsoLevel
	accessMode pgx.TxAccessMode
}

func newBeginOptions(options ...BeginOption) *beginOptions {
	o := &beginOptions{
		isoLevel:   pgx.Serializable,
		accessMode: pgx.ReadWrite,
	}

	for _, option := range options {
		option.Apply(o)
	}

	return o
}

// ToTxOptions converts the options to pgx.TxOptions.
func (o *beginOptions) ToTxOptions() pgx.TxOptions {
	return pgx.TxOptions{
		IsoLevel:   o.isoLevel,
		AccessMode: o.accessMode,
	}
}

// BeginIsoLevel describes an isolation level.
type BeginIsoLevel pgx.TxIsoLevel

// Known BeginIsoLevel values.
const (
	BeginIsoLevelSerializable    = BeginIsoLevel(pgx.Serializable)
	BeginIsoLevelRepeatableRead  = BeginIsoLevel(pgx.RepeatableRead)
	BeginIsoLevelReadCommitted   = BeginIsoLevel(pgx.ReadCommitted)
	BeginIsoLevelReadUncommitted = BeginIsoLevel(pgx.ReadUncommitted)
)

// Apply implements the BeginOption interface.
func (i BeginIsoLevel) Apply(o *beginOptions) {
	o.isoLevel = pgx.TxIsoLevel(i)
}

// BeginAccessMode describes an access mode.
type BeginAccessMode pgx.TxAccessMode

// Known BeginAccessMode values.
const (
	BeginAccessModeReadWrite = BeginAccessMode(pgx.ReadWrite)
	BeginAccessModeReadOnly  = BeginAccessMode(pgx.ReadOnly)
)

// Apply implements the BeginOption interface.
func (a BeginAccessMode) Apply(o *beginOptions) {
	o.accessMode = pgx.TxAccessMode(a)
}
