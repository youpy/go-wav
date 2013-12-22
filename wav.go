package wav

import (
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
	pos  uint32
}

type Wav struct {
	*WavData
	Format *WavFormat
}

type Sample struct {
	Values [2]int
}

func (wav *Wav) ReadSamples(params ...uint32) (samples []Sample, err error) {
	var bytes []byte
	var numSamples, n int

	if len(params) > 0 {
		numSamples = int(params[0])
	} else {
		numSamples = 2048
	}

	format := wav.Format
	numChannels := uint32(format.NumChannels)
	blockAlign := uint32(format.BlockAlign)
	bitsPerSample := format.BitsPerSample

	bytes = make([]byte, numChannels*uint32(numSamples)*(uint32(bitsPerSample)/8))
	n, err = wav.WavData.Read(bytes)

	if err != nil {
		return
	}

	numSamples = n / int(blockAlign)
	wav.WavData.pos += uint32(numSamples) * blockAlign
	samples = make([]Sample, numSamples)

	var offset int = 0

	for i := 0; i < numSamples; i++ {
		if bitsPerSample == 16 {
			for j := 0; j < int(numChannels); j++ {
				samples[i].Values[j] = int((int16(bytes[offset+(j*2)+1]) << 8) + int16(bytes[offset+(j*2)]))
			}
		} else {
			for j := 0; j < int(numChannels); j++ {
				samples[i].Values[j] = int(bytes[offset+j])
			}
		}

		offset += int(numChannels) * (int(bitsPerSample) / 8)
	}

	return
}

func (wav *Wav) IntValue(sample Sample, channel uint) int {
	return sample.Values[channel]
}

func (wav *Wav) FloatValue(sample Sample, channel uint) float64 {
	return float64(wav.IntValue(sample, channel)) / math.Pow(2, float64(wav.Format.BitsPerSample))
}
