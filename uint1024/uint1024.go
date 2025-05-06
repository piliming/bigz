package uint1024

import (
	"errors"
	"github.com/piliming/bigz/uint128"
	"github.com/piliming/bigz/uint256"
	"github.com/piliming/bigz/uint512"
	"math/big"
	"math/bits"
	"slices"
	"unsafe"
)

const bitCount = 1024
const byteCount = bitCount / 8
const uint64Count = byteCount / 8

func Zero() Uint1024 {
	return Uint1024{}
}

var one = From64(1)

func One() Uint1024 {
	return one
}

func Max() Uint1024 {
	return Uint1024{
		Lo: uint512.Max(),
		Hi: uint512.Max(),
	}
}

type Uint128 = uint128.Uint128
type Uint256 = uint256.Uint256
type Uint512 = uint512.Uint512

type Uint1024 struct {
	Lo Uint512
	Hi Uint512
}

func From512(v Uint512) Uint1024 {
	return Uint1024{
		Lo: v,
	}
}

func From256(v Uint256) Uint1024 {
	return Uint1024{
		Lo: uint512.From256(v),
	}
}

func From128(v Uint128) Uint1024 {
	return Uint1024{
		Lo: uint512.From128(v),
	}
}

func From64(v uint64) Uint1024 {
	return Uint1024{
		Lo: uint512.From64(v),
	}
}

func FromBig(i *big.Int) Uint1024 {
	u, _ := FromBigEx(i)
	return u
}

// FromBigEx converts *big.Int to 256-bit Uint256 value (eXtended version).
// Provides ok successful flag as a second return value.
// If input integer is negative or overflows 1024-bit then ok=false.
// If input is nil then zero 1024-bit returned.
func FromBigEx(i *big.Int) (Uint1024, bool) {
	switch {
	case i == nil:
		return Zero(), true // assuming nil === 0
	case i.Sign() < 0:
		return Zero(), false // value cannot be negative!
	case i.BitLen() > bitCount:
		return Max(), false // value overflows 1025-bit!
	}

	buf := i.Bytes()
	slices.Reverse(buf)

	arr := [byteCount]byte{}

	copy(arr[:], buf)

	return LoadLittleEndian(arr), true
}

func (u Uint1024) Big() *big.Int {
	buf := [byteCount]byte{}

	copy(buf[:], (*(*[byteCount]byte)(unsafe.Pointer(&u)))[:])

	slices.Reverse(buf[:])

	return new(big.Int).SetBytes(buf[:])
}

func (u Uint1024) IsZero() bool {
	return u.Lo.IsZero() && u.Hi.IsZero()
}

func (u Uint1024) Equals(v Uint1024) bool {
	return u.Lo.Equals(v.Lo) && u.Hi.Equals(v.Hi)
}

func (u Uint1024) Equals512(v Uint512) bool {
	return u.Lo.Equals(v) && u.Hi.IsZero()
}

func (u Uint1024) Cmp(v Uint1024) int {
	if h := u.Hi.Cmp(v.Hi); h != 0 {
		return h
	}
	return u.Lo.Cmp(v.Lo)
}

func (u Uint1024) Cmp512(v Uint512) int {
	if !u.Hi.IsZero() {
		return +1 // u > v
	}
	return u.Lo.Cmp(v)
}

///////////////////////////////////////////////////////////////////////////////
/// logical operators /////////////////////////////////////////////////////////

func (u Uint1024) Not() Uint1024 {
	return Uint1024{
		Lo: u.Lo.Not(),
		Hi: u.Hi.Not(),
	}
}

func (u Uint1024) AndNot(v Uint1024) Uint1024 {
	return Uint1024{
		Lo: u.Lo.AndNot(v.Lo),
		Hi: u.Hi.AndNot(v.Hi),
	}
}

func (u Uint1024) AndNot512(v Uint512) Uint1024 {
	return Uint1024{
		Lo: u.Lo.AndNot(v),
		Hi: u.Hi, // ^0 == ff..ff
	}
}

func (u Uint1024) And(v Uint1024) Uint1024 {
	return Uint1024{
		Lo: u.Lo.And(v.Lo),
		Hi: u.Hi.And(v.Hi),
	}
}

func (u Uint1024) And512(v Uint512) Uint1024 {
	return Uint1024{
		Lo: u.Lo.And(v),
		// Hi: Uint128{0, 0},
	}
}

