package time

import (
	"github.com/anton2920/gofa/slices"
	"github.com/anton2920/gofa/syscall"
)

/* From <time.h>. */
type Tm struct {
	Sec   int /* seconds after the minute [0-60] */
	Min   int /* minutes after the hour [0-59] */
	Hour  int /* hours since midnight [0-23] */
	Mday  int /* day of the month [1-31] */
	Mon   int /* months since January [0-11] */
	Year  int /* years since 1900 */
	Wday  int /* days since Sunday [0-6] */
	Yday  int /* days since January 1 [0-365] */
	Isdst int /* Daylight Savings Time flag */
}

const RFC822Len = 31

func Unix() int {
	var tp syscall.Timespec
	syscall.ClockGettime(syscall.CLOCK_REALTIME, &tp)
	return int(tp.Sec)
}

func UnixNs() int64 {
	var tp syscall.Timespec
	syscall.ClockGettime(syscall.CLOCK_REALTIME, &tp)
	return tp.Sec*1_000_000_000 + tp.Nsec
}

func ToTm(t int) Tm {
	var tm Tm

	daysSinceJan1st := [2][13]int{
		{0, 31, 59, 90, 120, 151, 181, 212, 243, 273, 304, 334, 365}, // 365 days, non-leap
		{0, 31, 60, 91, 121, 152, 182, 213, 244, 274, 305, 335, 366}, // 366 days, leap
	}

	t += 3 * 60 * 60 /* MSK timezone hack. */

	/* Re-bias from 1970 to 1601: 1970 - 1601 = 369 = 3*100 + 17*4 + 1 years (incl. 89 leap days) = (3*100*(365+24/100) + 17*4*(365+1/4) + 1*365)*24*3600 seconds. */
	sec := t + 11644473600

	wday := (sec/86400 + 1) % 7 /* day of week */

	/* Remove multiples of 400 years (incl. 97 leap days). */
	quadricentennials := sec / 12622780800 /* 400*365.2425*24*3600 .*/
	sec %= 12622780800

	/* Remove multiples of 100 years (incl. 24 leap days), can't be more than 3 (because multiples of 4*100=400 years (incl. leap days) have been removed). */
	centennials := sec / 3155673600 /* 100*(365+24/100)*24*3600. */
	if centennials > 3 {
		centennials = 3
	}
	sec -= centennials * 3155673600

	/* Remove multiples of 4 years (incl. 1 leap day), can't be more than 24 (because multiples of 25*4=100 years (incl. leap days) have been removed). */
	quadrennials := sec / 126230400 /*  4*(365+1/4)*24*3600. */
	if quadrennials > 24 {
		quadrennials = 24
	}
	sec -= quadrennials * 126230400

	/* Remove multiples of years (incl. 0 leap days), can't be more than 3 (because multiples of 4 years (incl. leap days) have been removed). */
	annuals := sec / 31536000 // 365*24*3600
	if annuals > 3 {
		annuals = 3
	}
	sec -= annuals * 31536000

	/* Calculate the year and find out if it's leap. */
	year := 1601 + quadricentennials*400 + centennials*100 + quadrennials*4 + annuals
	var leap int
	if (year%4 == 0) && ((year%100 != 0) || (year%400 == 0)) {
		leap = 1
	} else {
		leap = 0
	}

	/* Calculate the day of the year and the time. */
	yday := sec / 86400
	sec %= 86400
	hour := sec / 3600
	sec %= 3600
	min := sec / 60
	sec %= 60

	/* Calculate the month. */
	var month, mday int = 1, 1
	for ; month < 13; month++ {
		if yday < daysSinceJan1st[leap][month] {
			mday += yday - daysSinceJan1st[leap][month-1]
			break
		}
	}

	tm.Sec = sec          /*  [0,59]. */
	tm.Min = min          /*  [0,59]. */
	tm.Hour = hour        /*  [0,23]. */
	tm.Mday = mday        /*  [1,31]  (day of month). */
	tm.Mon = month - 1    /*  [0,11]  (month). */
	tm.Year = year - 1900 /*  70+     (year since 1900). */
	tm.Wday = wday        /*  [0,6]   (day since Sunday AKA day of week). */
	tm.Yday = yday        /*  [0,365] (day since January 1st AKA day of year). */
	tm.Isdst = -1         /*  daylight saving time flag. */

	return tm
}

func PutTmRFC822(buf []byte, tm Tm) int {
	var n, ndigits int

	var wdays = [...]string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	var months = [...]string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}

	n += copy(buf[n:], wdays[tm.Wday])
	buf[n] = ','
	buf[n+1] = ' '
	n += 2

	if tm.Mday < 10 {
		buf[n] = '0'
		n++
	}
	ndigits = slices.PutInt(buf[n:], tm.Mday)
	n += ndigits
	buf[n] = ' '
	n++

	n += copy(buf[n:], months[tm.Mon])
	buf[n] = ' '
	n++

	ndigits = slices.PutInt(buf[n:], tm.Year+1900)
	n += ndigits
	buf[n] = ' '
	n++

	if tm.Hour < 10 {
		buf[n] = '0'
		n++
	}
	ndigits = slices.PutInt(buf[n:], tm.Hour)
	n += ndigits
	buf[n] = ':'
	n++

	if tm.Min < 10 {
		buf[n] = '0'
		n++
	}
	ndigits = slices.PutInt(buf[n:], tm.Min)
	n += ndigits
	buf[n] = ':'
	n++

	if tm.Sec < 10 {
		buf[n] = '0'
		n++
	}
	ndigits = slices.PutInt(buf[n:], tm.Sec)
	n += ndigits
	buf[n] = ' '
	n++

	buf[n] = '+'
	buf[n+1] = '0'
	buf[n+2] = '3'
	buf[n+3] = '0'
	buf[n+4] = '0'

	return n + 5
}
