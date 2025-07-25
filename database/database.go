package database

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/anton2920/gofa/bits"
	"github.com/anton2920/gofa/errors"
	"github.com/anton2920/gofa/syscall"
	"github.com/anton2920/gofa/util"
)

type ID int32

type RecordHeader struct {
	ID    ID
	Flags bits.Flags
}

type DB struct {
	Version uint32
	FD      int32
}

const Version uint32 = 0x0

const (
	VersionOffset = 0
	NextIDOffset  = VersionOffset + int64(unsafe.Sizeof(Version))
	DataOffset    = NextIDOffset + int64(unsafe.Sizeof(ID(0)))
)

const (
	MinValidID = 0
	MaxValidID = (1 << 31) - 1
)

const FlagsCustomOffset = 16

const (
	FlagsNone    = bits.Flags(0)
	FlagsDeleted = bits.Flags(1 << (iota - 1))
)

/* NOTE(anton2920): to bypass check in runtime·adjustpoiners. */
const MinValidPointer = 4096

var NotFound = errors.New("not found")

/* Offset2String performs s.Ptr += base-MinValidPointer. */
//go:nosplit
func Offset2String(s string, base *byte) string {
	return *(*string)(unsafe.Pointer(&reflect.StringHeader{Data: uintptr(util.Noescape(util.PtrAdd(unsafe.Pointer(base), int(uintptr(unsafe.Pointer(util.StringData(s)))-MinValidPointer)))), Len: len(s)}))
}

/* Offset2Slice performs s.Ptr += base-MinValidPointer. */
//go:nosplit
func Offset2Slice(s []byte, base *byte) []byte {
	if len(s) == 0 {
		return s
	}
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{Data: uintptr(util.Noescape(util.PtrAdd(unsafe.Pointer(base), int(uintptr(unsafe.Pointer(&s[0]))-MinValidPointer)))), Len: len(s), Cap: cap(s)}))
}

/* String2Offset performs s.Ptr = offset+MinValidPointer. */
//go:nosplit
func String2Offset(s string, offset int) string {
	return *(*string)(unsafe.Pointer(&reflect.StringHeader{Data: uintptr(offset) + MinValidPointer, Len: len(s)}))
}

/* Slice2Offset performs s.Ptr = offset+MinValidPointer. */
//go:nosplit
func Slice2Offset(s []byte, offset int) []byte {
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{Data: uintptr(offset) + MinValidPointer, Len: len(s), Cap: cap(s)}))
}

//go:nosplit
func String2DBString(ds *string, ss string, data []byte, n int) int {
	nbytes := copy(data[n:], ss)
	*ds = String2Offset(ss, n)
	return nbytes
}

//go:nosplit
func Slice2DBSlice(ds *[]byte, ss []byte, size int, alignment int, data []byte, n int) int {
	if len(ss) == 0 {
		return 0
	}

	start := util.AlignUp(n, alignment)
	nbytes := copy(data[start:], *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{Data: uintptr(unsafe.Pointer(&ss[0])), Len: len(ss) * size, Cap: len(ss) * size})))
	*ds = Slice2Offset(ss, start)
	return nbytes + (start - n)
}

func ID2Slice(x *ID) []byte {
	size := int(unsafe.Sizeof(*x))
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{Data: uintptr(unsafe.Pointer(x)), Len: size, Cap: size}))
}

func uint322Slice(x *uint32) []byte {
	size := int(unsafe.Sizeof(*x))
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{Data: uintptr(unsafe.Pointer(x)), Len: size, Cap: size}))
}

func Open(path string) (*DB, error) {
	var err error

	db := new(DB)

	db.FD, err = syscall.Open(path, syscall.O_RDWR|syscall.O_CREAT, 0644)
	if err != nil {
		return nil, err
	}

	n, err := syscall.Pread(db.FD, uint322Slice(&db.Version), VersionOffset)
	if err != nil {
		syscall.Close(db.FD)
		return nil, err
	}
	if n < int(unsafe.Sizeof(db.Version)) {
		db.Version = Version

		_, err := syscall.Pwrite(db.FD, uint322Slice(&db.Version), VersionOffset)
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

	_, err := syscall.Pread(db.FD, ID2Slice(&id), NextIDOffset)
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
	_, err := syscall.Pwrite(db.FD, ID2Slice(&id), NextIDOffset)
	if err != nil {
		return fmt.Errorf("failed to write next ID: %w", err)
	}
	return nil
}

func GetOffsetForID(id ID, size int) int64 {
	return int64(int(id)*size) + DataOffset
}

func Read(db *DB, id ID, t unsafe.Pointer, size int) error {
	offset := int64(int(id)*size) + DataOffset

	n, err := syscall.Pread(db.FD, *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{Data: uintptr(t), Len: size, Cap: size})), offset)
	if err != nil {
		return fmt.Errorf("failed to read record from DB: %w", err)
	}
	if n < size {
		return NotFound
	}

	return nil
}

func ReadMany(db *DB, pos *int64, ts []byte, size int) (int, error) {
	if *pos < DataOffset {
		*pos = DataOffset
	}

	n, err := syscall.Pread(db.FD, *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{Data: uintptr(unsafe.Pointer(&ts[0])), Len: len(ts) * size, Cap: len(ts) * size})), *pos)
	if err != nil {
		return 0, fmt.Errorf("failed to read records from DB: %w", err)
	}
	*pos += int64(n)

	return n / size, nil
}

func Write(db *DB, id ID, t unsafe.Pointer, size int) error {
	offset := int64(int(id)*size) + DataOffset

	_, err := syscall.Pwrite(db.FD, *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{Data: uintptr(t), Len: size, Cap: size})), offset)
	if err != nil {
		return fmt.Errorf("failed to write record to DB: %w", err)
	}
	return nil
}