func (u Uint1024) Or(v Uint1024) Uint1024 {
	return Uint1024{
		Lo: u.Lo.Or(v.Lo),
		Hi: u.Hi.Or(v.Hi),
	}
}

func (u Uint1024) Or512(v Uint512) Uint1024 {
	return Uint1024{
		Lo: u.Lo.Or(v),
		Hi: u.Hi,
	}
}

func (u Uint1024) Xor(v Uint1024) Uint1024 {
	return Uint1024{
		Lo: u.Lo.Xor(v.Lo),
		Hi: u.Hi.Xor(v.Hi),
	}
}

func (u Uint1024) Xor512(v Uint512) Uint1024 {
	return Uint1024{
		Lo: u.Lo.Xor(v),
		Hi: u.Hi,
	}
}

///////////////////////////////////////////////////////////////////////////////
/// arithmetic operators //////////////////////////////////////////////////////

func Add(x, y Uint1024, carry uint64) (sum Uint1024, carryOut uint64) {
	sum.Lo, carryOut = uint512.Add(x.Lo, y.Lo, carry)
	sum.Hi, carryOut = uint512.Add(x.Hi, y.Hi, carryOut)
	return
}

func (u Uint1024) Add(v Uint1024) Uint1024 {
	sum, _ := Add(u, v, 0)
	return sum
}

func (u Uint1024) Add512(v Uint512) Uint1024 {
	lo, c0 := uint512.Add(u.Lo, v, 0)
	return Uint1024{Lo: lo, Hi: u.Hi.Add(uint512.From64(c0))}
}

func Sub(x, y Uint1024, borrow uint64) (diff Uint1024, borrowOut uint64) {
	diff.Lo, borrowOut = uint512.Sub(x.Lo, y.Lo, borrow)
	diff.Hi, borrowOut = uint512.Sub(x.Hi, y.Hi, borrowOut)
	return
}

func (u Uint1024) Sub(v Uint1024) Uint1024 {
	diff, _ := Sub(u, v, 0)
	return diff
}

func (u Uint1024) Sub512(v Uint512) Uint1024 {
	lo, b0 := uint512.Sub(u.Lo, v, 0)
	return Uint1024{Lo: lo, Hi: u.Hi.Sub(uint512.From64(b0))}
}

func Mul(x, y Uint1024) (hi, lo Uint1024) {
	lo.Hi, lo.Lo = uint512.Mul(x.Lo, y.Lo)
	hi.Hi, hi.Lo = uint512.Mul(x.Hi, y.Hi)
	t0, t1 := uint512.Mul(x.Lo, y.Hi)
	t2, t3 := uint512.Mul(x.Hi, y.Lo)

	var c0, c1 uint64
	lo.Hi, c0 = uint512.Add(lo.Hi, t1, 0)
	lo.Hi, c1 = uint512.Add(lo.Hi, t3, 0)
	hi.Lo, c0 = uint512.Add(hi.Lo, t0, c0)
	hi.Lo, c1 = uint512.Add(hi.Lo, t2, c1)
	hi.Hi = hi.Hi.Add256(uint256.From64(c0 + c1))
	return
}

func (u Uint1024) Mul(v Uint1024) Uint1024 {
	hi, lo := uint512.Mul(u.Lo, v.Lo)
	hi = hi.Add(u.Hi.Mul(v.Lo))
	hi = hi.Add(u.Lo.Mul(v.Hi))
	return Uint1024{Lo: lo, Hi: hi}
}

func (u Uint1024) Mul512(v Uint512) Uint1024 {
	hi, lo := uint512.Mul(u.Lo, v)
	return Uint1024{
		Lo: lo,
		Hi: hi.Add(u.Hi.Mul(v)),
	}
}

func (u Uint1024) Div(v Uint1024) Uint1024 {
	q, _ := u.QuoRem(v)
	return q
}

func (u Uint1024) Div512(v Uint512) Uint1024 {
	q, _ := u.QuoRem512(v)
	return q
}

func (u Uint1024) Div256(v Uint256) Uint1024 {
	q, _ := u.QuoRem256(v)
	return q
}

func (u Uint1024) Div128(v Uint128) Uint1024 {
	q, _ := u.QuoRem128(v)
	return q
}

