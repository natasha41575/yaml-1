/*
Copyright 2022 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package yaml

import (
	"encoding/json"
	"fmt"
	"testing"

	"gopkg.in/yaml.v2"
)

func newBenchmarkObject() interface{} {
	data := struct {
		Object map[string]interface{}
		Items  []interface{}
	}{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "PodList",
		},
		Items: []interface{}{},
	}
	for i := 0; i < 1000; i++ {
		item := struct {
			Object map[string]interface{}
		}{
			Object: map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "Pod",
				"metadata": map[string]interface{}{
					"name":      fmt.Sprintf("pod%d", i),
					"namespace": "ns",
					"labels": map[string]interface{}{
						"first-label":  "12",
						"second-label": "label-value",
					},
				},
			},
		}
		data.Items = append(data.Items, item)
	}
	return data
}

func newBenchmarkYAML() ([]byte, error) {
	return yaml.Marshal(newBenchmarkObject())
}

func BenchmarkMarshal(b *testing.B) {
	obj := newBenchmarkObject()
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			result, err := Marshal(obj)
			if err != nil {
				b.Errorf("error marshaling YAML: %v", err)
			}
			b.SetBytes(int64(len(result)))
		}
	})
}

func BenchmarkUnmarshal(b *testing.B) {
	yamlBytes, err := newBenchmarkYAML()
	if err != nil {
		b.Fatalf("error initializing YAML: %v", err)
	}
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var result interface{}
			if err = Unmarshal(yamlBytes, &result); err != nil {
				b.Errorf("error unmarshaling YAML: %v", err)
			}
		}
	})
	b.SetBytes(int64(len(yamlBytes)))
}

func BenchmarkUnmarshalStrict(b *testing.B) {
	yamlBytes, err := newBenchmarkYAML()
	if err != nil {
		b.Fatalf("error initializing YAML: %v", err)
	}
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var result interface{}
			if err = UnmarshalStrict(yamlBytes, &result); err != nil {
				b.Errorf("error unmarshaling YAML (Strict): %v", err)
			}
		}
	})
	b.SetBytes(int64(len(yamlBytes)))
}

func BenchmarkJSONToYAML(b *testing.B) {
	yamlBytes, err := newBenchmarkYAML()
	if err != nil {
		b.Fatalf("error initializing YAML: %v", err)
	}
	jsonBytes, err := YAMLToJSON(yamlBytes)
	if err != nil {
		b.Fatalf("error initializing JSON: %v", err)
	}
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			result, err := JSONToYAML(jsonBytes)
			if err != nil {
				b.Errorf("error converting JSON to YAML: %v", err)
			}
			b.SetBytes(int64(len(result)))
		}
	})
}

func BenchmarkYAMLtoJSON(b *testing.B) {
	yamlBytes, err := newBenchmarkYAML()
	if err != nil {
		b.Fatalf("error initializing YAML: %v", err)
	}
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			result, err := YAMLToJSON(yamlBytes)
			if err != nil {
				b.Errorf("error converting YAML to JSON: %v", err)
			}
			b.SetBytes(int64(len(result)))
		}
	})
}

func BenchmarkYAMLtoJSONStrict(b *testing.B) {
	yamlBytes, err := newBenchmarkYAML()
	if err != nil {
		b.Fatalf("error initializing YAML: %v", err)
	}
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			result, err := YAMLToJSONStrict(yamlBytes)
			if err != nil {
				b.Errorf("error converting YAML to JSON (Strict): %v", err)
			}
			b.SetBytes(int64(len(result)))
		}
	})
}

func BenchmarkJSONObjectToYAMLObject(b *testing.B) {
	yamlBytes, err := newBenchmarkYAML()
	if err != nil {
		b.Fatalf("error initializing YAML: %v", err)
	}
	jsonBytes, err := YAMLToJSON(yamlBytes)
	if err != nil {
		b.Fatalf("error initializing JSON: %v", err)
	}
	var m map[string]interface{}
	err = json.Unmarshal(jsonBytes, &m)
	if err != nil {
		b.Fatalf("error initializing map: %v", err)
	}
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			JSONObjectToYAMLObject(m)
		}
	})
}
