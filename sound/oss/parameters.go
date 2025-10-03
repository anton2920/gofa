package oss

import (
	"fmt"
	"unsafe"

	"github.com/anton2920/gofa/audio/wave"
	"github.com/anton2920/gofa/ints"
	"github.com/anton2920/gofa/syscall"
)

type DeviceParameters struct {
	Format       int
	Channels     int
	SamplingRate int
}

/* From <sys/soundcard.h>. */
const (
	// FormatQuery    = 0x00000000 /* Return current format */
	FormatMuLaw    = 0x00000001 /* Logarithmic mu-law */
	FormatALaw     = 0x00000002 /* Logarithmic A-law */
	FormatIMAADPCM = 0x00000004 /* A 4:1 compressed format where 16-bit squence represented using the the average 4 bits per sample */
	FormatU8       = 0x00000008 /* Unsigned 8-bit */
	FormatS16LE    = 0x00000010 /* Little endian signed 16-bit */
	FormatS16BE    = 0x00000020 /* Big endian signed 16-bit */
	FormatS8       = 0x00000040 /* Signed 8-bit */
	FormatU16LE    = 0x00000080 /* Little endian unsigned 16-bit */
	FormatU16BE    = 0x00000100 /* Big endian unsigned 16-bit */
	FormatMPEG     = 0x00000200 /* MPEG MP2/MP3 audio */
	FormatAC3      = 0x00000400 /* Dolby Digital AC3 */

	/*
	 * 32-bit formats below used for 24-bit audio data where the data is stored
	 * in the 24 most significant bits and the least significant bits are not used
	 * (should be set to 0).
	 */
	FormatS32LE = 0x00001000 /* Little endian signed 32-bit */
	FormatS32BE = 0x00002000 /* Big endian signed 32-bit */
	FormatU32LE = 0x00004000 /* Little endian unsigned 32-bit */
	FormatU32BE = 0x00008000 /* Big endian unsigned 32-bit */
	FormatS24LE = 0x00010000 /* Little endian signed 24-bit */
	FormatS24BE = 0x00020000 /* Big endian signed 24-bit */
	FormatU24LE = 0x00040000 /* Little endian unsigned 24-bit */
	FormatU24BE = 0x00080000 /* Big endian unsigned 24-bit */
)

func MergeDeviceParameters(params ...DeviceParameters) DeviceParameters {
	var result DeviceParameters

	for i := 0; i < len(params); i++ {
		param := &params[i]

		ints.Replace(&result.Format, param.Format)
		ints.Replace(&result.Channels, param.Channels)
		ints.Replace(&result.SamplingRate, param.SamplingRate)
	}

	return result
}

func SetDeviceParameters(fd int32, params DeviceParameters) error {
	format := int32(params.Format)
	channels := int32(params.Channels)
	speed := int32(params.SamplingRate)

	if err := syscall.Ioctl(fd, SNDCTL_DSP_SETFMT, unsafe.Pointer(&format)); err != nil {
		return fmt.Errorf("failed to set sample format: %v", err)
	}
	if err := syscall.Ioctl(fd, SNDCTL_DSP_CHANNELS, unsafe.Pointer(&channels)); err != nil {
		return fmt.Errorf("failed to set number of channels: %v", err)
	}
	if err := syscall.Ioctl(fd, SNDCTL_DSP_SPEED, unsafe.Pointer(&speed)); err != nil {
		return fmt.Errorf("failed to set sampling rate: %v", err)
	}

	return nil
}

func Channels(channels int) DeviceParameters {
	return DeviceParameters{Channels: channels}
}

func Format(format int) DeviceParameters {
	return DeviceParameters{Format: format}
}

func SamplingRate(hz int) DeviceParameters {
	return DeviceParameters{SamplingRate: hz}
}

func WaveHeader(header wave.Header) DeviceParameters {
	var result DeviceParameters

	result.Channels = int(header.NumChannels)
	result.SamplingRate = int(header.SampleRate)

	switch header.BitsPerSample {
	case 8:
		result.Format = FormatS8
	case 16:
		result.Format = FormatS16LE
	case 24:
		result.Format = FormatS24LE
	case 32:
		result.Format = FormatS32LE
	}

	return result
}
