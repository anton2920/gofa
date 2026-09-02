package kv

import (
	"errors"

	"github.com/anton2920/gofa/database"
	"github.com/anton2920/gofa/l10n"
)

type Tx struct {
	l10n.Language
}

func (tx *Tx) Commit() error {
	return nil
}

func (tx *Tx) Rollback() error {
	return nil
}

func (tx *Tx) Add(value []byte) (database.ID, error) {
	return 0, errors.New(tx.L("not implemented"))
}

func (tx *Tx) Get(key database.ID, value []byte) error {
	return errors.New(tx.L("not implemented"))
}

func (tx *Tx) Has(key database.ID, values []byte) error {
	return errors.New(tx.L("not implemented"))
}

func (tx *Tx) Set(key database.ID, value []byte) error {
	return errors.New(tx.L("not implemented"))
}
