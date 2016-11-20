package event

import (
	"testing"
	"fmt"
	"github.com/shawnritchie/gokju/structs"
	"reflect"
	"time"
)

type benchmarkKey int

const (
	seqKey benchmarkKey = iota
	timestampKey
	stringKey
)

func (k benchmarkKey)ToInt() int {
	return int(k)
}

var containeKeyDef structs.MetaDataDefinition = structs.MetaDataDefinition {
	Keys: structs.MetaDataMap{seqKey: reflect.TypeOf((*uint64)(nil)).Elem(),
		timestampKey: reflect.TypeOf((*time.Time)(nil)).Elem(),
		stringKey: reflect.TypeOf((*string)(nil)).Elem(),
	},
	Generator:func(i int) structs.MetaDataIdentifier{
		return benchmarkKey(i)
	},
}

func match(match, with [][]structs.MetaDataIdentifier) bool {
	if (len(match) != len(with)) {
		return false;
	}

	for i, v := range match {
		if (len(v) != len(with[i])) {
			return false
		}
		for j, s := range v {
			if (s != with[i][j]) {
				return false
			}
		}
	}
	return true;
}

func TestGetAllKeyCombinations(t *testing.T) {
	expectedResults := [][]structs.MetaDataIdentifier{
		[]structs.MetaDataIdentifier{benchmarkKey(0)}, []structs.MetaDataIdentifier{benchmarkKey(1)}, []structs.MetaDataIdentifier{benchmarkKey(2)},
		[]structs.MetaDataIdentifier{benchmarkKey(0), benchmarkKey(1)}, []structs.MetaDataIdentifier{benchmarkKey(1), benchmarkKey(2)},
		[]structs.MetaDataIdentifier{benchmarkKey(0), benchmarkKey(1), benchmarkKey(2)},
	}

	testResults := combinations(containeKeyDef)

	if !match(expectedResults, testResults) {
		t.Errorf("The results generated by getAllKeySets do not match\n")
		t.Errorf("Expected: %v\n", expectedResults)
		t.Errorf("Generated: %v\n", testResults)
	} else {
		fmt.Printf("Expected: %v\n", expectedResults)
		fmt.Printf("Generated: %v\n", testResults)
	}

}

func TestPermutations(t *testing.T) {
	expectedResults := [][]structs.MetaDataIdentifier{
		[]structs.MetaDataIdentifier{benchmarkKey(0), benchmarkKey(1), benchmarkKey(2)},
		[]structs.MetaDataIdentifier{benchmarkKey(1), benchmarkKey(0), benchmarkKey(2)},
		[]structs.MetaDataIdentifier{benchmarkKey(2), benchmarkKey(1), benchmarkKey(0)},
		[]structs.MetaDataIdentifier{benchmarkKey(1), benchmarkKey(2), benchmarkKey(0)},
		[]structs.MetaDataIdentifier{benchmarkKey(2), benchmarkKey(0), benchmarkKey(1)},
		[]structs.MetaDataIdentifier{benchmarkKey(0), benchmarkKey(2), benchmarkKey(1)},
	}

	testResults := permutations([]structs.MetaDataIdentifier{benchmarkKey(0), benchmarkKey(1), benchmarkKey(2)})

	if !match(expectedResults, testResults) {
		t.Errorf("The results generated by getAllKeySets do not match\n")
		t.Errorf("Expected: %v\n", expectedResults)
		t.Errorf("Generated: %v\n", testResults)
	} else {
		fmt.Printf("Expected: %v\n", expectedResults)
		fmt.Printf("Generated: %v\n", testResults)
	}
}

func TestAllPossibleOutcomes(t *testing.T) {
	expectedResults := [][]structs.MetaDataIdentifier{
		[]structs.MetaDataIdentifier{benchmarkKey(0)}, []structs.MetaDataIdentifier{benchmarkKey(1)}, []structs.MetaDataIdentifier{benchmarkKey(2)},
		[]structs.MetaDataIdentifier{benchmarkKey(0), benchmarkKey(1)}, []structs.MetaDataIdentifier{benchmarkKey(1), benchmarkKey(0)},
		[]structs.MetaDataIdentifier{benchmarkKey(1), benchmarkKey(2)}, []structs.MetaDataIdentifier{benchmarkKey(2), benchmarkKey(1)},
		[]structs.MetaDataIdentifier{benchmarkKey(0), benchmarkKey(1), benchmarkKey(2)},
		[]structs.MetaDataIdentifier{benchmarkKey(1), benchmarkKey(0), benchmarkKey(2)},
		[]structs.MetaDataIdentifier{benchmarkKey(2), benchmarkKey(1), benchmarkKey(0)},
		[]structs.MetaDataIdentifier{benchmarkKey(1), benchmarkKey(2), benchmarkKey(0)},
		[]structs.MetaDataIdentifier{benchmarkKey(2), benchmarkKey(0), benchmarkKey(1)},
		[]structs.MetaDataIdentifier{benchmarkKey(0), benchmarkKey(2), benchmarkKey(1)},
	}

	testResults := allPossibleOutcomes(containeKeyDef)
	if !match(expectedResults, testResults) {
		t.Errorf("The results generated by getAllKeySets do not match\n")
		t.Errorf("Expected: %v\n", expectedResults)
		t.Errorf("Generated: %v\n", testResults)
	} else {
		fmt.Printf("Expected: %v\n", expectedResults)
		fmt.Printf("Generated: %v\n", testResults)
	}
}

func TestExtractMethodDefinitions(t *testing.T) {
	methodDefinitions := extractMethodDefinitions(containeKeyDef)
	if (len(methodDefinitions) != 14) {
		t.Errorf("The amount of method definitions generated is incorrect Expected[14] Generated[%v]", len(methodDefinitions))
	}
}