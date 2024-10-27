package main

import (
	"math/big"
)


type Field big.Int

func NewField(p *big.Int) *Field {
	result := new(big.Int)
	result.Set(p)
	return (*Field)(result)
}

func (f *Field) InplaceAdd(a, b *big.Int, extraTerms ...*big.Int) {
	a.Add(a, b)
	for _, term := range extraTerms {
		a.Add(a, term)
	}
	a.Mod(a, (*big.Int)(f))
}

func (f *Field) Add(a, b *big.Int, extraTerms ...*big.Int) *big.Int {
	result := new(big.Int)
	result.Add(a, b)
	for _, term := range extraTerms {
		result.Add(result, term)
	}
	result.Mod(result, (*big.Int)(f))
	return result
}

func (f *Field) InplaceSub(a, b *big.Int, extraTerms ...*big.Int) {
	a.Sub(a, b)
	for _, term := range extraTerms {
		a.Sub(a, term)
	}
	a.Mod(a, (*big.Int)(f))
}

func (f *Field) Sub(a, b *big.Int, extraTerms ...*big.Int) *big.Int {
	result := new(big.Int)
	result.Sub(a, b)
	for _, term := range extraTerms {
		result.Sub(result, term)
	}
	result.Mod(result, (*big.Int)(f))
	return result
}

func (f *Field) InplaceMul(a, b *big.Int) {
	a.Mul(a, b)
	a.Mod(a, (*big.Int)(f))
}

func (f *Field) Mul(a, b *big.Int) *big.Int {
	result := new(big.Int)
	result.Mul(a, b)
	result.Mod(result, (*big.Int)(f))
	return result
}

func (f *Field) ModInvert(a *big.Int) {
	a.ModInverse(a, (*big.Int)(f))
}

func (f *Field) ModInverse(a *big.Int) *big.Int {
	result := new(big.Int)
	result.ModInverse(a, (*big.Int)(f))
	return result
}

func (f *Field) InplaceDiv(a, b *big.Int) {
	a.Mul(a, f.ModInverse(b))
}

func (f *Field) Div(a, b *big.Int) *big.Int {
	result := new(big.Int)
	result.Mul(a, f.ModInverse(b))
	return result
}

func (f *Field) InplaceExp(a, b *big.Int) {
	a.Exp(a, b, (*big.Int)(f))
}

func (f *Field) Exp(a, b *big.Int) *big.Int {
	result := new(big.Int)
	result.Exp(a, b, (*big.Int)(f))
	return result
}

func (f *Field) InplaceSqrt(a *big.Int) {
	check := a.ModSqrt(a, (*big.Int)(f))
	if check == nil {
		panic("no square root")
	}
}

func (f *Field) Sqrt(a *big.Int) *big.Int {
	result := new(big.Int)
	check := result.ModSqrt(a, (*big.Int)(f))
	if check == nil {
		panic("no square root")
	}
	return result
}

func (f *Field) InplaceNeg(a *big.Int) {
	a.Neg(a)
	a.Mod(a, (*big.Int)(f))
}

func (f *Field) Neg(a *big.Int) *big.Int {
	result := new(big.Int)
	result.Neg(a)
	result.Mod(result, (*big.Int)(f))
	return result
}