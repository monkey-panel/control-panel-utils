package utils

import (
	"encoding/json"
	"testing"
)

func TestConfig(t *testing.T) {
	var config ConfigStruct
	if err := ConfigRead("./_test.json", &config); err != nil {
		t.Fatal(err)
	}
}

func TestMergeObj(t *testing.T) {
	mergeValue := mergeObj(map[string]any{
		"a": "a",
		"b": "b",
		"c": []string{"test"},
		"d": "a",
		"e": map[string]any{"a": "a"},
	}, map[string]any{
		"a": "b",
		"b": "c",
		"c": "c",
		"d": []string{"test2"},
		"e": map[string]any{
			"a": "a",
			"b": "b",
			"c": []string{"test"},
			"d": "a",
		},
	})

	// {
	//     "a": "b",
	//     "b": "c",
	//     "c": "c",
	//     "d": ["test2"],
	//     "e": {
	//         "a": "a",
	//         "b": "b",
	//         "c": ["test"],
	//         "d": "a"
	//     }
	// }
	data, _ := json.Marshal(mergeValue)
	if string(data) != `{"a":"a","b":"b","c":"c","d":["test2"],"e":{"a":"a","b":"b","c":["test"],"d":"a"}}` {
		t.Fatal("mergeObj failed")
	}
}
