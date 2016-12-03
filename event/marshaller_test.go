package event

import (
	"testing"
	"encoding/json"
)

type (
	MarshallEvent struct {
		Event `json:"-" id:"JsonEvent" v:"0"`
		Name string `json:"Name"`
		Surname string `json:"Surname"`
	}
)

func TestJsonMarshal(t *testing.T) {
	expectedOutput := "{\"Event\":{\"Name\":\"Shawn\",\"Surname\":\"Ritchie\"},\"MetaData\":{\"0\":1},\"Identifier\":\"JsonEvent\",\"Version\":0}"
	c := NewContainer(MarshallEvent{Name:"Shawn", Surname:"Ritchie"}, MetaData{ seqKey: uint64(1) })

	b, err := json.Marshal(c)
	if (err != nil) {
		t.Errorf(err.Error())
	}

	if (string(b) != expectedOutput) {
		t.Errorf("Expected output did not match %s ", string(b))
	}
}

func TestJsonUnmarshal(t *testing.T) {
	e := MarshallEvent{Name:"Shawn", Surname:"Ritchie"}
	s := uint64(1)
	c := NewContainer(e, MetaData{ seqKey: s })
	b, err := json.Marshal(c)
	if (err != nil) {
		t.Errorf(err.Error())
	}

	u := &struct {
		Event `json:"-" id:"JsonEvent" v:"0"`
		Name string `json:"Name"`
		Surname string `json:"Surname"`
	}{}

	var payload *Container = new(Container)
	payload.UnmarshalJSON(containeKeyDef, u, b)

	if e.Name != u.Name && e.Surname != u.Surname && payload.MetaData[seqKey] != s {
		t.Errorf("Error deserializing Values Name: %s, Surname: %s, SeqKey: %v", u.Name, u.Surname, payload.MetaData[seqKey])
	}
}
