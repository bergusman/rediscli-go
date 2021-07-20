package resp

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

var streamTest = []interface{}{
	"OK",
	errors.New("error message"),
	1000,
	0,
	-1,
	[]byte("foobar"),
	[]byte(""),
	[]byte(nil),
	nil,
	[]interface{}{
		[]byte("foo"),
		[]byte("bar"),
	},
	[]interface{}{
		1, 2, 3,
	},
	[]interface{}{},
	[]interface{}(nil),
}

var streamEncoded = []string{
	"+OK\r\n",
	"-error message\r\n",
	":1000\r\n",
	":0\r\n",
	":-1\r\n",
	"$6\r\nfoobar\r\n",
	"$0\r\n\r\n",
	"$-1\r\n",
	"$-1\r\n",
	"*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n",
	"*3\r\n:1\r\n:2\r\n:3\r\n",
	"*0\r\n",
	"*-1\r\n",
}

func TestEncoder(t *testing.T) {
	for i := 0; i <= len(streamTest); i++ {
		var buf bytes.Buffer
		enc := NewEncoder(&buf)
		for j, v := range streamTest[0:i] {
			if err := enc.Encode(v); err != nil {
				t.Fatalf("encode #%d: %v", j, err)
			}
		}
		if have, want := buf.String(), nlines(streamEncoded, i); have != want {
			t.Errorf("encoding %d items: mismatch", i)
			diff(t, []byte(have), []byte(want))
			break
		}
	}
}

func nlines(ss []string, n int) string {
	return strings.Join(ss[0:n], "")
}

func diff(t *testing.T, a, b []byte) {
	for i := 0; ; i++ {
		if i >= len(a) || i >= len(b) || a[i] != b[i] {
			j := i - 10
			if j < 0 {
				j = 0
			}
			t.Errorf("diverge at %d: «%s» vs «%s»", i, trim(a[j:]), trim(b[j:]))
			return
		}
	}
}

func trim(b []byte) []byte {
	if len(b) > 20 {
		return b[0:20]
	}
	return b
}
