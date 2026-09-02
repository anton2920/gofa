package database

import "github.com/anton2920/gofa/bits"

type ID int64

type RecordHeader struct {
	ID    ID
	Flags bits.Flags64

	CreatedAt int64
	UpdatedAt int64
	DeletedAt int64
}
