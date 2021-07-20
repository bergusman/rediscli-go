package resp

import (
	"errors"
	"testing"
)

func TestMarshalSimpleString(t *testing.T) {
	r, err := Marshal("OK")
	if err != nil {
		t.Fatal(err)
	}
	if string(r) != "+OK\r\n" {
		t.Errorf("got %q", r)
	}
}

func TestMarshalEmptySimpleString(t *testing.T) {
	r, err := Marshal("")
	if err != nil {
		t.Fatal(err)
	}
	if string(r) != "+\r\n" {
		t.Errorf("got %q", r)
	}
}

func TestMarshalError(t *testing.T) {
	r, err := Marshal(errors.New("Error message"))
	if err != nil {
		t.Fatal(err)
	}
	if string(r) != "-Error message\r\n" {
		t.Errorf("got %q", r)
	}
}

func TestMarshalIntegers(t *testing.T) {
	r, err := Marshal(1000)
	if err != nil {
		t.Fatal(err)
	}
	if string(r) != ":1000\r\n" {
		t.Errorf("got %q", r)
	}

	r, err = Marshal(0)
	if err != nil {
		t.Fatal(err)
	}
	if string(r) != ":0\r\n" {
		t.Errorf("got %q", r)
	}

	r, err = Marshal(-1)
	if err != nil {
		t.Fatal(err)
	}
	if string(r) != ":-1\r\n" {
		t.Errorf("got %q", r)
	}
}

func TestMarshalBulkString(t *testing.T) {
	r, err := Marshal([]byte("foobar"))
	if err != nil {
		t.Fatal(err)
	}
	if string(r) != "$6\r\nfoobar\r\n" {
		t.Errorf("got %q", r)
	}
}

func TestMarshalEmptyBulkString(t *testing.T) {
	r, err := Marshal([]byte(""))
	if err != nil {
		t.Fatal(err)
	}
	if string(r) != "$0\r\n\r\n" {
		t.Errorf("got %q", r)
	}
}

func TestMarshalNullBulkString(t *testing.T) {
	r, err := Marshal([]byte(nil))
	if err != nil {
		t.Fatal(err)
	}
	if string(r) != "$-1\r\n" {
		t.Errorf("got %q", r)
	}
}

func TestMarshalNull(t *testing.T) {
	r, err := Marshal(nil)
	if err != nil {
		t.Fatal(err)
	}
	if string(r) != "$-1\r\n" {
		t.Errorf("got %q", r)
	}
}

func TestMarshalArray(t *testing.T) {
	r, err := Marshal([]interface{}{
		[]byte("foo"),
		[]byte("bar"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if string(r) != "*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n" {
		t.Errorf("got %q", r)
	}

	r, err = Marshal([]interface{}{
		1, 2, 3,
	})
	if err != nil {
		t.Fatal(err)
	}
	if string(r) != "*3\r\n:1\r\n:2\r\n:3\r\n" {
		t.Errorf("got %q", r)
	}
}

func TestMarshalEmptyArray(t *testing.T) {
	r, err := Marshal(make([]interface{}, 0))
	if err != nil {
		t.Fatal(err)
	}
	if string(r) != "*0\r\n" {
		t.Errorf("got %q", r)
	}
}

func TestMarshalNullArray(t *testing.T) {
	r, err := Marshal([]interface{}(nil))
	if err != nil {
		t.Fatal(err)
	}
	if string(r) != "*-1\r\n" {
		t.Errorf("got %q", r)
	}
}

func TestMarshalFailure(t *testing.T) {
	_, err := Marshal(struct{}{})
	if err == nil {
		t.Error("awaiting error")
	}
}

func BenchmarkMarshalString(b *testing.B) {
	v := "OK"
	for i := 0; i < b.N; i++ {
		Marshal(v)
	}
}

func BenchmarkMarshalIntegerArray(b *testing.B) {
	v := []interface{}{
		1, 2, 3,
	}
	for i := 0; i < b.N; i++ {
		Marshal(v)
	}
}

func BenchmarkMarshalBulkStringArray(b *testing.B) {
	v := []interface{}{
		[]byte("foo"),
		[]byte("bar"),
	}
	for i := 0; i < b.N; i++ {
		Marshal(v)
	}
}
