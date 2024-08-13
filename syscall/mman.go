package syscall

const (
	/* From <sys/mman.h>. */

	/*
	 * Protections are chosen from these bits, or-ed together
	 */
	PROT_NONE  = 0x00 /* no permissions */
	PROT_READ  = 0x01 /* pages can be read */
	PROT_WRITE = 0x02 /* pages can be written */
	PROT_EXEC  = 0x04 /* pages can be executed */

	/*
	 * Flags contain sharing type and options.
	 * Sharing types; choose one.
	 */
	MAP_SHARED  = 0x0001 /* share changes */
	MAP_PRIVATE = 0x0002 /* changes are private */

	/*
	 * Other flags
	 */
	MAP_FIXED        = 0x0010 /* map addr must be exactly as requested */
	MAP_HASSEMAPHORE = 0x0200 /* region may contain semaphores */
	MAP_STACK        = 0x0400 /* region grows down, like a stack */
	MAP_NOSYNC       = 0x0800 /* page to but do not sync underlying file */

	/*
	 * Mapping type
	 */
	MAP_FILE = 0x0000 /* map from file (default) */
	MAP_ANON = 0x1000 /* allocated from memory, swap space */

	/*
	 * Extended flags
	 */
	MAP_GUARD         = 0x00002000 /* reserve but don't map address range */
	MAP_EXCL          = 0x00004000 /* for MAP_FIXED, fail if address is used */
	MAP_NOCORE        = 0x00020000 /* dont include these pages in a coredump */
	MAP_PREFAULT_READ = 0x00040000 /* prefault mapping for reading */

	MAP_32BIT = 0x00080000 /* map in the low 2GB of address space */

	MAP_ALIGNMENT_SHIFT = 24
	MAP_ALIGNED_SUPER   = 1 << MAP_ALIGNMENT_SHIFT /* align on a superpage */

	/*
	 * shmflags for shm_open2()
	 */
	SHM_ALLOW_SEALING = 0x00000001
	SHM_GROW_ON_WRITE = 0x00000002
	SHM_LARGEPAGE     = 0x00000004

	SHM_LARGEPAGE_ALLOC_DEFAULT = 0
	SHM_LARGEPAGE_ALLOC_NOWAIT  = 1
	SHM_LARGEPAGE_ALLOC_HARD    = 2
)
