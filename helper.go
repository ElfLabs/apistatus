package apistatus

import (
	"bytes"
	"encoding/json"
)

// WriteData 写入数据
func WriteData(buf *bytes.Buffer, data any) error {
	switch v := data.(type) {
	case nil:
		return nil
	case []byte:
		_, err := buf.Write(v)
		return err
	case string:
		_, err := buf.WriteString(v)
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
