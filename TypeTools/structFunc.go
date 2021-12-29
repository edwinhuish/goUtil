package TypeTools

import (
	"encoding/json"
	"bytes"
	"fmt"
	"strconv"
)

func OutJson(val interface{}) string {
	b, _ := json.Marshal(val)
	return string(b)
}
func OutStructure(val []byte, data interface{}) error {
	decode := json.NewDecoder(bytes.NewReader(val))
	decode.UseNumber()
	return decode.Decode(&data)
}

//int转float64，分转元
func IntToFloat64(i int) (result float64) {
	result, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", float64(i)/float64(100)), 64)
	return
}
//int64转float64，分转元
func Int64ToFloat64(i int64) (result float64) {
	result, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", float64(i)/float64(100)), 64)
	return
}

//float64转int，元转分
func Float64ToInt(f float64) (result int) {
	return int(f * 100)
}
//float64转int，元转分
func Float64ToInt64(f float64) (result int64) {
	return int64(f * 100)
}

//string转float64
func StringToFloat64(s string) (result float64) {
	result, _ = strconv.ParseFloat(s, 64)
	return
}

//float64转string
func Float64ToString(s float64) (result string) {
	return strconv.FormatFloat(s, 'E', -1, 64)
}