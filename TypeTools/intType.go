package TypeTools

import (
	"encoding/base64"
	"encoding/binary"
	"math"
	"strconv"
)

var order binary.ByteOrder = binary.BigEndian

func IntToByte(num int) []byte {
	return UIntToByte(uint32(num))
}
func UIntToByte(num uint32) []byte {
	bs := make([]byte, 4)
	order.PutUint32(bs, num)
	return bs
}

func Int64ToByte(num int64) []byte {
	return UInt64ToByte(uint64(num))
}
func UInt64ToByte(num uint64) []byte {
	bs := make([]byte, 8)
	order.PutUint64(bs, num)
	return bs
}
func Base64Int64(num int64) string {
	return base64.StdEncoding.EncodeToString(Int64ToByte(num))
}
func EncodeInt(num uint32) string {
	a := strconv.FormatInt(int64(num), 36)
	for len(a) < 7 {
		a += "X"
	}
	return a
}
func Float64ToByte(num float64) []byte {
	return UInt64ToByte(math.Float64bits(num))
}
func ByteToFloat64(arr []byte) float64 {
	return math.Float64frombits(ByteToUInt64(arr))
}
