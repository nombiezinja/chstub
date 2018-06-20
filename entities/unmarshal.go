package entities

import (
	"github.com/francoispqt/gojay"
	u "github.com/nombiezinja/chstub/utils"
)

func (t *Payload) UnmarshalJSONObject(dec *gojay.Decoder, key string) error {
	switch key {
	case "ID":
		return dec.AddInt(&t.ID)
	case "Colour":
		return dec.AddString(&t.Colour)
	}
	return nil
}

func (t *Payload) NKeys() int {
	return 2
}

func GojayUnmarshal(data []byte) Payload {
	var p Payload
	err := gojay.UnmarshalJSONObject(data, &p)
	u.FailOnError(err, "Failed to unmarshall")
	return p
}
