package fs

type VFS interface {
	Open(path string, flags int32, perms ...uint16) (VFile, error)
	OpenAt(f VFile, path string, flags int32, perms ...uint16) (VFile, error)
	CreateDirectory(path string, perms uint16) error
}

type VFile interface {
	Read(buf []byte) (int, error)
	Write(buf []byte) (int, error)
	Close() error

	ReadAt(buf []byte, pos int64) (int, error)
	WriteAt(buf []byte, pos int64) (int, error)
	WriteAtEx(buf []byte, pos int64, lockHeld bool) (int, error)

	Size() (int, error)
	SizeEx(lockHeld bool) (int, error)

	Sync() error
	Truncate(n int64) error

	Lock()
	Unlock()

	VFS() VFS
}
