package syscall

/* From <sys/stat.h>. */
type Stat_t struct {
	Dev       uint   /* inode's device */
	Ino       uint   /* inode's number */
	Nlink     uint64 /* number of hard links */
	Mode      uint16 /* inode protection mode */
	_         int16
	Uid       uint32 /* user ID of the file's owner */
	Gid       uint32 /* group ID of the file's group */
	_         int32
	Rdev      uint64   /* device type */
	Atime     Timespec /* time of last access */
	Mtime     Timespec /* time of last data modification */
	Ctime     Timespec /* time of last file status change */
	Birthtime Timespec /* time of file creation */
	Size      int      /* file size, in bytes */
	Blocks    int      /* blocks allocated for file */
	Blksize   int32    /* optimal blocksize for I/O */
	Flags     uint32   /* user defined flags for file */
	Gen       uint64   /* file generation number */
	_         [10]int
}
