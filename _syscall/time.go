package syscall

/* From <time.h>. */
type Timespec struct {
	Sec, Nsec int64
}

type Timeval struct {
	Sec, Usec int64
}

const (
	/* See <sys/_clock_id.h>. */
	CLOCK_REALTIME = 0
)
