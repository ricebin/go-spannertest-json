package spannertest

import (
	"cloud.google.com/go/spanner"
	"context"
	"fmt"
	"google.golang.org/api/iterator"
	"math/rand"
	"testing"
	"time"
)

func TestIntegration_sum(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	tableName := "SumTestTable"
	client, adminClient, _, cleanup := makeClient(t)
	defer cleanup()
	if err := updateDDL(t, adminClient,
		`CREATE TABLE `+tableName+` (
			Name STRING(50) NOT NULL,
			Num INT64,
		) PRIMARY KEY (Name)`); err != nil {
		t.Fatal(err)
	}

	for _, nv := range []int64{4, 7, 9} {
		rowKey := fmt.Sprintf("rowKey_%v", rand.Int())
		m := spanner.InsertOrUpdate(tableName,
			[]string{"Name", "Num"},
			[]interface{}{rowKey, nv})

		if _, err := client.Apply(context.Background(), []*spanner.Mutation{m}); err != nil {
			t.Fatal(err)
		}
	}

	count, err := readCount(client)
	if err != nil {
		t.Fatal(err)
	}

	if count != 20 {
		t.Fatalf("count expected: 20 but got %v", count)
	}
}

func readCount(client *spanner.Client) (int64, error) {
	stmt := spanner.Statement{
		SQL: `
				SELECT SUM(Num) as Total
				FROM SumTestTable
			`,
	}

	iter := client.Single().Query(context.Background(), stmt)
	defer iter.Stop()

	for {
		row, err := iter.Next()

		if err == iterator.Done {
			return 0, fmt.Errorf("done not expected")
		}

		var out int64
		if err := row.Columns(&out); err != nil {
			return 0, err
		}
		return out, nil
	}
	return 0, fmt.Errorf("rows not found")
}
