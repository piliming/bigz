package uint512

import (
	"github.com/piliming/bigz/uint128"
	"github.com/piliming/bigz/uint256"

	"errors"
	"math/big"
	"math/bits"
	"slices"
	"unsafe"
)

const bitCount = 512
const byteCount = bitCount / 8
const uint64Count = byteCount / 8

// Note, Zero and Max are functions just to make read-only values.
// We cannot define constants for structures, and global variables
// are unacceptable because it will be possible to change them.

// Zero is the lowest possible Uint512 value.
func Zero() Uint512 {
	return From64(0)
}

// One is the lowest non-zero Uint512 value.
func One() Uint512 {
	return From64(1)
}

// Max is the largest possible Uint512 value.
func Max() Uint512 {
	return Uint512{
		Lo: uint256.Max(),
		Hi: uint256.Max(),
	}
}

type Uint128 = uint128.Uint128
type Uint256 = uint256.Uint256

type Uint512 struct {
	Lo Uint256
	Hi Uint256
}

func From256(v Uint256) Uint512 {
	return Uint512{Lo: v}
}

func From128(v Uint128) Uint512 {
	return Uint512{Lo: uint256.From128(v)}
}

// From64 converts 64-bit value v to a Uint256 value.
// Upper 128-bit half will be zero.
func From64(v uint64) Uint512 {
	return From128(uint128.From64(v))
}

// FromBig converts *big.Int to 256-bit Uint256 value ignoring overflows.
// If input integer is nil or negative then return Zero.
// If input interger overflows 256-bit then return Max.
func FromBig(i *big.Int) Uint512 {
	u, _ := FromBigEx(i)
	return u
}

// FromBigEx converts *big.Int to 256-bit Uint256 value (eXtended version).
// Provides ok successful flag as a second return value.
// If input integer is negative or overflows 512-bit then ok=false.
// If input is nil then zero 512-bit returned.
func FromBigEx(i *big.Int) (Uint512, bool) {
	switch {
	case i == nil:
		return Zero(), true // assuming nil === 0
	case i.Sign() < 0:
		return Zero(), false // value cannot be negative!
	case i.BitLen() > bitCount:
		return Max(), false // value overflows 512-bit!
	}

	buf := i.Bytes()
	slices.Reverse(buf)

	arr := [byteCount]byte{}

	copy(arr[:], buf)

	return LoadLittleEndian(arr), true
}

func (u Uint512) Big() *big.Int {

	buf := [byteCount]byte{}

	copy(buf[:], (*(*[byteCount]byte)(unsafe.Pointer(&u)))[:])

	slices.Reverse(buf[:])

	return new(big.Int).SetBytes(buf[:])
}

func (u Uint512) IsZero() bool {
	return u.Lo.IsZero() && u.Hi.IsZero()
}

func (u Uint512) Equals(v Uint512) bool {
	return u.Lo.Equals(v.Lo) && u.Hi.Equals(v.Hi)
}

func (u Uint512) Equals256(v Uint256) bool {
	return u.Lo.Equals(v) && u.Hi.IsZero()
}

func (u Uint512) Cmp(v Uint512) int {
	if h := u.Hi.Cmp(v.Hi); h != 0 {
		return h
	}
	return u.Lo.Cmp(v.Lo)
}

func (u Uint512) Cmp256(v Uint256) int {
	if !u.Hi.IsZero() {
		return +1 // u > v
	}
	return u.Lo.Cmp(v)
}

///////////////////////////////////////////////////////////////////////////////
/// logical operators /////////////////////////////////////////////////////////

func (u Uint512) Not() Uint512 {
	return Uint512{
		Lo: u.Lo.Not(),
		Hi: u.Hi.Not(),
	}
}

func (u Uint512) AndNot(v Uint512) Uint512 {
	return Uint512{
		Lo: u.Lo.AndNot(v.Lo),
		Hi: u.Hi.AndNot(v.Hi),
	}
}

func (u Uint512) AndNot256(v Uint256) Uint512 {
	return Uint512{
		Lo: u.Lo.AndNot(v),
		Hi: u.Hi, // ^0 == ff..ff
	}
}