func (u Uint1024) Div64(v uint64) Uint1024 {
	q, _ := u.QuoRem64(v)
	return q
}

func (u Uint1024) Mod(v Uint1024) Uint1024 {
	_, r := u.QuoRem(v)
	return r
}

func (u Uint1024) Mod512(v Uint512) Uint512 {
	_, r := u.QuoRem512(v)
	return r
}

func (u Uint1024) Mod256(v Uint256) Uint256 {
	_, r := u.QuoRem256(v)
	return r
}

func (u Uint1024) Mod128(v Uint128) Uint128 {
	_, r := u.QuoRem128(v)
	return r
}

func (u Uint1024) Mod64(v uint64) uint64 {
	_, r := u.QuoRem64(v)
	return r
}

func (u Uint1024) QuoRem(v Uint1024) (Uint1024, Uint1024) {
	if v.Hi.IsZero() {
		q, r := u.QuoRem512(v.Lo)
		return q, From512(r)
	}

	// generate a "trial quotient," guaranteed to be
	// within 1 of the actual quotient, then adjust.
	n := uint(v.Hi.LeadingZeros())
	u1, v1 := u.Rsh(1), v.Lsh(n)
	tq, _ := uint512.Div(u1.Hi, u1.Lo, v1.Hi)
	tq = tq.Rsh(bitCount/2 - 1 - n)
	if !tq.IsZero() {
		tq = tq.Sub256(uint256.One())
	}

	q, r := From512(tq), u.Sub(v.Mul512(tq))
	if r.Cmp(v) >= 0 {
		q = q.Add512(uint512.One())
		r = r.Sub(v)
	}

	return q, r
}

func (u Uint1024) QuoRem512(v Uint512) (Uint1024, Uint512) {
	if u.Hi.Cmp(v) < 0 {
		lo, r := uint512.Div(u.Hi, u.Lo, v)
		return Uint1024{Lo: lo}, r
	}

	hi, r := uint512.Div(uint512.Zero(), u.Hi, v)
	lo, r := uint512.Div(r, u.Lo, v)

	return Uint1024{Lo: lo, Hi: hi}, r
}

func (u Uint1024) QuoRem256(v Uint256) (q Uint1024, r Uint256) {
	q.Hi, r = u.Hi.QuoRem256(v)

	q.Lo.Hi, r = uint256.Div(r, u.Lo.Hi, v)
	q.Lo.Lo, r = uint256.Div(r, u.Lo.Lo, v)

	return
}

func (u Uint1024) QuoRem128(v Uint128) (q Uint1024, r Uint128) {
	q.Hi, r = u.Hi.QuoRem128(v)

	q.Lo.Hi.Hi, r = uint128.Div(r, u.Lo.Hi.Hi, v)
	q.Lo.Hi.Lo, r = uint128.Div(r, u.Lo.Hi.Lo, v)
	q.Lo.Lo.Hi, r = uint128.Div(r, u.Lo.Lo.Hi, v)
	q.Lo.Lo.Lo, r = uint128.Div(r, u.Lo.Lo.Lo, v)

	return
}

func (u Uint1024) QuoRem64(v uint64) (q Uint1024, r uint64) {
	q.Hi, r = u.Hi.QuoRem64(v)

	q.Lo.Hi.Hi.Hi, r = bits.Div64(r, u.Lo.Hi.Hi.Hi, v)
	q.Lo.Hi.Hi.Lo, r = bits.Div64(r, u.Lo.Hi.Hi.Lo, v)

	q.Lo.Hi.Lo.Hi, r = bits.Div64(r, u.Lo.Hi.Lo.Hi, v)
	q.Lo.Hi.Lo.Lo, r = bits.Div64(r, u.Lo.Hi.Lo.Lo, v)

	q.Lo.Lo.Hi.Hi, r = bits.Div64(r, u.Lo.Lo.Hi.Hi, v)
	q.Lo.Lo.Hi.Lo, r = bits.Div64(r, u.Lo.Lo.Hi.Lo, v)

	q.Lo.Lo.Lo.Hi, r = bits.Div64(r, u.Lo.Lo.Lo.Hi, v)
	q.Lo.Lo.Lo.Lo, r = bits.Div64(r, u.Lo.Lo.Lo.Lo, v)

	return
}

