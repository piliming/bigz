package uint512

import (
	"fmt"
	"github.com/Pilatuz/bigz/uint256"
	"math/big"
)

// FromString parses input string as a Uint256 value.
func FromString(s string) (Uint512, error) {
	var u Uint512
	_, err := fmt.Sscan(s, &u)
	return u, err
}

func (u Uint512) String() string {
	if u.Hi.IsZero() {
		return u.Lo.String()
	}

	// log10(2^512) ≈ 154, so a buffer of 156 is enough
	buf := []byte("000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000") // len=156
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

// Format does custom formatting of 256-bit value.
func (u Uint512) Format(s fmt.State, ch rune) {
	u.Big().Format(s, ch) // via big.Int, unefficient! consider to optimize
}

// Scan implements fmt.Scanner.
func (u *Uint512) Scan(s fmt.ScanState, ch rune) error {
	i := new(big.Int) // via big.Int, unefficient! consider to optimize
	if err := i.Scan(s, ch); err != nil {
		return err
	}

	v, ok := FromBigEx(i)
	if !ok {
		return fmt.Errorf("out of 256-bit range")
	}

	*u = v
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface.
func (u Uint512) MarshalText() (text []byte, err error) {
	return u.Big().MarshalText() // via big.Int, unefficient! consider to optimize
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (u *Uint512) UnmarshalText(text []byte) error {
	// via big.Int, unefficient! consider to optimize
	i := new(big.Int)
	if err := i.UnmarshalText(text); err != nil {
		return err
	}
	v, ok := FromBigEx(i)
	if !ok {
		return fmt.Errorf("%q overflows 256-bit integer", text)
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

// LoadLittleEndian loads 256-bit value from byte slice in little-endian byte order.
// It panics if byte slice length is less than 32.
func LoadLittleEndian(b [64]byte) Uint512 {
	return Uint512{
		Lo: uint256.LoadLittleEndian(b[:32]),
		Hi: uint256.LoadLittleEndian(b[32:]),
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
