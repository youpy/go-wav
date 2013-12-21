package wav

import (
	"errors"
	"io"
	"math"
)

const (
	AudioFormatPCM = 1
)

type WavFormat struct {
	AudioFormat   uint16
	NumChannels   uint16
	SampleRate    uint32
	ByteRate      uint32
	BlockAlign    uint16
	BitsPerSample uint16
}

type WavData struct {
	io.Reader
	Size uint32
}

type Wav struct {
	*WavData
	Format *WavFormat
}

type Sample struct {
	BitsPerSample uint16
	NumChannels   uint16
	values        []int
}

func (wav *Wav) ReadSample() (sample Sample, err error) {
	format := wav.Format
	numChannels := format.NumChannels
	bitsPerSample := format.BitsPerSample
	blockAlign := format.BlockAlign

	values := make([]int, format.NumChannels)
	bytes := make([]byte, format.BlockAlign)

	n, err := wav.WavData.Read(bytes)
	if err != nil {
		return
	}

	if n != int(blockAlign) {
		err = errors.New("Invalid length of bytes")
		return
	}

	for i := 0; i < int(numChannels); i++ {
		bytesForChannel := bytes[(i * int(bitsPerSample) / 8) : (i*int(bitsPerSample)/8)+(int(bitsPerSample)/8)]

		for j := 0; j < len(bytesForChannel); j++ {
			values[i] += int(bytesForChannel[j]) << uint(8*j)
		}
	}

	sample = Sample{bitsPerSample, numChannels, values}

	return
}

func (s Sample) IntValue(channel uint) int {
	return s.values[channel]
}

func (s Sample) FloatValue(channel uint) float64 {
	return float64(s.values[channel]) / math.Pow(2, float64(s.BitsPerSample))
}
