package oss

import (
	"errors"

	"github.com/anton2920/gofa/context"
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

	fd, ok := freebsd.Open(&context.Context{}, path, int32(mode), 0)
	if !ok {
		return d, errors.New("failed ot open audio device")
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
		freebsd.Close(&context.Context{}, fd)
		return d, errors.New("failed to set device parameters")
	}

	d.Handle = os.Handle(fd)
	d.Params = result

	return d, nil
}

func (d *Device) Close() bool {
	return os.CloseHandle(&context.Context{}, d.Handle)
}

func (d *Device) Write(buf []byte) (int, bool) {
	return os.WriteToFile(&context.Context{}, d.Handle, buf)
}
