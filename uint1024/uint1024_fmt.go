package uint1024

import (
	"fmt"
	"github.com/piliming/bigz/uint512"
	"math/big"
)

// FromString parses input string as a Uint256 value.
func FromString(s string) (Uint1024, error) {
	var u Uint1024
	_, err := fmt.Sscan(s, &u)
	return u, err
}

func (u Uint1024) String() string {
	if u.Hi.IsZero() {
		return u.Lo.String()
	}

	buf := []byte("000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
	for i := len(buf); ; i -= 19 {
		q, r := u.QuoRem64(1e19)
		var n int
		for ; r != 0; r /= 10 {
			n++
			buf[i-n] += byte(r % 10)
		}

		if q.IsZero() {
			return string(buf[i-n:])
		}
		u = q
	}
}

// Format does custom formatting of 1024-bit value.
func (u Uint1024) Format(s fmt.State, ch rune) {
	u.Big().Format(s, ch) // via big.Int, unefficient! consider to optimize
}

// Scan implements fmt.Scanner.
func (u *Uint1024) Scan(s fmt.ScanState, ch rune) error {
	i := new(big.Int) // via big.Int, unefficient! consider to optimize
	if err := i.Scan(s, ch); err != nil {
		return err
	}

	v, ok := FromBigEx(i)
	if !ok {
		return fmt.Errorf("out of 1024-bit range")
	}

	*u = v
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface.
func (u Uint1024) MarshalText() (text []byte, err error) {
	return u.Big().MarshalText() // via big.Int, unefficient! consider to optimize
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (u *Uint1024) UnmarshalText(text []byte) error {
	// via big.Int, unefficient! consider to optimize
	i := new(big.Int)
	if err := i.UnmarshalText(text); err != nil {
		return err
	}
	v, ok := FromBigEx(i)
	if !ok {
		return fmt.Errorf("%q overflows 1024-bit integer", text)
	}
	*u = v
	return nil
}

//// StoreLittleEndian stores 256-bit value in byte slice in little-endian byte order.
//// It panics if byte slice length is less than 32.
//func StoreLittleEndian(b []byte, u Uint256) {
//	uint128.StoreLittleEndian(b[:16], u.Lo)
//	uint128.StoreLittleEndian(b[16:], u.Hi)
//}
//
//// StoreBigEndian stores 256-bit value in byte slice in big-endian byte order.
//// It panics if byte slice length is less than 32.
//func StoreBigEndian(b []byte, u Uint256) {
//	uint128.StoreBigEndian(b[:16], u.Hi)
//	uint128.StoreBigEndian(b[16:], u.Lo)
//}

// LoadLittleEndian loads 1024-bit value from byte slice in little-endian byte order.
// It panics if byte slice length is less than 32.
func LoadLittleEndian(b [128]byte) Uint1024 {
	return Uint1024{
		Lo: uint512.LoadLittleEndian([64]byte(b[:64])),
		Hi: uint512.LoadLittleEndian([64]byte(b[64:])),
	}
}

//// LoadBigEndian loads 256-bit value from byte slice in big-endian byte order.
//// It panics if byte slice length is less than 32.
//func LoadBigEndian(b []byte) Uint256 {
//	return Uint256{
//		Lo: uint128.LoadBigEndian(b[16:]),
//		Hi: uint128.LoadBigEndian(b[:16]),
//	}
//}
