package main

import (
	"context"
	"github.com/simon-engledew/seed"
	"github.com/simon-engledew/seed/consumers"
	"github.com/simon-engledew/seed/distribution"
	"github.com/simon-engledew/seed/inspectors/mysql_schema"
	"os"
	"strings"
)

func main() {
	i := mysql_schema.InspectMySQLSchema(strings.NewReader(`
	CREATE TABLE owners (
		id BIGINT UNSIGNED,
		name VARCHAR(255)
	);

    CREATE TABLE cats (
		id BIGINT UNSIGNED,
		owner_id BIGINT,
		name VARCHAR(255)
    );
	`))
	schema, err := seed.Build(i)
	if err != nil {
		panic(err)
	}
	generator := schema.Generator(context.Background(), consumers.MySQLInsertWriter(os.Stdout, 100))
	// generate between 100 and 300 owners
	generator.Insert("owners", distribution.Range(100, 200), func(g *seed.RowGenerator) {
		// generate cats for 3/10 owners
		g.Insert("cats", distribution.Ratio(0.3))
	})
	if err := generator.Wait(); err != nil {
		panic(err)
	}
}
