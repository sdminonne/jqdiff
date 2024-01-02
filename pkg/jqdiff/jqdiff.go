package jqdiff

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type Jqndiffer interface {
	Compare(ref, actual []byte) ([]Diff, error)
}

type jqdiff struct {
}

func NewJqdiff() *jqdiff {
	return &jqdiff{}
}

func (j *jqdiff) Compare(ref, actual []byte) ([]Diff, error) {
	var refData any
	err := json.Unmarshal([]byte(ref), &refData)
	if err != nil {
		return []Diff{}, fmt.Errorf("can't unmarshal reference data: %v", err)
	}

	var actualData any
	err = json.Unmarshal([]byte(actual), &actualData)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal actual data: %v", err)
	}
	return j.compareAnyData("", refData, actualData), nil
}

func (j *jqdiff) compareAnyData(element string, refData, actualData any) []Diff {
	switch r := refData.(type) {
	case nil: // to handle json 'null'
		if reflect.TypeOf(actualData) == nil {
			return []Diff{}
		}
		if len(element) == 0 {
			element = "."
		}
		return []Diff{typedDiff[any, any]{
			Selector:  element,
			Reference: refData,
			Actual:    actualData,
			Kind:      DifferentType}} // typeError since actualData is not a bool
	case bool:
		if a, ok := actualData.(bool); ok {
			return j.compareDataBool(element, r, a)
		}
		return []Diff{typedDiff[bool, any]{
			Selector:  element,
			Reference: refData.(bool),
			Actual:    actualData,
			Kind:      DifferentType, // typeError since actualData is not a bool
		}}
	case float64:
		if a, ok := actualData.(float64); ok {
			return j.compareDataFloat64(element, r, a)
		}
		return []Diff{typedDiff[float64, any]{
			Selector:  element,
			Reference: refData.(float64),
			Actual:    actualData,
			Kind:      DifferentType, // typeError since actualData is not a float64
		}}
	case string:
		if a, ok := actualData.(string); ok {
			return j.compareDataStrings(element, r, a)
		}
		return []Diff{
			typedDiff[string, any]{
				Selector:  element,
				Reference: refData.(string),
				Actual:    actualData,
				Kind:      DifferentType, // typeError since actualData is not a string
			}}
	case []interface{}:
		if a, ok := actualData.([]interface{}); ok {
			return j.compareDataArrays(element, r, a)
		}
		return []Diff{
			typedDiff[[]interface{}, any]{
				Selector:  element,
				Reference: refData.([]interface{}),
				Actual:    actualData,
				Kind:      DifferentType, // typeError since actualData is not a []interface{}
			}}
	case map[string]interface{}:
		if a, ok := actualData.(map[string]interface{}); ok {
			return j.compareDataMaps(element, r, a)
		}
		return []Diff{typedDiff[map[string]interface{}, any]{
			Selector:  element,
			Reference: refData.(map[string]interface{}),
			Actual:    actualData,
			Kind:      DifferentType, // typeError since actualData is not a map[string]interface{}
		}}
	default:
		fmt.Printf("-> Hey dude I cannot handle type: %#v\n", reflect.TypeOf(refData))
		panic("unexpected type")
	}

}

func (j *jqdiff) compareDataBool(element string, ref, actual bool) []Diff {
	if ref != actual {
		return []Diff{
			typedDiff[bool, bool]{
				Selector:  element,
				Actual:    actual,
				Reference: ref,
				Kind:      DifferentValue},
		}
	}
	return []Diff{}
}

func (j *jqdiff) compareDataFloat64(element string, ref, actual float64) []Diff {
	if ref != actual {
		return []Diff{
			typedDiff[float64, float64]{
				Selector:  element,
				Actual:    actual,
				Reference: ref,
				Kind:      DifferentValue},
		}
	}
	return []Diff{}
}

func (j *jqdiff) compareDataStrings(element string, ref, actual string) []Diff {
	if len(element) == 0 {
		element = "."
	}
	if ref != actual {
		return []Diff{
			typedDiff[string, string]{
				Selector:  element,
				Actual:    actual,
				Reference: ref,
				Kind:      DifferentValue},
		}
	}
	return []Diff{}
}

func (j *jqdiff) compareDataArrays(element string, ref, actual []interface{}) []Diff {
	// needs to run both side
	diffs := []Diff{}
	refLength := len(ref)
	actualLength := len(actual)
	if len(element) == 0 {
		element = "."
	}
	for i := 0; i < refLength; i++ {
		if i < actualLength {
			currentElement := fmt.Sprintf("%s[%d]", element, i)
			diffs = append(diffs, j.compareAnyData(currentElement, ref[i], actual[i])...)
		}

	}
	return diffs
}

func (j *jqdiff) compareDataMaps(element string, ref, actual map[string]interface{}) []Diff {
	// need to run both side
	diffs := []Diff{}
	for refKey, refValue := range ref {
		currentSelector := element + "." + refKey
		actualValue, ok := actual[refKey]
		if ok {
			diffs = append(diffs, j.compareAnyData(currentSelector, refValue, actualValue)...)
			continue
		}
		diffs = append(diffs, typedDiff[any, any]{
			Selector:  currentSelector,
			Reference: refValue,
			Actual:    nil,
			Kind:      DifferentValue, // missing key/value in actual
		})
	}
	// TODO
	//	for actualKey, actualValue := range actual {
	//		refValue, ok := ref[actualKey]
	//	}
	return diffs
}
