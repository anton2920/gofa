package oss

import (
	"fmt"

	"github.com/anton2920/gofa/os"
	"github.com/anton2920/gofa/os/posix/freebsd"
)

type Device struct {
	Handle os.Handle
	Params DeviceParameters
}

type Mode int32

const (
	ModeInput = Mode(iota)
	ModeOutput
	ModeInputOutput
)

func Open(path string, mode Mode, params ...DeviceParameters) (Device, error) {
	var d Device

	fd, err := freebsd.Open(path, int32(mode), 0)
	if err != nil {
		return d, fmt.Errorf("failed ot open audio device: %v", err)
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
		freebsd.Close(fd)
		return d, fmt.Errorf("failed to set device parameters: %v", err)
	}

	d.Handle = os.Handle(fd)
	d.Params = result

	return d, nil
}

func (d *Device) Close() error {
	return os.CloseHandle(d.Handle)
}

func (d *Device) Write(buf []byte) (int, error) {
	return os.WriteToFile(d.Handle, buf)
}
