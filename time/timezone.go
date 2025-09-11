package time

type Timezone int32

const (
	TimezoneNone = Timezone(iota - 13)
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	TimezoneUTC
	_
	_
	TimezoneMSK
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	TimezoneCount
)

var Timezone2String = map[int]string{
	-12: "UTC-12",
	-11: "UTC-11",
	-10: "UTC-10",
	-9:  "UTC-9",
	-8:  "UTC-8",
	-7:  "UTC-7",
	-6:  "UTC-6",
	-5:  "UTC-5",
	-4:  "UTC-4",
	-3:  "UTC-3",
	-2:  "UTC-2",
	-1:  "UTC-1",
	0:   "UTC+0",
	1:   "UTC+1",
	2:   "UTC+2",
	3:   "UTC+3",
	4:   "UTC+4",
	5:   "UTC+5",
	6:   "UTC+6",
	7:   "UTC+7",
	8:   "UTC+8",
	9:   "UTC+9",
	10:  "UTC+10",
	11:  "UTC+11",
	12:  "UTC+12",
	13:  "UTC+13",
	14:  "UTC+14",
	15:  "UTC+15",
	16:  "UTC+16",
}
