package intel

import (
	"unsafe"

	"github.com/anton2920/gofa/debug"
)

var (
	HighestBasicFunction    uint32
	HighestExtendedFunction uint32

	VendorString string

	Stepping      int
	Model         int /* NOTE(anton2920): this contains both Model and ExtendedModel fields. */
	Family        int /* NOTE(anton2920): this contains both Family and ExtendedFamily fields. */
	ProcessorType int /* NOTE(anton2920): 00 - Original OEM Processor, 01 - Intel OverDrive® Processor, 10 - Dual processor. */

	BrandIndex  int
	BrandString string

	CPUHz Cycles
)

func init() {
	{
		a, b, c, d := CPUID(0x0, 0)
		HighestBasicFunction = a

		vendor := make([]uint32, 3)
		vendor[0] = b
		vendor[1] = d
		vendor[2] = c
		VendorString = string(unsafe.Slice((*byte)(unsafe.Pointer(&vendor[0])), len(vendor)*int(unsafe.Sizeof(vendor[0]))))
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
		debug.Printf("[gofa/intel]: %s Family %X Model %X Stepping %X Type %b", VendorString, Family, Model, Stepping, ProcessorType)

		BrandIndex = int(index & 0xFF)
	}

	{
		HighestExtendedFunction, _, _, _ = CPUID(0x80000000, 0)
		if HighestExtendedFunction > 0x80000004 {
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
			BrandString = string(unsafe.Slice((*byte)(unsafe.Pointer(&brand[0])), len(brand)*int(unsafe.Sizeof(brand[0]))))
		}
		debug.Printf("[gofa/intel]: %s", BrandString)
	}

	{
		denominator, numerator, coreHz, _ := CPUID(0x15, 0)
		if (numerator != 0) && (coreHz != 0) {
			CPUHz = (Cycles(coreHz) * Cycles(numerator)) / Cycles(denominator)
		} else if coreHz == 0 {
			signature := (Family << 8) | (Model)
			switch signature {
			case 0x0655: /* Intel® Xeon® Scalable Processor Family. */
				CPUHz = (Cycles(25_000_000) * Cycles(numerator)) / Cycles(denominator)
			case 0x064E, 0x065E, 0x068E, 0x069E: /* 6th, 7th, 8th and 9th generation Intel® Core™ processors. */
				CPUHz = (Cycles(24_000_000) * Cycles(numerator)) / Cycles(denominator)
			case 0x65C: /* Next Generation Intel Atom® processors based on Goldmont Microarchitecture. */
				CPUHz = (Cycles(19_200_000) * Cycles(numerator)) / Cycles(denominator)
			}
		}
		if CPUHz != 0 {
			debug.Printf("[gofa/intel]: CPU Frequency %dHz", CPUHz)
		}
	}
}
