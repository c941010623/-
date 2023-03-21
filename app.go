package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
)

// 編碼
func jsonToMessagePack(jsonBytes []byte) ([]byte, error) {
	// 將 JSON 資料轉換為 Go 中的 map[string]interface{} 型別
	var jsonObj map[string]interface{}
	err := json.Unmarshal(jsonBytes, &jsonObj)
	if err != nil {
		return nil, err
	}

	// 將 Go 中的資料結構轉換為 MessagePack 格式
	var buf bytes.Buffer
	err = writeValue(&buf, jsonObj)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func writeValue(buf *bytes.Buffer, val interface{}) error {
	switch v := val.(type) {
	case bool:
		if v {
			buf.WriteByte(0xc3) // true
		} else {
			buf.WriteByte(0xc2) // false
		}
	case int:
		if v >= 0 && v <= 127 {
			buf.WriteByte(byte(v))
		} else if v >= -32 && v <= -1 {
			buf.WriteByte(byte(0xe0 | (v & 0x1f)))
		} else if v >= -128 && v <= 127 {
			buf.WriteByte(0xd0) // int8
			binary.Write(buf, binary.BigEndian, int8(v))
		} else if v >= -32768 && v <= 32767 {
			buf.WriteByte(0xd1) // int16
			binary.Write(buf, binary.BigEndian, int16(v))
		} else if v >= -2147483648 && v <= 2147483647 {
			buf.WriteByte(0xd2) // int32
			binary.Write(buf, binary.BigEndian, int32(v))
		} else {
			buf.WriteByte(0xd3) // int64
			binary.Write(buf, binary.BigEndian, int64(v))
		}
	case float64:
		buf.WriteByte(0xcb) // float64
		binary.Write(buf, binary.BigEndian, v)
	case string:
		b := []byte(v)
		l := len(b)
		if l <= 31 {
			buf.WriteByte(byte(0xa0 | l))
		} else if l <= 255 {
			buf.WriteByte(0xd9) // str8
			buf.WriteByte(byte(l))
		} else if l <= 65535 {
			buf.WriteByte(0xda) // str16
			binary.Write(buf, binary.BigEndian, uint16(l))
		} else {
			buf.WriteByte(0xdb) // str32
			binary.Write(buf, binary.BigEndian, uint32(l))
		}
		buf.Write(b)
	case []interface{}:
		l := len(v)
		if l <= 15 {
			buf.WriteByte(byte(0x90 | l))
		} else if l <= 65535 {
			buf.WriteByte(0xdc) // array16
			binary.Write(buf, binary.BigEndian, uint16(l))
		} else {
			buf.WriteByte(0xdd) // array32
			binary.Write(buf, binary.BigEndian, uint32(l))
		}
		for _, e := range v {
			err := writeValue(buf, e)
			if err != nil {
				return err
			}
		}
	case map[string]interface{}:
		l := len(v)
		if l <= 15 {
			buf.WriteByte(byte(0x80 | l))
		} else if l <= 65535 {
			buf.WriteByte(0xde) // map16
			binary.Write(buf, binary.BigEndian, uint16(l))
		} else {
			buf.WriteByte(0xdf) // map32
			binary.Write(buf, binary.BigEndian, uint32(l))
		}
		for k, e := range v {
			err := writeValue(buf, k)
			if err != nil {
				return err
			}
			err = writeValue(buf, e)
			if err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("不支持的類型: %T", v)
	}
	return nil
}

func main() {
	jsonStr := `{"name":"Alice","age":20,"score":[80,85,90]}`
	jsonBytes := []byte(jsonStr)

	msgPackBytes, err := jsonToMessagePack(jsonBytes)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("JSON: %s\n", jsonStr)
	fmt.Printf("MessagePack: %x\n", msgPackBytes)
}