func Div(hi, lo, y Uint1024) (quo, rem Uint1024) {
	if y.IsZero() {
		panic(errors.New("integer divide by zero"))
	}
	if y.Cmp(hi) <= 0 {
		panic(errors.New("integer overflow"))
	}

	s := uint(y.LeadingZeros())
	y = y.Lsh(s)

	un32 := hi.Lsh(s).Or(lo.Rsh(512 - s))
	un10 := lo.Lsh(s)
	q1, rhat := un32.QuoRem512(y.Hi)
	r1 := From512(rhat)

	for !q1.Hi.IsZero() || q1.Mul512(y.Lo).Cmp(Uint1024{Hi: r1.Lo, Lo: un10.Hi}) > 0 {
		q1 = q1.Sub(One())
		r1 = r1.Add512(y.Hi)
		if !r1.Hi.IsZero() {
			break
		}
	}

	un21 := Uint1024{Hi: un32.Lo, Lo: un10.Hi}.Sub(q1.Mul(y))
	q0, rhat := un21.QuoRem512(y.Hi)
	r0 := From512(rhat)

	for !q0.Hi.IsZero() || q0.Mul512(y.Lo).Cmp(Uint1024{Hi: r0.Lo, Lo: un10.Lo}) > 0 {
		q0 = q0.Sub(One())
		r0 = r0.Add512(y.Hi)
		if !r0.Hi.IsZero() {
			break
		}
	}

	return Uint1024{Hi: q1.Lo, Lo: q0.Lo},
		Uint1024{Hi: un21.Lo, Lo: un10.Lo}.
			Sub(q0.Mul(y)).Rsh(s)
}

