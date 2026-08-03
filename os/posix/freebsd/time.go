package freebsd

/* From <sys/_timespec.h>. */
type Timespec struct {
	Sec  int
	Nsec int
}

/* From <sys/_timeval.h>. */
type Timeval struct {
	Sec  int
	Usec int
}

const (
	/* From <sys/_clock_id.h>. */
	CLOCK_REALTIME = 0
)
