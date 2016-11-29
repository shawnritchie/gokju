package event

import (
	"reflect"
)

type Container struct {
	Event    Eventer
	MetaData MetaData
}

type (
	ContainerContext interface {
		MapFunctionType(fxType reflect.Type) func(v reflect.Value, c Container)
		AllFunctionType() []reflect.Type
	}

	DefaultContainerContext struct {
		metaDataDef       MetaDataDefinition
		methodDefinitions map[reflect.Type]func(v reflect.Value, c Container)
	}
)

func (c *DefaultContainerContext)MapFunctionType(fxType reflect.Type) func(v reflect.Value, c Container) {
	return c.methodDefinitions[fxType]
}

func (c *DefaultContainerContext)AllFunctionType() []reflect.Type {
	keys := make([]reflect.Type, len(c.methodDefinitions))
	i := 0
	for  k, _ := range c.methodDefinitions {
		keys[i] = k
		i++
	}
	return keys
}

func NewContainerContext(definition MetaDataDefinition) DefaultContainerContext {
	return DefaultContainerContext{
		metaDataDef: definition,
		methodDefinitions: extractMethodDefinitions(definition),
	}
}

func combinations(t MetaDataDefinition) [][]MetaDataIdentifier {
	results := make([][]MetaDataIdentifier, ((t.Len() + 1)*t.Len())/2)

	c := 0
	for i := 0; i < t.Len(); i++ {
		for j := 0; j < t.Len()-i; j++ {
			results[c] = make([]MetaDataIdentifier, i+1)
			for k := 0; k <= i; k++ {
				results[c][k] = t.Get(k + j)
			}
			c++
		}
	}

	return results;
}

func permutations(set []MetaDataIdentifier) [][]MetaDataIdentifier {
	results := make([][]MetaDataIdentifier, 0)

	var permutationsTailRecursive func(n int, input []MetaDataIdentifier)
	permutationsTailRecursive = func(n int, s []MetaDataIdentifier){
		local := make([]MetaDataIdentifier, len(s))
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

func allPossibleOutcomes(t MetaDataDefinition) [][]MetaDataIdentifier {
	results := make([][]MetaDataIdentifier, 0)
	for _, s := range combinations(t) {
		for _, a := range permutations(s) {
			results = append(results, a)
		}
	}
	return results
}

func fxDefinition(t *MetaDataDefinition, paramKeys []MetaDataIdentifier) reflect.Type {
	types := make([]reflect.Type, len(paramKeys)+2)
	types[0] = reflect.TypeOf((*interface{})(nil)).Elem()
	types[1] = reflect.TypeOf((*Eventer)(nil)).Elem()
	for i, key := range paramKeys {
		types[i+2] = (*t).Type(key)
	}

	return reflect.FuncOf(types, []reflect.Type{}, false)
}

func extractMethodDefinitions(t MetaDataDefinition) map[reflect.Type]func(v reflect.Value, c Container) {
	allParamKeys := append(allPossibleOutcomes(t), []MetaDataIdentifier{})
	results := make(map[reflect.Type]func(v reflect.Value, c Container), len(allParamKeys))
	for _, paramKeys := range allParamKeys {

		copyKeys:= make([]MetaDataIdentifier, len(paramKeys))
		copy(copyKeys, paramKeys)
		params := make([]reflect.Value, len(paramKeys)+1)

		results[fxDefinition(&t, paramKeys)] = func(v reflect.Value, c Container) {
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
