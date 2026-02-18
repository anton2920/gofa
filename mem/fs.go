package mem

import (
	"sync"
	"sync/atomic"

	"github.com/anton2920/gofa/errors"
	"github.com/anton2920/gofa/io/fs"
	"github.com/anton2920/gofa/syscall"
)

type FS struct{}

type File struct {
	sync.RWMutex

	Data []byte
	Pos  int64

	Flags int
}

var _ fs.VFile = new(File)

func NewFS() fs.VFS {
	return FS{}
}

func (fs FS) Open(path string, flags int32, perms ...uint16) (fs.VFile, error) {
	return &File{Flags: int(flags)}, nil
}

func (fs FS) OpenAt(f fs.VFile, path string, flags int32, perms ...uint16) (fs.VFile, error) {
	return fs.Open(path, flags, perms...)
}

func (fs FS) CreateDirectory(path string, perms uint16) error {
	return errors.NotImplemented
}

func (f *File) Read(buf []byte) (int, error) {
	f.RLock()
	pos := atomic.AddInt64(&f.Pos, int64(len(buf)))
	n := copy(buf, f.Data[pos-int64(len(buf)):])
	f.RUnlock()
	return n, nil
}

func (f *File) Write(buf []byte) (int, error) {
	if (f.Flags & syscall.O_APPEND) == syscall.O_APPEND {
		f.Lock()
		atomic.StoreInt64(&f.Pos, int64(len(f.Data)))
		f.Data = append(f.Data, buf...)
		f.Unlock()
		return len(buf), nil
	} else {
		return -1, errors.NotImplemented
	}

}

func (f *File) Close() error {
	return nil
}

func (f *File) ReadAt(buf []byte, pos int64) (int, error) {
	f.RLock()
	defer f.RUnlock()

	if ((pos == 0) && (len(f.Data) == 0)) || (pos > int64(len(f.Data))) {
		return -1, errors.New("reading past the EOF")
	}

	return copy(buf, f.Data[pos:]), nil
}

func (f *File) WriteAt(buf []byte, pos int64) (int, error) {
	return f.WriteAtEx(buf, pos, false)
}

func (f *File) WriteAtEx(buf []byte, pos int64, lockHeld bool) (int, error) {
	var end []byte

	if !lockHeld {
		f.Lock()
		defer f.Unlock()
	}

	if pos+int64(len(buf)) < int64(len(f.Data)) {
		end = f.Data[pos+int64(len(buf)):]
	}
	f.Data = append(f.Data[pos:], buf...)
	f.Data = append(f.Data[pos+int64(len(buf)):], end...)

	return len(buf), nil
}

func (f *File) Size() (int, error) {
	return f.SizeEx(false)
}

func (f *File) SizeEx(lockHeld bool) (int, error) {
	if !lockHeld {
		f.Lock()
		defer f.Unlock()
	}
	return len(f.Data), nil
}

func (f *File) Sync() error {
	return nil
}

func (f *File) Truncate(n int64) error {
	size, _ := f.Size()
	if n == int64(size) {
		return nil
	} else if n > int64(size) {
		return errors.NotImplemented
	}

	f.Lock()
	f.Data = f.Data[:n]
	atomic.StoreInt64(&f.Pos, int64(len(f.Data)))
	f.Unlock()

	return nil
}

func (f *File) VFS() fs.VFS {
	return NewFS()
}
