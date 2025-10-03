package wave

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/anton2920/gofa/os"
)

type Header struct {
	/* RIFF chunk descriptor. */
	ChunkID   [4]byte /* ASCII "RIFF". */
	ChunkSize uint32  /* 36 + SubChunk2Size. */
	Format    [4]byte /* ASCII "WAVE". */

	/* 'fmt' sub-chunk. */
	SubChunk1ID   [4]byte /* ASCII "fmt ". */
	SubChunk1Size uint32  /* 16. */
	AudioFormat   uint16  /* 1. */
	NumChannels   uint16  /* Mono = 1, Stereo = 2, etc. */
	SampleRate    uint32  /* 8000, 44100, etc. */
	ByteRate      uint32  /* SampleRate * BlockAlign. */
	BlockAlign    uint16  /* NumChannels * BitsPerSample / 8. */
	BitsPerSample uint16  /* 8, 16, etc. */

	/* 'data' sub-chunk. */
	SubChunk2ID   [4]byte /* ASCII "data". */
	SubChunk2Size uint32  /* NumSamples * BlockAlign. */
}

func ReadHeaderFromFile(f os.Handle, header *Header) error {
	n, err := os.Read(f, *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{Data: uintptr(unsafe.Pointer(header)), Len: int(unsafe.Sizeof(*header)), Cap: int(unsafe.Sizeof(*header))})))
	if err != nil {
		return fmt.Errorf("failed to read WAVE header from file: %v", err)
	}
	if n != int64(unsafe.Sizeof(*header)) {
		return fmt.Errorf("read incorrect number of bytes (%d != %d)", n, unsafe.Sizeof(*header))
	}
	return nil
}
