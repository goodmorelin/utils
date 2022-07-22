package pkg

import (
	"github.com/shopspring/decimal"
	"math/big"
	"strings"
)

const MaxPrec = 18

var (
	BigIntEighteen = String2BigInt("1", MaxPrec) // 10^18
	BigIntZero     = String2BigInt("0", MaxPrec) // 0
	maxDecimal     = [18]byte{
		'0', '0', '0', '0', '0', '0',
		'0', '0', '0', '0', '0', '0',
		'0', '0', '0', '0', '0', '0',
	}
)

/*c=a+b*/
func String2BigIntAdd(a, b string, exp int) string {
	if exp < 0 {
		exp = MaxPrec
	}

	aa := String2BigInt(a, exp)
	bb := String2BigInt(b, exp)
	c := new(big.Int).Add(aa, bb)

	return BigInt2String(c, exp)
}

/*c=a-b*/
func String2BigIntSub(a, b string, exp int) string {
	if exp < 0 {
		exp = MaxPrec
	}

	aa := String2BigInt(a, exp)
	bb := String2BigInt(b, exp)
	c := new(big.Int).Sub(aa, bb)

	return BigInt2String(c, exp)
}

/*c=a*b*/
func String2BigIntMul(a, b string, exp int) string {
	if exp < 0 {
		exp = MaxPrec
	}

	aa := String2BigInt(a, exp)
	bb := String2BigInt(b, exp)
	cc := new(big.Int).Div(new(big.Int).Mul(aa, bb), BigIntEighteen)

	return BigInt2String(cc, exp)
}

/*c=a/b*/
func String2BigIntDiv(a, b string, exp int) string {
	if exp < 0 {
		exp = MaxPrec
	}

	bb := String2BigInt(b, exp)
	if bb.Cmp(BigIntZero) == 0 {
		// 分母不能为0
		return BigInt2String(BigIntZero, exp)
	}
	aa := String2BigInt(a, exp)
	c := new(big.Int).Div(new(big.Int).Mul(aa, BigIntEighteen), bb)

	return BigInt2String(c, exp)
}

/*a compare b*/
func String2BigIntCmp(a, b string, exp int) int {
	aa := String2BigInt(a, exp)
	bb := String2BigInt(b, exp)

	//   -1 if a <  b
	//    0 if a == b
	//   +1 if a >  b
	return aa.Cmp(bb)
}

/*abs*/
func String2BigIntAbs(a string, exp int) string {
	if exp < 0 {
		exp = MaxPrec
	}

	aa := String2BigInt(a, exp)
	if aa.Sign() == -1 {
		bb := BigIntSub(BigIntZero, aa)
		return BigInt2String(bb, exp)
	}

	return BigInt2String(aa, exp)
}

// String2BigInt convert string ---> big.Int
func String2BigInt(s string, maxDecimalLength int) *big.Int {
	var index = strings.IndexByte(s, '.')
	// 小数部分所有数字
	var right = maxDecimal
	var rightSlice = right[:maxDecimalLength]

	// 整数
	if index == -1 {
		var i big.Int
		i.SetString(s+string(rightSlice), 10)
		return &i
	}

	// 整数部分
	var left = s[:index]

	// 小数部分
	var d = s[index+1:]
	// 小数部分数字位数
	var l = len(d)

	// 如果小数精度超长，则截取
	if l > maxDecimalLength {
		d = d[:maxDecimalLength]
	}

	// 拷贝小数部分有效数字位数
	copy(rightSlice, d)

	var total = left + string(rightSlice)
	var i big.Int
	i.SetString(total, 10)
	return &i
}

// BigInt2Stringing convert big.Int ---> string
func BigInt2String(h *big.Int, exp int) string {
	if exp < 0 {
		exp = MaxPrec
	}
	return decimal.NewFromBigInt(h, int32(-exp)).String()
}

/*c=a+b*/
func BigIntAdd(a, b *big.Int) *big.Int {
	return new(big.Int).Add(a, b)
}

/*c=a-b*/
func BigIntSub(a, b *big.Int) *big.Int {
	return new(big.Int).Sub(a, b)
}

/*c=a*b*/
func BigIntMul(a, b *big.Int) *big.Int {
	return new(big.Int).Div(new(big.Int).Mul(a, b), BigIntEighteen)
}

/*c=a/b*/
func BigIntDiv(a, b *big.Int) *big.Int {
	// 分母不能为0
	if b.Sign() == 0 {
		return BigIntZero
	}

	return new(big.Int).Div(new(big.Int).Mul(a, BigIntEighteen), b)
}

/*a compare b*/
func BigIntCmp(a, b *big.Int) int {
	//   -1 if a <  b
	//    0 if a == b
	//   +1 if a >  b
	return a.Cmp(b)
}

/*abs*/
func BigIntAbs(a *big.Int) *big.Int {
	if a.Sign() == -1 {
		return BigIntSub(BigIntZero, a)
	}

	return a
}
