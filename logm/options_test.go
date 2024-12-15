package logm

import (
	"testing"

	"github.com/ibrt/golang-utils/fixturez"
	. "github.com/onsi/gomega"
)

type OptionsSuite struct {
	// intentionally empty
}

func TestOptionsSuite(t *testing.T) {
	fixturez.RunSuite(t, &OptionsSuite{})
}

func (*OptionsSuite) TestEmitOptions(g *WithT) {
	g.Expect(newEmitOptions(
		EmitA(1, 2, 3),
		EmitM("k1", "v1"),
		EmitMetadata{"k2": "v2"})).
		To(Equal(&emitOptions{
			args: []any{1, 2, 3},
			metadata: map[string]any{
				"k1": "v1",
				"k2": "v2",
			},
		}))
}

func (*OptionsSuite) TestBeginOptions(g *WithT) {
	g.Expect(newBeginOptions(
		BeginM("k1", "v1"),
		BeginMetadata{"k2": "v2"},
		BeginErrM("ek1", "ev1"),
		BeginErrMetadata{"ek2": "ev2"})).
		To(Equal(&beginOptions{
			metadata: map[string]any{
				"k1": "v1",
				"k2": "v2",
			},
			errMetadata: map[string]any{
				"ek1": "ev1",
				"ek2": "ev2",
			},
		}))
}
