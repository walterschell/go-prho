package main
import (
	"testing"
)

func TestCurve(t *testing.T) {
	curve, err := C64Params.Curve()
	if err != nil {
		t.Fatalf("error converting curve: %v", err)
	}
	G := curve.G()

	k1 := curve.RandomScalar()
	k2 := curve.RandomScalar()

	P1 := G.Mul(k1)
	P2 := G.Mul(k2)

	k1P2 := P2.Mul(k1)
	k2P1 := P1.Mul(k2)

	if !k1P2.Equal(k2P1) {
		t.Fatalf("expected k1*P2 == k2*P1")
	}
}