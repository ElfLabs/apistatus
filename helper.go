package apistatus

import (
	"bytes"
	"encoding/json"
	"strconv"
)

// WriteData 写入数据
func WriteData(buf *bytes.Buffer, data any) (err error) {
	switch v := data.(type) {
	case nil:
		_, err = buf.WriteString("null")
		return
	case []byte:
		if json.Valid(v) {
			_, err = buf.Write(v)
		} else {
			buf.WriteByte('"')
			_, err = buf.Write(v)
			buf.WriteByte('"')
		}
		return err
	case string:
		_, err = buf.WriteString(strconv.Quote(v))
		return err
	default:
		dataBytes, err := json.Marshal(data)
		if err != nil {
			return err
		}
		_, err = buf.Write(dataBytes)
		return err
	}
}
