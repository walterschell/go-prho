package main
import (
	"testing"
)
var tinyParameters = CurveParams{
	A:     "63",
	B:     "33",
	P:     "103",
	BitSize: 7,
	N: "97",
	Gx:   "90",
	Gy:   "2",
	Name: "tiny",
}

func TestPollardsRho(t *testing.T) {
	curve, err := tinyParameters.Curve()
	//curve, err := C64Params.Curve()
	if err != nil {
		t.Errorf("Error parsing curve params: %v", err)
	}
	secretKey, publicKey := curve.NewKeypair()

	computedSecretKey := PollardsRho(publicKey)
	if computedSecretKey.Cmp(secretKey) != 0 {
		t.Logf("Computed Secret key * G = %v (expected %v)", curve.G().Mul(computedSecretKey), publicKey)
		t.Errorf("Secret key mismatch. Expected %v, got %v", secretKey, computedSecretKey)
	}
}