func (u Uint1024) Lsh(n uint) Uint1024 {
	if n == 0 {
		return u
	}
	if n >= bitCount {
		return Zero()
	}

	if n > 512 {
		return Uint1024{
			Lo: Uint512{},
			Hi: u.Lo.Lsh(n - 512),
		}
	}

	if n > 256 {
		return Uint1024{
			Lo: Uint512{
				Lo: Uint256{},
				Hi: u.Lo.Lo.Lsh(n - 256),
			},
			Hi: Uint512{
				Lo: u.Lo.Hi.Lsh(n - 256).Or(u.Lo.Lo.Rsh(512 - n)),
				Hi: u.Hi.Lo.Lsh(n - 256).Or(u.Lo.Hi.Rsh(512 - n)),
			},
		}
	}

	if n > 128 {
		return Uint1024{
			Lo: Uint512{
				Lo: Uint256{
					Lo: Uint128{},
					Hi: u.Lo.Lo.Lo.Lsh(n - 128),
				},
				Hi: Uint256{
					Lo: u.Lo.Lo.Hi.Lsh(n - 128).Or(u.Lo.Lo.Lo.Rsh(256 - n)),
					Hi: u.Lo.Hi.Lo.Lsh(n - 128).Or(u.Lo.Lo.Hi.Rsh(256 - n)),
				},
			},
			Hi: Uint512{
				Lo: Uint256{
					Lo: u.Lo.Hi.Hi.Lsh(n - 128).Or(u.Lo.Hi.Lo.Rsh(256 - n)),
					Hi: u.Hi.Lo.Lo.Lsh(n - 128).Or(u.Lo.Hi.Hi.Rsh(256 - n)),
				},
				Hi: Uint256{
					Lo: u.Hi.Lo.Hi.Lsh(n - 128).Or(u.Hi.Lo.Lo.Rsh(256 - n)),
					Hi: u.Hi.Hi.Lo.Lsh(n - 128).Or(u.Hi.Lo.Hi.Rsh(256 - n)),
				},
			},
		}
	}

	if n > 64 {
		return Uint1024{
			Lo: Uint512{
				Lo: Uint256{
					Lo: Uint128{
						Lo: 0,
						Hi: u.Lo.Lo.Lo.Lo << (n - 64),
					},
					Hi: Uint128{
						Lo: u.Lo.Lo.Lo.Hi<<(n-64) | u.Lo.Lo.Lo.Lo>>(128-n),
						Hi: u.Lo.Lo.Hi.Lo<<(n-64) | u.Lo.Lo.Lo.Hi>>(128-n),
					},
				},
				Hi: Uint256{
					Lo: Uint128{
						Lo: u.Lo.Lo.Hi.Hi<<(n-64) | u.Lo.Lo.Hi.Lo>>(128-n),
						Hi: u.Lo.Hi.Lo.Lo<<(n-64) | u.Lo.Lo.Hi.Hi>>(128-n),
					},
					Hi: Uint128{
						Lo: u.Lo.Hi.Lo.Hi<<(n-64) | u.Lo.Hi.Lo.Lo>>(128-n),
						Hi: u.Lo.Hi.Hi.Lo<<(n-64) | u.Lo.Hi.Lo.Hi>>(128-n),
					},
				},
			},
			Hi: Uint512{
				Lo: Uint256{
					Lo: Uint128{
						Lo: u.Lo.Hi.Hi.Hi<<(n-64) | u.Lo.Hi.Hi.Lo>>(128-n),
						Hi: u.Hi.Lo.Lo.Lo<<(n-64) | u.Lo.Hi.Hi.Hi>>(128-n),
					},
					Hi: Uint128{
						Lo: u.Hi.Lo.Lo.Hi<<(n-64) | u.Hi.Lo.Lo.Lo>>(128-n),
						Hi: u.Hi.Lo.Hi.Lo<<(n-64) | u.Hi.Lo.Lo.Hi>>(128-n),
					},
				},
				Hi: Uint256{
					Lo: Uint128{
						Lo: u.Hi.Lo.Hi.Hi<<(n-64) | u.Hi.Lo.Hi.Lo>>(128-n),
						Hi: u.Hi.Hi.Lo.Lo<<(n-64) | u.Hi.Lo.Hi.Hi>>(128-n),
					},
					Hi: Uint128{
						Lo: u.Hi.Hi.Lo.Hi<<(n-64) | u.Hi.Hi.Lo.Lo>>(128-n),
						Hi: u.Hi.Hi.Hi.Lo<<(n-64) | u.Hi.Hi.Lo.Hi>>(128-n),
					},
				},
			},
		}
	}

	return Uint1024{
		Lo: Uint512{
			Lo: Uint256{
				Lo: Uint128{
					Lo: u.Lo.Lo.Lo.Lo << n,
					Hi: u.Lo.Lo.Lo.Hi<<n | u.Lo.Lo.Lo.Lo>>(64-n),
				},
				Hi: Uint128{
					Lo: u.Lo.Lo.Hi.Lo<<n | u.Lo.Lo.Lo.Hi>>(64-n),
					Hi: u.Lo.Lo.Hi.Hi<<n | u.Lo.Lo.Hi.Lo>>(64-n),
				},
			},
			Hi: Uint256{
				Lo: Uint128{
					Lo: u.Lo.Hi.Lo.Lo<<n | u.Lo.Lo.Hi.Hi>>(64-n),
					Hi: u.Lo.Hi.Lo.Hi<<n | u.Lo.Hi.Lo.Lo>>(64-n),
				},
				Hi: Uint128{
					Lo: u.Lo.Hi.Hi.Lo<<n | u.Lo.Hi.Lo.Hi>>(64-n),
					Hi: u.Lo.Hi.Hi.Hi<<n | u.Lo.Hi.Hi.Lo>>(64-n),
				},
			},
		},
		Hi: Uint512{
			Lo: Uint256{
				Lo: Uint128{
					Lo: u.Hi.Lo.Lo.Lo<<n | u.Lo.Hi.Hi.Hi>>(64-n),
					Hi: u.Hi.Lo.Lo.Hi<<n | u.Hi.Lo.Lo.Lo>>(64-n),
				},
				Hi: Uint128{
					Lo: u.Hi.Lo.Hi.Lo<<n | u.Hi.Lo.Lo.Hi>>(64-n),
					Hi: u.Hi.Lo.Hi.Hi<<n | u.Hi.Lo.Hi.Lo>>(64-n),
				},
			},
			Hi: Uint256{
				Lo: Uint128{
					Lo: u.Hi.Hi.Lo.Lo<<n | u.Hi.Lo.Hi.Hi>>(64-n),
					Hi: u.Hi.Hi.Lo.Hi<<n | u.Hi.Hi.Lo.Lo>>(64-n),
				},
				Hi: Uint128{
					Lo: u.Hi.Hi.Hi.Lo<<n | u.Hi.Hi.Lo.Hi>>(64-n),
					Hi: u.Hi.Hi.Hi.Hi<<n | u.Hi.Hi.Hi.Lo>>(64-n),
				},
			},
		},
	}
}

