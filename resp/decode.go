package resp

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strconv"
)

func Unmarshal(data []byte) (interface{}, error) {
	dec := decoder{r: bufio.NewReader(bytes.NewReader(data))}
	return dec.decode()
}

type decoder struct {
	r *bufio.Reader
}

func (dec *decoder) decode() (interface{}, error) {
	c, err := dec.r.ReadByte()
	if err != nil {
		return nil, err
	}

	switch c {
	case '+':
		s, err := dec.readString()
		if err != nil {
			return nil, err
		}
		return s, nil

	case '-':
		s, err := dec.readString()
		if err != nil {
			return nil, err
		}
		return errors.New(s), nil

	case ':':
		i, err := dec.readInteger()
		if err != nil {
			return nil, err
		}
		return i, nil

	case '$':
		n, err := dec.readInteger()
		if err != nil {
			return nil, err
		}
		if n < -1 {
			return nil, errors.New("count less then -1")
		}
		if n == -1 {
			return nil, nil
		}
		bulk := make([]byte, n)
		if n > 0 {
			_, err = io.ReadFull(dec.r, bulk)
			if err != nil {
				return nil, err
			}
		}
		err = dec.readCRLF()
		if err != nil {
			return nil, err
		}
		return bulk, nil

	case '*':
		n, err := dec.readInteger()
		if err != nil {
			return nil, err
		}
		if n < -1 {
			return nil, errors.New("count less then -1")
		}
		if n == -1 {
			return []interface{}(nil), nil
		}
		if n == 0 {
			return make([]interface{}, 0), nil
		}
		var arr []interface{}
		for i := 0; i < n; i++ {
			v, err := dec.decode()
			if err != nil {
				return nil, err
			}
			arr = append(arr, v)
		}
		return arr, nil

	default:
		return nil, errors.New("unsupported type first byte " + quoteChar(c))
	}
}

func (dec *decoder) readString() (string, error) {
	s, err := dec.r.ReadString('\r')
	if err != nil {
		return "", err
	}
	err = dec.readLF()
	if err != nil {
		return "", err
	}

	s = s[:len(s)-1]
	return string(s), nil
}

func (dec *decoder) readInteger() (int, error) {
	s, err := dec.r.ReadString('\r')
	if err != nil {
		return 0, err
	}
	err = dec.readLF()
	if err != nil {
		return 0, err
	}

	s = s[:len(s)-1]
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}

	return i, nil
}

func (dec *decoder) readCRLF() error {
	b, err := dec.r.ReadByte()
	if err != nil {
		return err
	}
	if b != '\r' {
		return errors.New(`awaiting '\r' but read ` + quoteChar(b))
	}
	return dec.readLF()
}

func (dec *decoder) readLF() error {
	b, err := dec.r.ReadByte()
	if err != nil {
		return err
	}
	if b != '\n' {
		return errors.New(`awaiting '\n' but read ` + quoteChar(b))
	}
	return nil
}

func quoteChar(c byte) string {
	if c == '\'' {
		return `'\''`
	}
	if c == '"' {
		return `'"'`
	}
	s := strconv.Quote(string(c))
	return "'" + s[1:len(s)-1] + "'"
}
