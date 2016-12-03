package event

import (
	"testing"
	"encoding/json"
	"fmt"
)

type (
	Event1 struct {
		Event `json:"-" id:"JsonEvent" v:"0"`
		Name string `json:"Name"`
		Surname string `json:"Surname"`
	}

	Event3 struct {
		Event `json:"-" id:"JsonEvent" v:"2"`
		FullName string `json:"FullName"`
	}
)

func TestUpcaster_Chaining(t *testing.T) {
	e := Event1{Name:"Shawn", Surname:"Ritchie"}
	s := uint64(1)
	c := NewContainer(e, MetaData{ seqKey: s })
	b, err := json.Marshal(c)
	if (err != nil) {
		t.Errorf(err.Error())
	}
	fmt.Println(string(b))

	u := &struct {
		Event `json:"-" id:"JsonEvent" v:"0"`
		Name string `json:"Name"`
		Surname string `json:"Surname"`
	}{}

	var payload *Container = new(Container)
	payload.UnmarshalJSON(containeKeyDef, u, b)

	interceptorV1 := Upcaster{
		UpcasterIdentity: UpcasterIdentity{
			Identifier:Identifier("JsonEvent"),
			Version:Version(0),
		},
		Intercept:func(in *Container) {
			u := in.Event.(*struct {
				Event `json:"-" id:"JsonEvent" v:"0"`
				Name string `json:"Name"`
				Surname string `json:"Surname"`
			})


			in.Event = &struct {
				Event `json:"-" id:"JsonEvent" v:"1"`
				FirstName string `json:"FirstName"`
				LastName string `json:"LastName"`
			} {
				FirstName:u.Name,
				LastName:u.Surname,
			}
		},
	}

	interceptorV2 := Upcaster{
		UpcasterIdentity: UpcasterIdentity{
			Identifier:Identifier("JsonEvent"),
			Version:Version(1),
		},
		Intercept:func(in *Container) {
			u := in.Event.(*struct {
				Event `json:"-" id:"JsonEvent" v:"1"`
				FirstName string `json:"FirstName"`
				LastName string `json:"LastName"`
			})


			in.Event = &Event3{
				FullName:fmt.Sprintf("%s %s", u.FirstName, u.LastName),
			}
		},
	}

	interceptorV1.Intercept(payload)
	interceptorV2.Intercept(payload)

	if (payload.Event.(*Event3).FullName != "Shawn Ritchie") {
		t.Errorf("Expected Value 'Shawn Ritchie', but produced '%v'", payload.Event.(*Event3).FullName)
	}
}