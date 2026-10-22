package main

import (
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/AD-Tech-Projects/Ledger/packages/infra/database/models"
)

func main() {
	// ANCHOR - Load all the database models
	stmts, err := gormschema.New("postgres").Load(
		models.AccountHolderModel{},
		models.AccountModel{},
		models.TransactionModel{},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}

	io.WriteString(os.Stdout, stmts)
}
