package resp

import (
	"bytes"
	"fmt"
	"strconv"
)

func Marshal(v interface{}) ([]byte, error) {
	switch v := v.(type) {
	case string:
		return []byte("+" + v + "\r\n"), nil
	case error:
		return []byte("-" + v.Error() + "\r\n"), nil
	case int:
		return []byte(":" + strconv.Itoa(v) + "\r\n"), nil
	case []byte:
		if v == nil {
			return []byte("$-1\r\n"), nil
		}

		var buf bytes.Buffer
		buf.WriteString("$" + strconv.Itoa(len(v)) + "\r\n")
		buf.Write(v)
		buf.WriteString("\r\n")
		return buf.Bytes(), nil
	case []interface{}:
		if v == nil {
			return []byte("*-1\r\n"), nil
		}
		if len(v) == 0 {
			return []byte("*0\r\n"), nil
		}

		var buf bytes.Buffer
		buf.WriteString("*" + strconv.Itoa(len(v)) + "\r\n")

		for _, vv := range v {
			b, err := Marshal(vv)
			if err != nil {
				return nil, err
			}
			buf.Write(b)
		}

		return buf.Bytes(), nil
	case nil:
		return []byte("$-1\r\n"), nil
	default:
		return nil, fmt.Errorf("unsupported type: %T", v)
	}
}