func (u Uint512) And(v Uint512) Uint512 {
	return Uint512{
		Lo: u.Lo.And(v.Lo),
		Hi: u.Hi.And(v.Hi),
	}
}

func (u Uint512) And256(v Uint256) Uint512 {
	return Uint512{
		Lo: u.Lo.And(v),
		// Hi: Uint128{0, 0},
	}
}

func (u Uint512) Or(v Uint512) Uint512 {
	return Uint512{
		Lo: u.Lo.Or(v.Lo),
		Hi: u.Hi.Or(v.Hi),
	}
}

func (u Uint512) Or256(v Uint256) Uint512 {
	return Uint512{
		Lo: u.Lo.Or(v),
		Hi: u.Hi,
	}
}

func (u Uint512) Xor(v Uint512) Uint512 {
	return Uint512{
		Lo: u.Lo.Xor(v.Lo),
		Hi: u.Hi.Xor(v.Hi),
	}
}

func (u Uint512) Xor256(v Uint256) Uint512 {
	return Uint512{
		Lo: u.Lo.Xor(v),
		Hi: u.Hi,
	}
}

///////////////////////////////////////////////////////////////////////////////
/// arithmetic operators //////////////////////////////////////////////////////

func Add(x, y Uint512, carry uint64) (sum Uint512, carryOut uint64) {
	sum.Lo, carryOut = uint256.Add(x.Lo, y.Lo, carry)
	sum.Hi, carryOut = uint256.Add(x.Hi, y.Hi, carryOut)
	return
}

func (u Uint512) Add(v Uint512) Uint512 {
	sum, _ := Add(u, v, 0)
	return sum
}

func (u Uint512) Add256(v Uint256) Uint512 {
	lo, c0 := uint256.Add(u.Lo, v, 0)
	return Uint512{Lo: lo, Hi: u.Hi.Add128(Uint128{Lo: c0})}
}

func Sub(x, y Uint512, borrow uint64) (diff Uint512, borrowOut uint64) {
	diff.Lo, borrowOut = uint256.Sub(x.Lo, y.Lo, borrow)
	diff.Hi, borrowOut = uint256.Sub(x.Hi, y.Hi, borrowOut)
	return
}

func (u Uint512) Sub(v Uint512) Uint512 {
	diff, _ := Sub(u, v, 0)
	return diff
}

func (u Uint512) Sub256(v Uint256) Uint512 {
	lo, b0 := uint256.Sub(u.Lo, v, 0)
	return Uint512{Lo: lo, Hi: u.Hi.Sub128(Uint128{Lo: b0})}
}

func Mul(x, y Uint512) (hi, lo Uint512) {
	lo.Hi, lo.Lo = uint256.Mul(x.Lo, y.Lo)
	hi.Hi, hi.Lo = uint256.Mul(x.Hi, y.Hi)
	t0, t1 := uint256.Mul(x.Lo, y.Hi)
	t2, t3 := uint256.Mul(x.Hi, y.Lo)

	var c0, c1 uint64
	lo.Hi, c0 = uint256.Add(lo.Hi, t1, 0)
	lo.Hi, c1 = uint256.Add(lo.Hi, t3, 0)
	hi.Lo, c0 = uint256.Add(hi.Lo, t0, c0)
	hi.Lo, c1 = uint256.Add(hi.Lo, t2, c1)
	hi.Hi = hi.Hi.Add128(uint128.From64(c0 + c1))
	return
}

func (u Uint512) Mul(v Uint512) Uint512 {
	hi, lo := uint256.Mul(u.Lo, v.Lo)
	hi = hi.Add(u.Hi.Mul(v.Lo))
	hi = hi.Add(u.Lo.Mul(v.Hi))
	return Uint512{Lo: lo, Hi: hi}
}

func (u Uint512) Mul256(v Uint256) Uint512 {
	hi, lo := uint256.Mul(u.Lo, v)
	return Uint512{
		Lo: lo,
		Hi: hi.Add(u.Hi.Mul(v)),
	}
}

