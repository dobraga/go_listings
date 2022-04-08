package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Connect() *sql.DB {
	db, err := sql.Open("postgres", viper.Get("SQLALCHEMY_DATABASE_URI").(string))
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func listTables() []string {
	db := Connect()

	var tables []string
	var table string

	result, err := db.Query("SELECT tablename FROM pg_catalog.pg_tables;")
	Check(err)
	defer result.Close()

	for result.Next() {
		err := result.Scan(&table)
		if err != nil {
			log.Fatal(err)
		}

		tables = append(tables, table)
	}

	return tables
}

func CreateTables() {
	db := Connect()

	all_tables := listTables()

	var need_tables = []string{"metro", "imovel", "imovel_ativo"}

	for i, table := range need_tables {
		if !Contains(all_tables, table) {
			log.Info(fmt.Sprintf("Creating '%s' table", table))
			sql_file := fmt.Sprintf("sql/ddl/%d_%s.sql", i, table)
			query := readQuery(sql_file)
			_, err := db.Exec(query)
			Check(err)
			log.Info(fmt.Sprintf("Created '%s' table", table))
		}
	}

}

func readQuery(file string) string {
	query, err := os.ReadFile(file)
	Check(err)
	str_query := string(query)
	return str_query
}
