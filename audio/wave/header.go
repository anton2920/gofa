package wave

type RiffChunk struct {
	ID     [4]byte /* ASCII "RIFF". */
	Size   uint32  /* 36 + SubChunk2Size. */
	Format [4]byte /* ASCII "WAVE". */
}

type FmtChunk struct {
	ID            [4]byte /* ASCII "fmt ". */
	Size          uint32  /* 16. */
	AudioFormat   uint16  /* 1. */
	NumChannels   uint16  /* Mono = 1, Stereo = 2, etc. */
	SampleRate    uint32  /* 8000, 44100, etc. */
	ByteRate      uint32  /* SampleRate * BlockAlign. */
	BlockAlign    uint16  /* NumChannels * BitsPerSample / 8. */
	BitsPerSample uint16  /* 8, 16, etc. */
}

type DataChunk struct {
	ID   [4]byte /* ASCII "data". */
	Size uint32  /* FMT.NumSamples * FMT.BlockAlign. */
}

type SimpleHeader struct {
	RIFF   RiffChunk
	Format FmtChunk
	Data   DataChunk
}