func (u Uint512) Div(v Uint512) Uint512 {
	q, _ := u.QuoRem(v)
	return q
}

func (u Uint512) Div256(v Uint256) Uint512 {
	q, _ := u.QuoRem256(v)
	return q
}

func (u Uint512) Div128(v Uint128) Uint512 {
	q, _ := u.QuoRem128(v)
	return q
}

func (u Uint512) Div64(v uint64) Uint512 {
	q, _ := u.QuoRem64(v)
	return q
}

func (u Uint512) Mod(v Uint512) Uint512 {
	_, r := u.QuoRem(v)
	return r
}

func (u Uint512) Mod256(v Uint256) Uint256 {
	_, r := u.QuoRem256(v)
	return r
}

func (u Uint512) Mod128(v Uint128) Uint128 {
	_, r := u.QuoRem128(v)
	return r
}

func (u Uint512) Mod64(v uint64) uint64 {
	_, r := u.QuoRem64(v)
	return r
}

func (u Uint512) QuoRem(v Uint512) (Uint512, Uint512) {
	if v.Hi.IsZero() {
		q, r := u.QuoRem256(v.Lo)
		return q, From256(r)
	}

	// generate a "trial quotient," guaranteed to be
	// within 1 of the actual quotient, then adjust.
	n := uint(v.Hi.LeadingZeros())
	u1, v1 := u.Rsh(1), v.Lsh(n)
	tq, _ := uint256.Div(u1.Hi, u1.Lo, v1.Hi)
	tq = tq.Rsh(255 - n)
	if !tq.IsZero() {
		tq = tq.Sub128(uint128.One())
	}

	q, r := From256(tq), u.Sub(v.Mul256(tq))
	if r.Cmp(v) >= 0 {
		q = q.Add256(uint256.One())
		r = r.Sub(v)
	}

	return q, r
}

func (u Uint512) QuoRem256(v Uint256) (Uint512, Uint256) {
	if u.Hi.Cmp(v) < 0 {
		lo, r := uint256.Div(u.Hi, u.Lo, v)
		return Uint512{Lo: lo}, r
	}

	hi, r := uint256.Div(uint256.Zero(), u.Hi, v)
	lo, r := uint256.Div(r, u.Lo, v)
	return Uint512{Lo: lo, Hi: hi}, r
}

func (u Uint512) QuoRem128(v Uint128) (q Uint512, r Uint128) {
	q.Hi, r = u.Hi.QuoRem128(v)

	q.Lo.Hi, r = uint128.Div(r, u.Lo.Hi, v)
	q.Lo.Lo, r = uint128.Div(r, u.Lo.Lo, v)
	return
}

func (u Uint512) QuoRem64(v uint64) (q Uint512, r uint64) {
	q.Hi, r = u.Hi.QuoRem64(v)

	q.Lo.Hi.Hi, r = bits.Div64(r, u.Lo.Hi.Hi, v)
	q.Lo.Hi.Lo, r = bits.Div64(r, u.Lo.Hi.Lo, v)
	q.Lo.Lo.Hi, r = bits.Div64(r, u.Lo.Lo.Hi, v)
	q.Lo.Lo.Lo, r = bits.Div64(r, u.Lo.Lo.Lo, v)

	return
}

