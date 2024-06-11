package database

import (
	"fmt"
	"unsafe"

	"github.com/anton2920/gofa/errors"
	"github.com/anton2920/gofa/syscall"
	"github.com/anton2920/gofa/util"
)

type ID int32

type DB struct {
	Version uint32
	FD      int32
}

const Version uint32 = 0x0

const (
	VersionOffset = 0
	NextIDOffset  = VersionOffset + int64(unsafe.Sizeof(Version))
	DataOffset    = NextIDOffset + int64(unsafe.Sizeof(id))
)

const (
	MinValidID = 0
	MaxValidID = (1 << 31) - 1
)

/* NOTE(anton2920): to bypass check in runtime·adjustpoiners. */
const MinValidPointer = 4096

/* NOTE(anton2920): for sizeof. */
var id ID

var NotFound = errors.New("not found")

//go:nosplit
func Offset2String(s string, base *byte) string {
	return unsafe.String((*byte)(unsafe.Add(unsafe.Pointer(base), uintptr(unsafe.Pointer(unsafe.StringData(s)))-MinValidPointer)), len(s))
}

//go:nosplit
func Offset2Slice[T any](s []T, base *byte) []T {
	if len(s) == 0 {
		return s
	}
	return unsafe.Slice((*T)(unsafe.Add(unsafe.Pointer(base), uintptr(unsafe.Pointer(unsafe.SliceData(s)))-MinValidPointer)), len(s))
}

//go:nosplit
func String2Offset(s string, offset int) string {
	return unsafe.String((*byte)(unsafe.Pointer(uintptr(offset)+MinValidPointer)), len(s))
}

//go:nosplit
func Slice2Offset[T any](s []T, offset int) []T {
	return unsafe.Slice((*T)(unsafe.Pointer(uintptr(offset)+MinValidPointer)), len(s))
}

//go:nosplit
func String2DBString(ds *string, ss string, data []byte, n int) int {
	nbytes := copy(data[n:], ss)
	*ds = String2Offset(ss, n)
	return nbytes
}

//go:nosplit
func Slice2DBSlice[T any](ds *[]T, ss []T, data []byte, n int) int {
	if len(ss) == 0 {
		return 0
	}

	start := util.RoundUp(n, int(unsafe.Alignof(&ss[0])))
	nbytes := copy(data[start:], unsafe.Slice((*byte)(unsafe.Pointer(&ss[0])), len(ss)*int(unsafe.Sizeof(ss[0]))))
	*ds = Slice2Offset(ss, start)
	return nbytes + (start - n)
}

func Open(path string) (*DB, error) {
	var err error

	db := new(DB)

	db.FD, err = syscall.Open(path, syscall.O_RDWR|syscall.O_CREAT, 0644)
	if err != nil {
		return nil, err
	}

	n, err := syscall.Pread(db.FD, unsafe.Slice((*byte)(unsafe.Pointer(&db.Version)), unsafe.Sizeof(db.Version)), VersionOffset)
	if err != nil {
		syscall.Close(db.FD)
		return nil, err
	}
	if n < int(unsafe.Sizeof(db.Version)) {
		db.Version = Version

		_, err := syscall.Pwrite(db.FD, unsafe.Slice((*byte)(unsafe.Pointer(&db.Version)), unsafe.Sizeof(db.Version)), VersionOffset)
		if err != nil {
			syscall.Close(db.FD)
			return nil, err
		}
	} else if db.Version != Version {
		syscall.Close(db.FD)
		return nil, fmt.Errorf("incompatible DB file version %d, expected %d", db.Version, Version)
	}

	return db, nil
}

func Close(db *DB) error {
	return syscall.Close(db.FD)
}

func Drop(db *DB) error {
	if err := syscall.Ftruncate(db.FD, DataOffset); err != nil {
		return fmt.Errorf("failed to truncate file: %w", err)
	}
	if err := SetNextID(db, 0); err != nil {
		return fmt.Errorf("failed to set next ID: %w", err)
	}
	return nil
}

func GetNextID(db *DB) (ID, error) {
	var id ID

	_, err := syscall.Pread(db.FD, unsafe.Slice((*byte)(unsafe.Pointer(&id)), unsafe.Sizeof(id)), NextIDOffset)
	if err != nil {
		return -1, fmt.Errorf("failed to read next ID: %w", err)
	}

	return id, nil
}

/* TODO(anton2920): make that atomic. */
func IncrementNextID(db *DB) (ID, error) {
	id, err := GetNextID(db)
	if err != nil {
		return -1, err
	}
	if err := SetNextID(db, id+1); err != nil {
		return -1, err
	}
	return id, nil
}

func SetNextID(db *DB, id ID) error {
	_, err := syscall.Pwrite(db.FD, unsafe.Slice((*byte)(unsafe.Pointer(&id)), unsafe.Sizeof(id)), NextIDOffset)
	if err != nil {
		return fmt.Errorf("failed to write next ID: %w", err)
	}
	return nil
}

func Read[T any](db *DB, id ID, t *T) error {
	size := int(unsafe.Sizeof(*t))
	offset := int64(int(id)*size) + DataOffset

	n, err := syscall.Pread(db.FD, unsafe.Slice((*byte)(unsafe.Pointer(t)), size), offset)
	if err != nil {
		return fmt.Errorf("failed to read record from DB: %w", err)
	}
	if n < size {
		return NotFound
	}

	return nil
}

func ReadMany[T any](db *DB, pos *int64, ts []T) (int, error) {
	if *pos < DataOffset {
		*pos = DataOffset
	}
	size := int(unsafe.Sizeof(ts[0]))

	n, err := syscall.Pread(db.FD, unsafe.Slice((*byte)(unsafe.Pointer(unsafe.SliceData(ts))), len(ts)*size), *pos)
	if err != nil {
		return 0, fmt.Errorf("failed to read records from DB: %w", err)
	}
	*pos += int64(n)

	return n / size, nil
}

func Write[T any](db *DB, id ID, t *T) error {
	size := int(unsafe.Sizeof(*t))
	offset := int64(int(id)*size) + DataOffset

	_, err := syscall.Pwrite(db.FD, unsafe.Slice((*byte)(unsafe.Pointer(t)), size), offset)
	if err != nil {
		return fmt.Errorf("failed to write record to DB: %w", err)
	}
	return nil
}
