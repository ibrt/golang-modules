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
	g.Expect(newEmitOptions(EmitA(1, 2, 3), EmitM("k", "v"))).
		To(Equal(&emitOptions{
			args:     []any{1, 2, 3},
			metadata: map[string]any{"k": "v"},
		}))
}

func (*OptionsSuite) TestBeginOptions(g *WithT) {
	g.Expect(newBeginOptions(BeginM("k", "v"), BeginErrM("ek", "ev"))).
		To(Equal(&beginOptions{
			metadata:    map[string]any{"k": "v"},
			errMetadata: map[string]any{"ek": "ev"},
		}))
}
