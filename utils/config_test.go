package utils

import (
	"encoding/json"
	"testing"
)

// global configuration struct
type ConfigStruct struct {
	AllowOrigins []string `json:"allow_origins"`
	JWTTimeout   int64    `json:"jwt_timeout"` // hours
	Address      string   `json:"address"`
	EnableTLS    bool     `json:"enable_tls"`
	JWTKey       string   `json:"jwt_key"`
}

// default configuration
func (c ConfigStruct) Default() any {
	return ConfigStruct{
		AllowOrigins: []string{"*"},
		JWTTimeout:   24 * 14, // 14day
		Address:      "0.0.0.0:8000",
		EnableTLS:    true,
		JWTKey:       RandomString(64),
	}
}

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
