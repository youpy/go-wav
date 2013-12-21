package wav

import (
	"encoding/binary"
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
	var value16 int16
	var value8 int8

	format := wav.Format
	numChannels := format.NumChannels
	bitsPerSample := format.BitsPerSample

	values := make([]int, format.NumChannels)

	for i := 0; i < int(numChannels); i++ {
		if bitsPerSample == 16 {
			err = binary.Read(wav.WavData, binary.LittleEndian, &value16)
			if err != nil {
				return
			}

			values[i] = int(value16)
		} else { // 8bit
			err = binary.Read(wav.WavData, binary.LittleEndian, &value8)
			if err != nil {
				return
			}

			values[i] = int(value16)
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
