package kv

import (
	"github.com/anton2920/gofa/container/bplus"
	"github.com/anton2920/gofa/database"
	"github.com/anton2920/gofa/l10n"
)

type Database struct {
	bplus.Tree
}

func (db *Database) Add(value interface{}) (database.ID, error) {
	return 0, nil
}

func (db *Database) Begin(l l10n.Language) (*Tx, error) {
	return &Tx{Language: l}, nil
}
