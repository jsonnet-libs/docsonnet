package main

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRemarshal(t *testing.T) {
	o := Object{
		Help: "grafana.libsonnet is the offical Jsonnet library for Grafana",
		Fields: map[string]Field{
			"new": {Function: &Function{
				Name: "new",
				Help: "new returns Grafana resources with sane defaults",
			}},
			"addConfig": {Function: &Function{
				Name: "addConfig",
				Help: "addConfig adds config entries to grafana.ini",
			}},
			"datasource": {Object: &Object{
				Name: "datasource",
				Help: "ds-util makes creating datasources easy",
				Fields: map[string]Field{
					"new": {Function: &Function{
						Name: "new",
						Help: "new creates a new datasource",
					}},
				},
			}},
		},
	}

	data, err := json.Marshal(o)
	if err != nil {
		t.Fatal(err)
	}

	var got Object
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}

	if str := cmp.Diff(o, got); str != "" {
		t.Fatal(str)
	}
}
