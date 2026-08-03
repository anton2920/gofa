//go:build freebsd && 386
// +build freebsd,386

package freebsd

/* From <sys/stat.h>. */
type Stat_t struct {
	Dev          uint64 /* inode's device */
	Ino          uint64 /* inode's number */
	Nlink        uint64 /* number of hard links */
	Mode         uint16 /* inode protection mode */
	_            int16
	Uid          uint32 /* user ID of the file's owner */
	Gid          uint32 /* group ID of the file's group */
	_            int32
	Rdev         uint64 /* device type */
	AtimeExt     int32
	Atime        Timespec /* time of last access */
	MtimeExt     int32
	Mtime        Timespec /* time of last data modification */
	CtimeExt     int32
	Ctime        Timespec /* time of last file status change */
	BirthtimeExt int32
	Birthtime    Timespec /* time of file creation */
	Size         int64    /* file size, in bytes */
	Blocks       int64    /* blocks allocated for file */
	Blksize      int32    /* optimal blocksize for I/O */
	Flags        uint32   /* user defined flags for file */
	Gen          uint64   /* file generation number */
	_            [10]int
}
