package oss

import (
	"fmt"
	"unsafe"

	"github.com/anton2920/gofa/os"
	"github.com/anton2920/gofa/syscall"
)

type Device struct {
	os.Handle

	DeviceParameters
}

var (
	SNDCTL_DSP_SETFMT   = syscall.IOWR('P', 5, uint(unsafe.Sizeof(int32(0))))
	SNDCTL_DSP_CHANNELS = syscall.IOWR('P', 6, uint(unsafe.Sizeof(int32(0))))
	SNDCTL_DSP_SPEED    = syscall.IOWR('P', 2, uint(unsafe.Sizeof(int32(0))))
)

func Open(path string, params ...DeviceParameters) (*Device, error) {
	var d Device

	fd, err := syscall.Open(path, syscall.O_WRONLY, 0)
	if err != nil {
		return nil, fmt.Errorf("failed ot open audio device: %v", err)
	}

	result := MergeDeviceParameters(params...)
	if result.Format == 0 {
		result.Format = FormatS16LE
	}
	if result.Channels == 0 {
		result.Channels = 2
	}
	if result.SamplingRate == 0 {
		result.SamplingRate = 44100
	}

	if err := SetDeviceParameters(fd, result); err != nil {
		syscall.Close(fd)
		return nil, fmt.Errorf("failed to set device parameters: %v", err)
	}

	d.Handle = os.Handle(fd)
	d.DeviceParameters = result

	return &d, nil
}

func (d *Device) Close() error {
	return os.Close(d.Handle)
}

func (d *Device) Write(buf []byte) (int64, error) {
	return os.Write(d.Handle, buf)
}
