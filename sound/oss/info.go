package oss

import (
	"fmt"
	"unsafe"

	"github.com/anton2920/gofa/syscall"
)

type AudioBufInfo struct {
	Fragments  int32 /* # of avail. frags (partly used ones not counted) */
	FragsTotal int32 /* Total # of fragments allocated */
	FragSize   int32 /* Size of a fragment in bytes */
	Bytes      int32 /* Avail. space in bytes (includes partly used fragments). Note! 'bytes' could be more than fragments*fragsize */
}

var (
	SNDCTL_DSP_GETOSPACE = syscall.IOR('P', 12, uint(unsafe.Sizeof(AudioBufInfo{})))
	SNDCTL_DSP_GETISPACE = syscall.IOR('P', 13, uint(unsafe.Sizeof(AudioBufInfo{})))
)

func (d *Device) OutputBufferAvailableSpace() (int, error) {
	var ab AudioBufInfo
	if err := syscall.Ioctl(int32(d.Handle), SNDCTL_DSP_GETOSPACE, unsafe.Pointer(&ab)); err != nil {
		return -1, fmt.Errorf("failed to get output audio buf info: %v", err)
	}
	return int(ab.Bytes), nil
}
