// + build ignore

package uint1024

import (
	"fmt"
	"math/big"
	"math/bits"
	"math/rand"
	"slices"
	"testing"
)

func rand1024slice(count int) []Uint1024 {
	out := make([]Uint1024, count)
	for i := range out {
		out[i] = rand1024()
	}
	return out
}

func rand1024() Uint1024 {
	buf := [128]byte{}

	for i := 0; i < 128; i++ {
		buf[i] = byte(rand.Uint64() % 256)
	}

	return LoadLittleEndian(buf)
}

func assertString(t *testing.T, got, want string, prefix string) {
	if got != want {
		t.Fatalf("%s: \nwant=> %s \ngot => %s", prefix, want, got)
	}
}

func assertBool(t *testing.T, got, want bool, prefix string) {
	if got != want {
		t.Fatalf("%s: \nwant=> %v \ngot => %v", prefix, want, got)
	}
}

func assertInt(t *testing.T, got, want int, prefix string) {
	if got != want {
		t.Fatalf("%s: \nwant=> %v \ngot => %v", prefix, want, got)
	}
}

const loopTimes = 1000000

var maxBigIntUint1024 = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 1024), big.NewInt(1))

func TestUint1024_Add(t *testing.T) {
	for i := 0; i < loopTimes; i++ {

		val := rand1024()
		bigVal := val.Big()

		assertString(t, val.String(), bigVal.String(), "big")

		addVal := rand1024()
		addBigVal := addVal.Big()
		assertString(t, addVal.String(), addBigVal.String(), "big1")

		val = val.Add(addVal)
		bigVal = bigVal.Add(bigVal, addBigVal)
		bigVal = bigVal.And(bigVal, maxBigIntUint1024)
		assertString(t, val.String(), bigVal.String(), "add")
	}
}

func TestUint1024_Sub(t *testing.T) {
	for i := 0; i < loopTimes; i++ {
		val := rand1024()
		bigVal := val.Big()

		bigVal.Bytes()

		subVal := rand1024()
		subBigVal := subVal.Big()
		assertString(t, subVal.String(), subBigVal.String(), "big2")

		val = val.Sub(subVal)
		bigVal = bigVal.Sub(bigVal, subBigVal)
		bigVal = bigVal.And(bigVal, maxBigIntUint1024)
		assertString(t, val.String(), bigVal.String(), "sub")
	}
}

func TestUint1024_Mul(t *testing.T) {
	for i := 0; i < loopTimes; i++ {
		val := rand1024()
		bigVal := val.Big()

		mulVal := rand1024()
		mulBigVal := mulVal.Big()

		val = val.Mul(mulVal)

		bigVal = bigVal.Mul(bigVal, mulBigVal)
		bigVal = bigVal.And(bigVal, maxBigIntUint1024)

		assertString(t, val.String(), bigVal.String(), "mul")
	}
}

func TestUint1024_QRem64(t *testing.T) {
	for i := 0; i < loopTimes; i++ {

		val := rand1024()
		bigVal := val.Big()

		//assertString(t, val.String(), bigVal.String(), "big")

		divVal := uint64(rand.Uint32())
		divBigVal := big.NewInt(int64(divVal))
		assertString(t, fmt.Sprintf("%d", divVal), divBigVal.String(), "big1")

		qVal, rVal := val.QuoRem64(divVal)
		qBigVal, rBigVal := bigVal.QuoRem(bigVal, divBigVal, big.NewInt(0))
		//bigVal = bigVal.And(bigVal, maxBigIntUint1024)

		assertString(t, fmt.Sprintf("%d", rVal), rBigVal.String(), "r")

		assertString(t, qVal.String(), qBigVal.String(), "q")
	}
}

