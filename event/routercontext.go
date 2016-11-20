package event

import (
	"github.com/shawnritchie/gokju/structs"
	"reflect"
)

type RouterContext struct {
	metaDataDef       structs.MetaDataDefinition
	methodDefinitions map[reflect.Type]func(v reflect.Value, c EventContainer)
}

func NewRouterContext(definition structs.MetaDataDefinition) RouterContext {
	return RouterContext {
		metaDataDef: definition,
		methodDefinitions: extractMethodDefinitions(definition),
	}
}

func combinations(t structs.MetaDataDefinition) [][]structs.MetaDataIdentifier {
	results := make([][]structs.MetaDataIdentifier, ((t.Len() + 1)*t.Len())/2)

	c := 0
	for i := 0; i < t.Len(); i++ {
		for j := 0; j < t.Len()-i; j++ {
			results[c] = make([]structs.MetaDataIdentifier, i+1)
			for k := 0; k <= i; k++ {
				results[c][k] = t.Get(k + j)
			}
			c++
		}
	}

	return results;
}

func permutations(set []structs.MetaDataIdentifier) [][]structs.MetaDataIdentifier {
	results := make([][]structs.MetaDataIdentifier, 0)

	var permutationsTailRecursive func(n int, input []structs.MetaDataIdentifier)
	permutationsTailRecursive = func(n int, s []structs.MetaDataIdentifier){
		local := make([]structs.MetaDataIdentifier, len(s))
		copy(local, s)

		if (n == 1) {
			results = append(results, local)
		} else {
			for i := 0; i < n; i++ {
				permutationsTailRecursive(n - 1, local)

				if n % 2 == 0 {
					local[0], local[n - 1] = local[n - 1], local[0]
				} else {
					local[i], local[n - 1] = local[n - 1], local[i]
				}
			}
		}
	}

	permutationsTailRecursive(len(set), set)
	return results
}

func allPossibleOutcomes(t structs.MetaDataDefinition) [][]structs.MetaDataIdentifier {
	results := make([][]structs.MetaDataIdentifier, 0)
	for _, s := range combinations(t) {
		for _, a := range permutations(s) {
			results = append(results, a)
		}
	}
	return results
}

func fxDefinition(t *structs.MetaDataDefinition, paramKeys []structs.MetaDataIdentifier) reflect.Type {
	types := make([]reflect.Type, len(paramKeys)+2)
	types[0] = reflect.TypeOf((*interface{})(nil)).Elem()
	types[1] = reflect.TypeOf((*Eventer)(nil)).Elem()
	for i, key := range paramKeys {
		types[i+2] = (*t).Type(key)
	}

	return reflect.FuncOf(types, []reflect.Type{}, false)
}

func extractMethodDefinitions(t structs.MetaDataDefinition) map[reflect.Type]func(v reflect.Value, c EventContainer) {
	allParamKeys := append(allPossibleOutcomes(t), []structs.MetaDataIdentifier{})
	results := make(map[reflect.Type]func(v reflect.Value, c EventContainer), len(allParamKeys))
	for _, paramKeys := range allParamKeys {

		copyKeys:= make([]structs.MetaDataIdentifier, len(paramKeys))
		copy(copyKeys, paramKeys)
		params := make([]reflect.Value, len(paramKeys)+1)

		results[fxDefinition(&t, paramKeys)] = func(v reflect.Value, c EventContainer) {
			keys := copyKeys
			vals := params
			vals[0] = reflect.ValueOf(c.Event)
			for i, key := range keys {
				vals[i+1] = reflect.ValueOf(c.MetaData[key])
			}
			v.Call(vals)
		}
	}
	return results
}
