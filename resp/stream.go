package resp

import (
	"bufio"
	"io"
)

type Encoder struct {
	w   io.Writer
	err error
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

func (enc *Encoder) Encode(v interface{}) error {
	if enc.err != nil {
		return enc.err
	}
	b, err := Marshal(v)
	if err != nil {
		return err
	}
	if _, err = enc.w.Write(b); err != nil {
		enc.err = err
	}
	return err
}

type Decoder struct {
	dec decoder
	err error
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{dec: decoder{r: bufio.NewReader(r)}}
}

func (dec *Decoder) Decode() (interface{}, error) {
	if dec.err != nil {
		return nil, dec.err
	}
	return dec.dec.decode()
}