func TestUint1024_Div(t *testing.T) {
	for i := 0; i < loopTimes; i++ {
		val := rand1024()

		oriVal := val.String()
		bigVal := val.Big()
		oriBigVal := bigVal.String()

		divVal := rand1024()
		divBigVal := divVal.Big()
		val = val.Div(divVal)
		bigVal = bigVal.Div(bigVal, divBigVal)
		bigVal = bigVal.And(bigVal, maxBigIntUint1024)

		got, want := val.String(), bigVal.String()

		if got != want {
			fmt.Println("oriVal:   ", oriVal)
			fmt.Println("oriBigVal:", oriBigVal)

			fmt.Println("divVal:   ", divVal)
			fmt.Println("divBigVal:", divBigVal)

			fmt.Println("val:      ", val)
			fmt.Println("bigVal:   ", bigVal)

		}

		assertString(t, got, want, "div")
	}
}

func TestUint1024_Lsh(t *testing.T) {
	for i := 0; i < loopTimes; i++ {
		val := rand1024()
		bigVal := val.Big()

		shift := rand.Uint64() % 1030

		shiftVal := val.Lsh(uint(shift))
		bigShiftVal := new(big.Int).Lsh(bigVal, uint(shift))
		bigShiftVal = bigShiftVal.And(bigShiftVal, maxBigIntUint1024)
		got, want := shiftVal.String(), bigShiftVal.String()

		assertString(t, got, want, fmt.Sprintf("Lsh-%d", shift))
	}
}

func TestUint1024_Rsh(t *testing.T) {
	for i := 0; i < loopTimes; i++ {
		val := rand1024()
		bigVal := val.Big()

		//shift := rand.Uint64() % (uint64(len(val)) * 64)
		shift := rand.Uint64() % 1030

		shiftVal := val.Rsh(uint(shift))
		bigShiftVal := new(big.Int).Rsh(bigVal, uint(shift))
		bigShiftVal = bigShiftVal.And(bigShiftVal, maxBigIntUint1024)

		assertString(t, shiftVal.String(), bigShiftVal.String(), fmt.Sprintf("Rsh-%d", shift))
	}
}

func TestUint1024_Bits(t *testing.T) {
	for i := 0; i < loopTimes; i++ {
		val := rand1024()
		bigVal := val.Big()

		for i := 0; i < 1030; i++ {
			assertBool(t, val.Bit(i), bigVal.Bit(i) == 1, "Bit")
		}
	}
}

func TestUint1024_LeadingZeros(t *testing.T) {
	for i := 0; i < loopTimes; i++ {
		val := rand1024()

		bigVal := val.Big()

		assertString(t, val.String(), bigVal.String(), "LeadingZeros_Big")

		valLeadingZeros := val.LeadingZeros()

		bigBytes := bigVal.Bytes()

		diff := 128 - len(bigBytes)

		if diff > 0 {
			diffBytes := make([]byte, diff)
			bigBytes = append(diffBytes, bigBytes...)
		}

		//slices.Reverse(bigBytes)
		bigValLeadingZeros := 0

		for _, bigByte := range bigBytes {
			if bigByte == 0 {
				bigValLeadingZeros += 8
				continue
			}
			bigValLeadingZeros += bits.LeadingZeros8(bigByte)
			break
		}
		assertInt(t, valLeadingZeros, bigValLeadingZeros, "LeadingZeros")

	}
}

func TestUint1024_TrailingZeros(t *testing.T) {
	for i := 0; i < loopTimes; i++ {
		val := rand1024()

		bigVal := val.Big()

		assertString(t, val.String(), bigVal.String(), "TrailingZeros_Big")

		valTrailingZeros := val.TrailingZeros()

		bigBytes := bigVal.Bytes()

		diff := 128 - len(bigBytes)

		if diff > 0 {
			diffBytes := make([]byte, diff)
			bigBytes = append(diffBytes, bigBytes...)
		}

		slices.Reverse(bigBytes)
		bigValTrailingZeros := 0

		for _, bigByte := range bigBytes {
			if bigByte == 0 {
				bigValTrailingZeros += 8
				continue
			}
			bigValTrailingZeros += bits.TrailingZeros8(bigByte)
			break
		}
		assertInt(t, valTrailingZeros, bigValTrailingZeros, "TrailingZeros")

	}
}

func TestUint1024_FromBig(t *testing.T) {
	for i := 0; i < loopTimes; i++ {
		val := rand1024()
		fromBigVal := FromBig(val.Big())

		assertString(t, val.String(), fromBigVal.String(), "FromBig")

	}
}