func Div(hi, lo, y Uint512) (quo, rem Uint512) {
	if y.IsZero() {
		panic(errors.New("integer divide by zero"))
	}
	if y.Cmp(hi) <= 0 {
		panic(errors.New("integer overflow"))
	}

	s := uint(y.LeadingZeros())
	y = y.Lsh(s)

	un32 := hi.Lsh(s).Or(lo.Rsh(256 - s))
	un10 := lo.Lsh(s)
	q1, rhat := un32.QuoRem256(y.Hi)
	r1 := From256(rhat)

	for !q1.Hi.IsZero() || q1.Mul256(y.Lo).Cmp(Uint512{Hi: r1.Lo, Lo: un10.Hi}) > 0 {
		q1 = q1.Sub256(uint256.One())
		r1 = r1.Add256(y.Hi)
		if !r1.Hi.IsZero() {
			break
		}
	}

	un21 := Uint512{Hi: un32.Lo, Lo: un10.Hi}.Sub(q1.Mul(y))
	q0, rhat := un21.QuoRem256(y.Hi)
	r0 := From256(rhat)

	for !q0.Hi.IsZero() || q0.Mul256(y.Lo).Cmp(Uint512{Hi: r0.Lo, Lo: un10.Lo}) > 0 {
		q0 = q0.Sub256(uint256.One())
		r0 = r0.Add256(y.Hi)
		if !r0.Hi.IsZero() {
			break
		}
	}

	return Uint512{Hi: q1.Lo, Lo: q0.Lo},
		Uint512{Hi: un21.Lo, Lo: un10.Lo}.
			Sub(q0.Mul(y)).Rsh(s)
}

func (u Uint512) Lsh(n uint) Uint512 {
	if n == 0 {
		return u
	}
	if n >= bitCount {
		return Zero()
	}

	if n > 256 {
		return Uint512{
			Lo: Uint256{},
			Hi: u.Lo.Lsh(n - 256),
		}
	}

	if n > 128 {
		return Uint512{
			Lo: Uint256{
				Lo: Uint128{},
				Hi: u.Lo.Lo.Lsh(n - 128),
			},
			Hi: Uint256{
				Lo: u.Lo.Hi.Lsh(n - 128).Or(u.Lo.Lo.Rsh(256 - n)),
				Hi: u.Hi.Lo.Lsh(n - 128).Or(u.Lo.Hi.Rsh(256 - n)),
			},
		}
	}

	if n > 64 {
		return Uint512{
			Lo: Uint256{
				Lo: Uint128{
					Lo: 0,
					Hi: u.Lo.Lo.Lo << (n - 64),
				},
				Hi: Uint128{
					Lo: u.Lo.Lo.Hi<<(n-64) | u.Lo.Lo.Lo>>(128-n),
					Hi: u.Lo.Hi.Lo<<(n-64) | u.Lo.Lo.Hi>>(128-n),
				},
			},
			Hi: Uint256{
				Lo: Uint128{
					Lo: u.Lo.Hi.Hi<<(n-64) | u.Lo.Hi.Lo>>(128-n),
					Hi: u.Hi.Lo.Lo<<(n-64) | u.Lo.Hi.Hi>>(128-n),
				},
				Hi: Uint128{
					Lo: u.Hi.Lo.Hi<<(n-64) | u.Hi.Lo.Lo>>(128-n),
					Hi: u.Hi.Hi.Lo<<(n-64) | u.Hi.Lo.Hi>>(128-n),
				},
			},
		}
	}

	return Uint512{
		Lo: Uint256{
			Lo: Uint128{
				Lo: u.Lo.Lo.Lo << n,
				Hi: u.Lo.Lo.Hi<<n | u.Lo.Lo.Lo>>(64-n),
			},
			Hi: Uint128{
				Lo: u.Lo.Hi.Lo<<n | u.Lo.Lo.Hi>>(64-n),
				Hi: u.Lo.Hi.Hi<<n | u.Lo.Hi.Lo>>(64-n),
			},
		},
		Hi: Uint256{
			Lo: Uint128{
				Lo: u.Hi.Lo.Lo<<n | u.Lo.Hi.Hi>>(64-n),
				Hi: u.Hi.Lo.Hi<<n | u.Hi.Lo.Lo>>(64-n),
			},
			Hi: Uint128{
				Lo: u.Hi.Hi.Lo<<n | u.Hi.Lo.Hi>>(64-n),
				Hi: u.Hi.Hi.Hi<<n | u.Hi.Hi.Lo>>(64-n),
			},
		},
	}
}