func (u Uint1024) Rsh(n uint) Uint1024 {
	if n == 0 {
		return u
	}
	if n >= bitCount {
		return Zero()
	}

	if n > 512 {
		return Uint1024{
			Lo: u.Hi.Rsh(n - 512),
			Hi: Uint512{},
		}
	}

	if n > 256 {
		return Uint1024{
			Lo: Uint512{
				Lo: u.Lo.Hi.Rsh(n - 256).Or(u.Hi.Lo.Lsh(512 - n)),
				Hi: u.Hi.Lo.Rsh(n - 256).Or(u.Hi.Hi.Lsh(512 - n)),
			},
			Hi: Uint512{
				Lo: u.Hi.Hi.Rsh(n - 256),
				Hi: Uint256{},
			},
		}
	}

	if n > 128 {
		return Uint1024{
			Lo: Uint512{
				Lo: Uint256{
					Lo: u.Lo.Lo.Hi.Rsh(n - 128).Or(u.Lo.Hi.Lo.Lsh(256 - n)),
					Hi: u.Lo.Hi.Lo.Rsh(n - 128).Or(u.Lo.Hi.Hi.Lsh(256 - n)),
				},
				Hi: Uint256{
					Lo: u.Lo.Hi.Hi.Rsh(n - 128).Or(u.Hi.Lo.Lo.Lsh(256 - n)),
					Hi: u.Hi.Lo.Lo.Rsh(n - 128).Or(u.Hi.Lo.Hi.Lsh(256 - n)),
				},
			},
			Hi: Uint512{
				Lo: Uint256{
					Lo: u.Hi.Lo.Hi.Rsh(n - 128).Or(u.Hi.Hi.Lo.Lsh(256 - n)),
					Hi: u.Hi.Hi.Lo.Rsh(n - 128).Or(u.Hi.Hi.Hi.Lsh(256 - n)),
				},
				Hi: Uint256{
					Lo: u.Hi.Hi.Hi.Rsh(n - 128),
					Hi: Uint128{},
				},
			},
		}
	}

	if n > 64 {
		return Uint1024{
			Lo: Uint512{
				Lo: Uint256{
					Lo: Uint128{
						Lo: u.Lo.Lo.Lo.Hi>>(n-64) | u.Lo.Lo.Hi.Lo<<(128-n),
						Hi: u.Lo.Lo.Hi.Lo>>(n-64) | u.Lo.Lo.Hi.Hi<<(128-n),
					},
					Hi: Uint128{
						Lo: u.Lo.Lo.Hi.Hi>>(n-64) | u.Lo.Hi.Lo.Lo<<(128-n),
						Hi: u.Lo.Hi.Lo.Lo>>(n-64) | u.Lo.Hi.Lo.Hi<<(128-n),
					},
				},
				Hi: Uint256{
					Lo: Uint128{
						Lo: u.Lo.Hi.Lo.Hi>>(n-64) | u.Lo.Hi.Hi.Lo<<(128-n),
						Hi: u.Lo.Hi.Hi.Lo>>(n-64) | u.Lo.Hi.Hi.Hi<<(128-n),
					},
					Hi: Uint128{
						Lo: u.Lo.Hi.Hi.Hi>>(n-64) | u.Hi.Lo.Lo.Lo<<(128-n),
						Hi: u.Hi.Lo.Lo.Lo>>(n-64) | u.Hi.Lo.Lo.Hi<<(128-n),
					},
				},
			},
			Hi: Uint512{
				Lo: Uint256{
					Lo: Uint128{
						Lo: u.Hi.Lo.Lo.Hi>>(n-64) | u.Hi.Lo.Hi.Lo<<(128-n),
						Hi: u.Hi.Lo.Hi.Lo>>(n-64) | u.Hi.Lo.Hi.Hi<<(128-n),
					},
					Hi: Uint128{
						Lo: u.Hi.Lo.Hi.Hi>>(n-64) | u.Hi.Hi.Lo.Lo<<(128-n),
						Hi: u.Hi.Hi.Lo.Lo>>(n-64) | u.Hi.Hi.Lo.Hi<<(128-n),
					},
				},
				Hi: Uint256{
					Lo: Uint128{
						Lo: u.Hi.Hi.Lo.Hi>>(n-64) | u.Hi.Hi.Hi.Lo<<(128-n),
						Hi: u.Hi.Hi.Hi.Lo>>(n-64) | u.Hi.Hi.Hi.Hi<<(128-n),
					},
					Hi: Uint128{
						Lo: u.Hi.Hi.Hi.Hi >> (n - 64),
						Hi: 0,
					},
				},
			},
		}
	}

	return Uint1024{
		Lo: Uint512{
			Lo: Uint256{
				Lo: Uint128{
					Lo: u.Lo.Lo.Lo.Lo>>n | u.Lo.Lo.Lo.Hi<<(64-n),
					Hi: u.Lo.Lo.Lo.Hi>>n | u.Lo.Lo.Hi.Lo<<(64-n),
				},
				Hi: Uint128{
					Lo: u.Lo.Lo.Hi.Lo>>n | u.Lo.Lo.Hi.Hi<<(64-n),
					Hi: u.Lo.Lo.Hi.Hi>>n | u.Lo.Hi.Lo.Lo<<(64-n),
				},
			},
			Hi: Uint256{
				Lo: Uint128{
					Lo: u.Lo.Hi.Lo.Lo>>n | u.Lo.Hi.Lo.Hi<<(64-n),
					Hi: u.Lo.Hi.Lo.Hi>>n | u.Lo.Hi.Hi.Lo<<(64-n),
				},
				Hi: Uint128{
					Lo: u.Lo.Hi.Hi.Lo>>n | u.Lo.Hi.Hi.Hi<<(64-n),
					Hi: u.Lo.Hi.Hi.Hi>>n | u.Hi.Lo.Lo.Lo<<(64-n),
				},
			},
		},
		Hi: Uint512{
			Lo: Uint256{
				Lo: Uint128{
					Lo: u.Hi.Lo.Lo.Lo>>n | u.Hi.Lo.Lo.Hi<<(64-n),
					Hi: u.Hi.Lo.Lo.Hi>>n | u.Hi.Lo.Hi.Lo<<(64-n),
				},
				Hi: Uint128{
					Lo: u.Hi.Lo.Hi.Lo>>n | u.Hi.Lo.Hi.Hi<<(64-n),
					Hi: u.Hi.Lo.Hi.Hi>>n | u.Hi.Hi.Lo.Lo<<(64-n),
				},
			},
			Hi: Uint256{
				Lo: Uint128{
					Lo: u.Hi.Hi.Lo.Lo>>n | u.Hi.Hi.Lo.Hi<<(64-n),
					Hi: u.Hi.Hi.Lo.Hi>>n | u.Hi.Hi.Hi.Lo<<(64-n),
				},
				Hi: Uint128{
					Lo: u.Hi.Hi.Hi.Lo>>n | u.Hi.Hi.Hi.Hi<<(64-n),
					Hi: u.Hi.Hi.Hi.Hi >> n,
				},
			},
		},
	}
}

