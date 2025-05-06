package uint1024

import "github.com/piliming/bigz/uint256"

func Zero() Uint512 {
	return From64(0)
}

func One() Uint512 {
	return From64(1)
}

func Max() Uint512 {
	return Uint512{
		Lo: uint256.Max(),
		Hi: uint256.Max(),
	}
}

type Uint512 struct {
	Lo uint256.Uint256
	Hi uint256.Uint256
}

func From256(v uint256.Uint256) Uint512 {
	return Uint512{
		Lo: v,
	}
}

func From64(v uint64) Uint512 {
	return From256(uint256.From64(v))
}
