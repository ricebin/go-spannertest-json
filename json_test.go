package spannertest

import (
	"cloud.google.com/go/spanner"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func TestIntegration_Json(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	tableName := "JsonTestTable"
	client, adminClient, _, cleanup := makeClient(t)
	defer cleanup()
	if err := updateDDL(t, adminClient,
		`CREATE TABLE `+tableName+` (
			Name STRING(50) NOT NULL,
			jv JSON,
		) PRIMARY KEY (Name)`); err != nil {
		t.Fatal(err)
	}

	inJson := Metadata{
		Creators: []Creators{
			{
				Address: fmt.Sprintf("creator_%v", rand.Int()),
			},
		},
	}

	rowKey := fmt.Sprintf("rowKey_%v", rand.Int())
	m := spanner.InsertOrUpdate(tableName,
		[]string{"Name", "jv"},
		[]interface{}{rowKey, spanner.NullJSON{Value: inJson, Valid: true}})

	if _, err := client.Apply(context.Background(), []*spanner.Mutation{m}); err != nil {
		t.Fatal(err)
	}

	row, err := client.Single().ReadRow(context.Background(), tableName, spanner.Key{rowKey}, []string{"jv"})
	if err != nil {
		t.Fatal(err)
	}

	{
		var outStr string
		if err := row.Columns(&outStr); err == nil {
			t.Fatal("should not be able to decode json to string")
		}
	}

	{
		out := &spanner.NullJSON{}
		if err := row.Columns(out); err != nil {
			t.Fatal(err)
		}

		outJson := &Metadata{}
		if err := json.Unmarshal([]byte(out.String()), outJson); err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(toJson(inJson), toJson(outJson)) {
			t.Fatalf("wanted %v, but got %v", inJson, outJson)
		}
	}
}

func toJson(value interface{}) string {
	out, err := json.Marshal(value)
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}

type Metadata struct {
	Creators []Creators `json:"creators"`
}

type Creators struct {
	Address string `json:"address"`
}