func (u Uint512) Rsh(n uint) Uint512 {
	if n == 0 {
		return u
	}
	if n >= bitCount {
		return Zero()
	}

	if n > 256 {
		return Uint512{
			Lo: u.Hi.Rsh(n - 256),
			Hi: Uint256{},
		}
	}

	if n > 128 {
		return Uint512{
			Lo: Uint256{
				Lo: u.Lo.Hi.Rsh(n - 128).Or(u.Hi.Lo.Lsh(256 - n)),
				Hi: u.Hi.Lo.Rsh(n - 128).Or(u.Hi.Hi.Lsh(256 - n)),
			},
			Hi: Uint256{
				Lo: u.Hi.Hi.Rsh(n - 128),
				Hi: Uint128{},
			},
		}
	}

	if n > 64 {
		return Uint512{
			Lo: Uint256{
				Lo: Uint128{
					Lo: u.Lo.Lo.Hi>>(n-64) | u.Lo.Hi.Lo<<(128-n),
					Hi: u.Lo.Hi.Lo>>(n-64) | u.Lo.Hi.Hi<<(128-n),
				},
				Hi: Uint128{
					Lo: u.Lo.Hi.Hi>>(n-64) | u.Hi.Lo.Lo<<(128-n),
					Hi: u.Hi.Lo.Lo>>(n-64) | u.Hi.Lo.Hi<<(128-n),
				},
			},
			Hi: Uint256{
				Lo: Uint128{
					Lo: u.Hi.Lo.Hi>>(n-64) | u.Hi.Hi.Lo<<(128-n),
					Hi: u.Hi.Hi.Lo>>(n-64) | u.Hi.Hi.Hi<<(128-n),
				},
				Hi: Uint128{
					Lo: u.Hi.Hi.Hi >> (n - 64),
					Hi: 0,
				},
			},
		}
	}

	return Uint512{
		Lo: Uint256{
			Lo: Uint128{
				Lo: u.Lo.Lo.Lo>>n | u.Lo.Lo.Hi<<(64-n),
				Hi: u.Lo.Lo.Hi>>n | u.Lo.Hi.Lo<<(64-n),
			},
			Hi: Uint128{
				Lo: u.Lo.Hi.Lo>>n | u.Lo.Hi.Hi<<(64-n),
				Hi: u.Lo.Hi.Hi>>n | u.Hi.Lo.Lo<<(64-n),
			},
		},
		Hi: Uint256{
			Lo: Uint128{
				Lo: u.Hi.Lo.Lo>>n | u.Hi.Lo.Hi<<(64-n),
				Hi: u.Hi.Lo.Hi>>n | u.Hi.Hi.Lo<<(64-n),
			},
			Hi: Uint128{
				Lo: u.Hi.Hi.Lo>>n | u.Hi.Hi.Hi<<(64-n),
				Hi: u.Hi.Hi.Hi >> n,
			},
		},
	}
}

func (u Uint512) BitLen() int {
	if !u.Hi.IsZero() {
		return 256 + u.Hi.BitLen()
	}
	return u.Lo.BitLen()
}

func (u Uint512) LeadingZeros() int {
	if !u.Hi.IsZero() {
		return u.Hi.LeadingZeros()
	}
	return 256 + u.Lo.LeadingZeros()
}

func (u Uint512) TrailingZeros() int {
	if !u.Lo.IsZero() {
		return u.Lo.TrailingZeros()
	}
	return 256 + u.Hi.TrailingZeros()
}

func (u Uint512) OnesCount() int {
	return u.Lo.OnesCount() +
		u.Hi.OnesCount()
}

func (u Uint512) Reverse() Uint512 {
	return Uint512{
		Lo: u.Hi.Reverse(),
		Hi: u.Lo.Reverse(),
	}
}

func (u Uint512) ReverseBytes() Uint512 {
	return Uint512{
		Lo: u.Hi.ReverseBytes(),
		Hi: u.Lo.ReverseBytes(),
	}
}

func (u Uint512) Bit(n int) bool {
	if n < 0 || n >= bitCount {
		return false
	}
	word := n / 64
	bit := uint(n % 64)
	return ((*(*[uint64Count]uint64)(unsafe.Pointer(&u)))[word]>>bit)&1 == 1
}
