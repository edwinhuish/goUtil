package TypeTools

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"time"
)

func ByteToInt32(data []byte) int32 {
	buf := bytes.NewBuffer(data)
	var x int32
	binary.Read(buf, order, &x)
	return x
}

func ByteToInt64(data []byte) int64 {
	buf := bytes.NewBuffer(data)
	var x int64
	binary.Read(buf, order, &x)
	return x
}
func ByteToUInt64(data []byte) uint64 {
	buf := bytes.NewBuffer(data)
	var x uint64
	binary.Read(buf, order, &x)
	return x
}
func FormatJpeg(data []byte) string {
	return "data:image/jpg;base64," + base64.StdEncoding.EncodeToString(data)
}
func FormatPng(data []byte) string {
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(data)
}
func ByteToTime(data []byte) time.Time {
	d := ByteToInt64(data)
	return time.Unix(0, d)
}
func TimeToByte(data time.Time) []byte {
	return Int64ToByte(data.UnixNano())
}
