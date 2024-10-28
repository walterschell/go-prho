package elliptic

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/walterschell/go-prho/field"
)

type CurveParams struct {
	P       string
	N       string
	A       string
	B       string
	Gx      string
	Gy      string
	BitSize int
	Name    string
}

type Curve struct {
	p       *field.Field
	n       *big.Int
	a       *big.Int
	b       *big.Int
	gx      *big.Int
	gy      *big.Int
	bitSize int
	params  *CurveParams
}

func (c *CurveParams) Curve() (*Curve, error) {
	p, ok := new(big.Int).SetString(c.P, 0)
	if !ok {
		return nil, fmt.Errorf("error parsing P")
	}
	n, ok := new(big.Int).SetString(c.N, 0)
	if !ok {
		return nil, fmt.Errorf("error parsing N")
	}
	a, ok := new(big.Int).SetString(c.A, 0)
	if !ok {
		return nil, fmt.Errorf("error parsing A")
	}
	b, ok := new(big.Int).SetString(c.B, 0)
	if !ok {
		return nil, fmt.Errorf("error parsing B")
	}
	gx, ok := new(big.Int).SetString(c.Gx, 0)
	if !ok {
		return nil, fmt.Errorf("error parsing Gx")
	}
	gy, ok := new(big.Int).SetString(c.Gy, 0)
	if !ok {
		return nil, fmt.Errorf("error parsing Gy")
	}
	result := &Curve{
		p:       field.New(p),
		n:       n,
		a:       a,
		b:       b,
		gx:      gx,
		gy:      gy,
		bitSize: c.BitSize,
		params:  c,
	}

	if !result.IsOnCurve(gx, gy) {
		return nil, fmt.Errorf("invalid generator point")
	}
	return result, nil
}

func (c *Curve) N() *big.Int {
	return new(big.Int).Set(c.n)
}

func (c *Curve) IsOnCurve(x, y *big.Int) bool {
	// y^2 mod p = (x^3 + ax + b) mod p
	left := c.p.Exp(y, big.NewInt(2))
	right := c.p.Add(
		c.p.Exp(x, big.NewInt(3)),
		c.p.Mul(c.a, x),
		c.b)
	return left.Cmp(right) == 0
}

type Point struct {
	curve *Curve
	x     *big.Int
	y     *big.Int
}

func (c *Curve) G() *Point {
	return &Point{
		curve: c,
		x:     c.gx,
		y:     c.gy,
	}
}

func (c *Curve) Infinity() *Point {
	return &Point{
		curve: c,
		x:     new(big.Int),
		y:     new(big.Int),
	}
}

func (c *Curve) RandomScalar() *big.Int {
	result, err := rand.Int(rand.Reader, c.n)
	if err != nil {
		panic(err)
	}
	return result
}

func (c *Curve) NewKeypair() (*big.Int, *Point) {
	scalar := c.RandomScalar()
	point := c.G().Mul(scalar)
	return scalar, point
}

func (p *Point) IsOnCurve() bool {
	return p.curve.IsOnCurve(p.x, p.y) || (p.x.Cmp(new(big.Int)) == 0 && p.y.Cmp(new(big.Int)) == 0)
}

func (p *Point) Equal(q *Point) bool {
	return p.curve == q.curve && p.x.Cmp(q.x) == 0 && p.y.Cmp(q.y) == 0
}

func (p *Point) Add(q *Point) *Point {
	if p.curve != q.curve {
		panic("cannot add points on different curves")
	}
	// Assume both points are on the curve
	// Valid because all points are checked to be on the curve at creation time

	infinity := p.curve.Infinity()
	if p.Equal(infinity) {
		return q
	}
	if q.Equal(infinity) {
		return p
	}

	// If y1 != y2, and both points are on curve then q must be the inverse of p
	// and the result is the point at infinity
	if p.x.Cmp(q.x) == 0 && p.y.Cmp(q.y) != 0 {
		return infinity
	}

	var slope *big.Int

	if p.Equal(q) {
		// Slope = (3x^2 + a) / 2y
		numerator := p.curve.p.Add(
			p.curve.p.Mul(big.NewInt(3), p.curve.p.Exp(p.x, big.NewInt(2))),
			p.curve.a)
		denominator := p.curve.p.Mul(big.NewInt(2), p.y)
		slope = p.curve.p.Div(numerator, denominator)
	} else {
		// Slope = (y2 - y1) / (x2 - x1)
		numerator := p.curve.p.Sub(q.y, p.y)
		denominator := p.curve.p.Sub(q.x, p.x)
		slope = p.curve.p.Div(numerator, denominator)
	}

	// x = s^2 - rhs.x - lhs.x
	x := p.curve.p.Sub(
		p.curve.p.Exp(slope, big.NewInt(2)),
		q.x,
		p.x)

	// y = s(lhs.x - x) - lhs.y
	y := p.curve.p.Sub(p.x, x)
	p.curve.p.InplaceMul(y, slope)
	p.curve.p.InplaceSub(y, p.y)

	result := &Point{
		curve: p.curve,
		x:     x,
		y:     y,
	}
	if !result.IsOnCurve() {
		panic("result not on curve")
	}
	return result
}

func (p *Point) Mul(k *big.Int) *Point {
	result := p.curve.Infinity()
	sum := p
	for bitIndex := 0; bitIndex < k.BitLen(); bitIndex++ {
		if k.Bit(bitIndex) == 1 {
			result = result.Add(sum)
		}
		sum = sum.Add(sum)
	}
	return result
}

func (p *Point) Neg() *Point {
	return &Point{
		curve: p.curve,
		x:     p.x,
		y:     p.curve.p.Neg(p.y),
	}
}

func (p *Point) Sub(q *Point) *Point {
	return p.Add(q.Neg())
}

func (p *Point) Curve() *Curve {
	return p.curve
}

func (p *Point) String() string {
	return fmt.Sprintf("(%v, %v)", p.x, p.y)
}

func (p *Point) X() *big.Int {
	return new(big.Int).Set(p.x)
}

func (p *Point) Y() *big.Int {
	return new(big.Int).Set(p.y)
}

var C64Params = &CurveParams{
	A:       "0x197eacf564277a28",
	B:       "0x2869c4a069451233",
	P:       "0xfc477ce0dee80f77",
	BitSize: 64,
	N:       "0xfc477ce09e86adfb",
	Gx:      "0x56853b2bd6052661",
	Gy:      "0xcb116939be0710c9",
}