func (u Uint1024) BitLen() int {
	if !u.Hi.IsZero() {
		return 512 + u.Hi.BitLen()
	}
	return u.Lo.BitLen()
}

func (u Uint1024) LeadingZeros() int {
	if !u.Hi.IsZero() {
		return u.Hi.LeadingZeros()
	}
	return 512 + u.Lo.LeadingZeros()
}

func (u Uint1024) TrailingZeros() int {
	if !u.Lo.IsZero() {
		return u.Lo.TrailingZeros()
	}
	return 512 + u.Hi.TrailingZeros()
}

func (u Uint1024) OnesCount() int {
	return u.Lo.OnesCount() + u.Hi.OnesCount()
}

func (u Uint1024) Reverse() Uint1024 {
	return Uint1024{
		Lo: u.Hi.Reverse(),
		Hi: u.Lo.Reverse(),
	}
}

func (u Uint1024) ReverseBytes() Uint1024 {
	return Uint1024{
		Lo: u.Hi.ReverseBytes(),
		Hi: u.Lo.ReverseBytes(),
	}
}

func (u Uint1024) Bit(n int) bool {
	if n < 0 || n >= bitCount {
		return false
	}
	word := n / 64
	bit := uint(n % 64)
	return ((*(*[uint64Count]uint64)(unsafe.Pointer(&u)))[word]>>bit)&1 == 1
}
