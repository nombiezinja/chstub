package entities

import (
	"github.com/francoispqt/gojay"
	"github.com/nombiezinja/chstub/utils"
)

func (p *Payload) MarshalJSONObject(enc *gojay.Encoder) {
	enc.AddIntKey("ID", p.ID)
	enc.AddStringKey("Colour", p.Colour)
}

func (p *Payload) IsNil() bool {
	return p == nil
}

func GojayMarshal(p *Payload) []byte {
	marshalled, err := gojay.MarshalJSONObject(p)

	utils.FailOnError(err, "Marshal json failed")

	return marshalled
}
