package event

type (
	UpcasterIdentity struct {
		Identifier Identifier
		Version Version
	}

	Upcaster struct {
		UpcasterIdentity
		Intercept func(in *Container)
	}

	UpcasterRegistry map[UpcasterIdentity]Upcaster
	Upcasters []Upcaster
)


func (s Upcasters) Len() int {
	return len(s)
}

func (s Upcasters) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Upcasters) Less(i, j int) bool {
	return s[i].Version < s[j].Version
}
