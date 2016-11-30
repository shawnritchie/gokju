package event

import (
	"testing"
)

type (
	TestingIdentifierVersionEvent struct {
		Event `id:"TestingIdentifierVersionEvent" v:"3"`
		add int
		str string
		long uint64
	}

	NotAnEvent struct {
		add int
		str string
		long uint64
	}
)

func TestEventDefinition_EventID(t *testing.T) {
	id := EventIdentifier((*TestingIdentifierVersionEvent)(nil))
	if id != "TestingIdentifierVersionEvent" {
		t.Errorf("Identifier extracted was '%v' expected 'TestingIdentifierVersionEvent'", id)
	}
}

func TestEventDefinition_NotEventID(t *testing.T) {
	id := EventIdentifier((*NotAnEvent)(nil))
	if id != "event.NotAnEvent" {
		t.Errorf("Identifier extracted was '%v' expected 'TestingIdentifierVersionEvent'", id)
	}
}

func TestEventDefinition_Version(t *testing.T) {
	v := EventVersion((*TestingIdentifierVersionEvent)(nil))
	if v != 3 {
		t.Errorf("Version extracted was v%v expected v3", v)
	}
}

func TestEventDefinition_DefaultVersion(t *testing.T) {
	v := EventVersion((*NotAnEvent)(nil))
	if v != 0 {
		t.Errorf("Version extracted was v%v expected v0", v)
	}
}

