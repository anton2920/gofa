package oss

import (
	"errors"
	"unsafe"

	"github.com/anton2920/gofa/context"
	"github.com/anton2920/gofa/os/posix/freebsd"
)

type AudioBufInfo struct {
	Fragments  int32 /* # of avail. frags (partly used ones not counted) */
	FragsTotal int32 /* Total # of fragments allocated */
	FragSize   int32 /* Size of a fragment in bytes */
	Bytes      int32 /* Avail. space in bytes (includes partly used fragments). Note! 'bytes' could be more than fragments*fragsize */
}

var (
	SNDCTL_DSP_GETOSPACE = freebsd.IOR('P', 12, uint(unsafe.Sizeof(AudioBufInfo{})))
	SNDCTL_DSP_GETISPACE = freebsd.IOR('P', 13, uint(unsafe.Sizeof(AudioBufInfo{})))
)

func (d *Device) OutputBufferAvailableSpace() (int, error) {
	var ab AudioBufInfo
	if !freebsd.Ioctl(&context.Context{}, int32(d.Handle), SNDCTL_DSP_GETOSPACE, unsafe.Pointer(&ab)) {
		return -1, errors.New("failed to get output audio buf info")
	}
	return int(ab.Bytes), nil
}
