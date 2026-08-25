package intel

import (
	"unsafe"

	"github.com/anton2920/gofa/bytes"
	"github.com/anton2920/gofa/debug/debug_"
	"github.com/anton2920/gofa/ints"
	"github.com/anton2920/gofa/slices"
	"github.com/anton2920/gofa/strings"
)

var (
	HighestBasicFunction    uint32
	HighestExtendedFunction uint32

	VendorString string

	Stepping      int
	Model         int /* this contains both Model and ExtendedModel fields. */
	Family        int /* this contains both Family and ExtendedFamily fields. */
	ProcessorType int /* 00 - Original OEM Processor, 01 - Intel OverDrive® Processor, 10 - Dual processor. */

	BrandIndex  int
	BrandString string

	CPUHz uint64
)

var BrandIndex2BrandString = [...]string{
	0x01: "Intel(R) Celeron(R) processor",
	0x02: "Intel(R) Pentium(R) III processor",
	0x03: "Intel(R) Pentium(R) III Xeon(R) processor", /* if processor signature = 0x000006B1, then "Intel(R) Celeron(R) processor". */
	0x04: "Intel(R) Pentium(R) III processor",
	0x06: "Mobile Intel(R) Pentium(R) III processor-M",
	0x07: "Mobile Intel(R) Celeron(R) processor",
	0x08: "Intel(R) Pentium(R) 4 processor",
	0x09: "Intel(R) Pentium(R) 4 processor",
	0x0A: "Intel(R) Celeron(R) processor",
	0x0B: "Intel(R) Xeon(R) processor", /* if processor signature = 0x00000F13, "Intel(R) Xeon(R) processor". */
	0x0C: "Intel(R) Xeon(R) processor MP",
	0x0E: "Mobile Intel(R) Pentium(R) 4 processor-M", /* if processor signature = 0x00000F13, then "Intel(R) Xeon(R) processor" */
	0x0F: "Mobile Intel(R) Celeron(R) processor",
	0x11: "Mobile Genuine Intel(R) processor",
	0x12: "Intel(R) Celeron(R) M processor",
	0x13: "Mobile Intel(R) Celeron processor",
	0x14: "Intel(R) Celeron(R) processor",
	0x15: "Mobile Genuine Intel(R) processor",
	0x16: "Intel(R) Pentium(R) M processor",
	0x17: "Mobile Intel(R) Celeron(R) processor",
}

func init() {
	{
		a, b, c, d := CPUID(0x0, 0)
		HighestBasicFunction = a

		vendor := make([]uint32, 3)
		vendor[0] = b
		vendor[1] = d
		vendor[2] = c
		buffer := bytes.SliceFromUnsafePointer(unsafe.Pointer(&vendor[0]), len(vendor)*int(unsafe.Sizeof(vendor[0])))
		VendorString = string(buffer)
	}
	{
		info, index, _, _ := CPUID(0x1, 0)
		Stepping = int(info & 0xF)
		Family = int((info >> 8) & 0xF)
		if Family == 0 {
			Family += int((info >> 20) & 0x1F)
		}
		Model = int((info >> 4) & 0xF)
		if (Family == 0x6) || (Family == 0x0) {
			Model += (int((info>>16)&0xF) << 4)
		}
		ProcessorType = int((info >> 12) & 0x3)
		BrandIndex = int(index & 0xFF)
	}
	{
		hexdigits0 := [...]string{"", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F"}
		hexdigits := [...]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F"}
		bindigits := [...]string{"00", "01", "10", "11"}

		debug_.Println("[cpu/intel]: ", VendorString,
			" Family ", hexdigits0[(Family>>4)&0xF], hexdigits[Family&0xF],
			" Model ", hexdigits0[(Model>>4)&0xF], hexdigits[Model&0xF],
			" Stepping ", hexdigits0[(Stepping>>4)&0xF], hexdigits[Stepping&0xF],
			" Type ", bindigits[ProcessorType],
		)
	}

	{
		HighestExtendedFunction, _, _, _ = CPUID(0x80000000, 0)
		if (HighestExtendedFunction & 0x80000000) == 0 {
			/* Brand index method. */
			BrandString = BrandIndex2BrandString[BrandIndex]

		} else if HighestExtendedFunction >= 0x80000004 {
			/* Brand string method. */
			base := uint32(0x80000002)
			brand := make([]uint32, 12)
			for i := 0; i < 3; i++ {
				a, b, c, d := CPUID(base, 0)
				brand[4*i+0] = a
				brand[4*i+1] = b
				brand[4*i+2] = c
				brand[4*i+3] = d
				base++
			}
			buffer := bytes.SliceFromUnsafePointer(unsafe.Pointer(&brand[0]), len(brand)*int(unsafe.Sizeof(brand[0])))
			BrandString = string(buffer)

			/* Trimming '\x00' byte at the end. */
			for BrandString[len(BrandString)-1] == '\x00' {
				BrandString = BrandString[:len(BrandString)-1]
			}
		}
		if len(BrandString) > 0 {
			debug_.Println("[cpu/intel]: ", strings.TrimSpace(BrandString))
		}
	}

	{
		denominator, numerator, coreHz, _ := CPUID(0x15, 0)
		if (numerator != 0) && (coreHz != 0) {
			CPUHz = (uint64(coreHz) * uint64(numerator)) / uint64(denominator)
		} else if (denominator != 0) && (coreHz == 0) {
			signature := (Family << 8) | (Model)
			switch signature {
			case 0x0655: /* Intel® Xeon® Scalable Processor Family. */
				CPUHz = (uint64(25000000) * uint64(numerator)) / uint64(denominator)
			case 0x064E, 0x065E, 0x068E, 0x069E: /* 6th, 7th, 8th and 9th generation Intel® Core™ processors. */
				CPUHz = (uint64(24000000) * uint64(numerator)) / uint64(denominator)
			case 0x65C: /* Next Generation Intel Atom® processors based on Goldmont Microarchitecture. */
				CPUHz = (uint64(19200000) * uint64(numerator)) / uint64(denominator)
			}
		}
		if CPUHz > 0 {
			buf := make([]byte, ints.Bufsize)
			n := slices.PutUint64(buf, CPUHz)
			debug_.Println("[cpu/intel]: CPU frequency: ", bytes.AsString(buf[:n]), "Hz")
		}
	}
}
