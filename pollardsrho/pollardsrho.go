package pollardsrho

import (
	"fmt"
	"math/big"

	"github.com/walterschell/go-prho/elliptic"
)

func PollardsRho(p *elliptic.Point) *big.Int {
	// Given P, Q (Points) find x such that xP = Q
	// Find integers a,b,A,B such that aP + bQ = AP + BQ

	i0 := big.NewInt(0)
	i1 := big.NewInt(1)
	i2 := big.NewInt(2)
	i3 := big.NewInt(3)
	upperBound := big.NewInt(50)

	N := p.Curve().N()

	P := p.Curve().G()
	Q := p

	nextPoint := func(pt *elliptic.Point) *elliptic.Point {
		discriminant := pt.X()
		discriminant.Mod(discriminant, i3)
		if discriminant.Cmp(i0) == 0 {
			return pt.Add(P)
		}
		if discriminant.Cmp(i1) == 0 {
			return pt.Add(pt)
		}
		if discriminant.Cmp(i2) == 0 {
			return pt.Add(Q)
		}
		panic("invalid discriminant")
	}
	nextab := func(pt *elliptic.Point, a, b *big.Int) (*big.Int, *big.Int) {
		discriminant := pt.X()
		discriminant.Mod(discriminant, i3)

		A := new(big.Int).Set(a)
		B := new(big.Int).Set(b)

		if discriminant.Cmp(i0) == 0 {
			A.Add(A, i1).Mod(A, N)
			return A, B
		}
		if discriminant.Cmp(i1) == 0 {
			A.Mul(A, i2).Mod(A, N)
			B.Mul(B, i2).Mod(B, N)
			return A, B
		}
		if discriminant.Cmp(i2) == 0 {
			B.Add(B, i1).Mod(B, N)
			return A, B
		}
		panic("invalid discriminant")
	}

	iterationCount := new(big.Int)
	a0 := p.Curve().RandomScalar()
	b0 := p.Curve().RandomScalar()

	R0 := P.Mul(a0).Add(Q.Mul(b0))
	R1 := nextPoint(R0)
	a1, b1 := nextab(R0, a0, b0)

	R2 := nextPoint(R1)
	a2, b2 := nextab(R1, a1, b1)

	for !R1.Equal(R2) {
		// Tortise
		a1, b1 = nextab(R1, a1, b1)
		R1 = nextPoint(R1)

		// Hare
		a2, b2 = nextab(R2, a2, b2)
		R2 = nextPoint(R2)
		a2, b2 = nextab(R2, a2, b2)
		R2 = nextPoint(R2)

		if iterationCount.Cmp(upperBound) > 0 {
			panic("iteration count exceeded")
		}
	}
	if b2.Cmp(b1) == 0 {
		return new(big.Int)
	}
	// (a2 - a1) / (b1 - b2) mod n
	numerator := new(big.Int).Sub(a2, a1)
	numerator.Mod(numerator, N)

	denominator := new(big.Int).Sub(b1, b2)
	denominator.Mod(denominator, N)
	denominator.ModInverse(denominator, N)
	result := numerator.Mul(numerator, denominator)
	result.Mod(result, N)
	return result
}

func main() {
	fmt.Println("Hello, World!")
}
