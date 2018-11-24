package trancatemodel

import (
	"flag"
	"log"
	"os"
	"testing"
)

// TestMain model package setup/teardonw
func TestMain(m *testing.M) {
	flag.Parse()

	dbSetupCfg := DBConfig{
		DBName:  "gotodo",
		Host:    "localhost",
		Port:    "5432",
		SSLMode: "disable",
		User:    "gotodo_root",
	}
	tblSetupCfg := DBConfig{
		DBName:  "gotodo",
		Host:    "localhost",
		Port:    "5432",
		SSLMode: "disable",
		User:    "gotodo_api_test",
	}
	testSchema := "gotodo_api_test"
	testUser := "gotodo_api_test"

	TestDropSchema(dbSetupCfg, testSchema)

	if err := TestCreateSchema(dbSetupCfg, testSchema, testUser); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	if err := TestCreateTables(tblSetupCfg, "../schema"); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	code := m.Run()

	if err := TestDropSchema(dbSetupCfg, testSchema); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	os.Exit(code)
}
