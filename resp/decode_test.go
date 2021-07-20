package resp

import (
	"reflect"
	"testing"
)

func TestUnmarshalSimpleString(t *testing.T) {
	v, err := Unmarshal([]byte("+OK\r\n"))
	if err != nil {
		t.Fatal(err)
	}
	s, ok := v.(string)
	if !ok {
		t.Error("awaiting string type")
	}
	if s != "OK" {
		t.Errorf("unmarshaled %q; want \"OK\"", s)
	}
}

func TestUnmarshalEmptySimpleString(t *testing.T) {
	v, err := Unmarshal([]byte("+\r\n"))
	if err != nil {
		t.Fatal(err)
	}
	s, ok := v.(string)
	if !ok {
		t.Error("awaiting string type")
	}
	if s != "" {
		t.Errorf("unmarshaled %q; want \"\"", s)
	}
}

func TestUnmarshalError(t *testing.T) {
	v, err := Unmarshal([]byte("-Error message\r\n"))
	if err != nil {
		t.Fatal(err)
	}
	e, ok := v.(error)
	if !ok {
		t.Error("awaiting error type")
	}
	if e.Error() != "Error message" {
		t.Errorf("unmarshaled %q; want \"Error message\"", e)
	}
}

func TestUnmarshalInteger(t *testing.T) {
	v, err := Unmarshal([]byte(":1000\r\n"))
	if err != nil {
		t.Fatal(err)
	}
	i, ok := v.(int)
	if !ok {
		t.Error("awaiting int type")
	}
	if i != 1000 {
		t.Errorf("unmarshaled %v; want 1000", i)
	}
}

func TestUnmarshalBulkString(t *testing.T) {
	v, err := Unmarshal([]byte("$6\r\nfoobar\r\n"))
	if err != nil {
		t.Fatal(err)
	}
	b, ok := v.([]byte)
	if !ok {
		t.Error("awaiting byte slice type")
	}
	if string(b) != "foobar" {
		t.Errorf("unmarshaled %q; want \"foobar\"", b)
	}
}

func TestUnmarshalEmptyBulkString(t *testing.T) {
	v, err := Unmarshal([]byte("$0\r\n\r\n"))
	if err != nil {
		t.Fatal(err)
	}
	b, ok := v.([]byte)
	if !ok {
		t.Error("awaiting byte slice type")
	}
	if string(b) != "" {
		t.Errorf("unmarshaled %q; want \"\"", b)
	}
}

func TestUnmarshalNullBulkString(t *testing.T) {
	v, err := Unmarshal([]byte("$-1\r\n"))
	if err != nil {
		t.Fatal(err)
	}
	if v != nil {
		t.Errorf("unmarshaled %v; want nil", v)
	}
}

func TestUnmarshalArray(t *testing.T) {
	v, err := Unmarshal([]byte("*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"))
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(v, []interface{}{
		[]byte("foo"),
		[]byte("bar"),
	}) {
		t.Errorf("unmarshaled %v; want other", v)
	}
}

func TestUnmarshalEmptyArray(t *testing.T) {
	v, err := Unmarshal([]byte("*0\r\n"))
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(v, []interface{}{}) {
		t.Errorf("unmarshaled %v; want empty slice", v)
	}
}

func TestUnmarshalNullArray(t *testing.T) {
	v, err := Unmarshal([]byte("*-1\r\n"))
	if err != nil {
		t.Fatal(err)
	}
	a, ok := v.([]interface{})
	if !ok {
		t.Error("awaiting interface slice type")
	}
	if a != nil {
		t.Errorf("unmarshaled %v; want nil slice", v)
	}
}
