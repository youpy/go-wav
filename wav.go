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
	numChannels := int(format.NumChannels)
	blockAlign := int(format.BlockAlign)
	bitsPerSample := int(format.BitsPerSample)

	bytes = make([]byte, numSamples*blockAlign)
	n, err = wav.WavData.Read(bytes)

	if err != nil {
		return
	}

	numSamples = n / blockAlign
	wav.WavData.pos += uint32(numSamples * blockAlign)
	samples = make([]Sample, numSamples)

	offset := 0

	for i := 0; i < numSamples; i++ {
		if bitsPerSample == 16 {
			for j := 0; j < int(numChannels); j++ {
				soffset := offset + (j * numChannels)
				samples[i].Values[j] = int((int16(bytes[soffset+1]) << 8) + int16(bytes[soffset]))
			}
		} else {
			for j := 0; j < int(numChannels); j++ {
				samples[i].Values[j] = int(bytes[offset+j])
			}
		}

		offset += blockAlign
	}

	return
}

func (wav *Wav) IntValue(sample Sample, channel uint) int {
	return sample.Values[channel]
}

func (wav *Wav) FloatValue(sample Sample, channel uint) float64 {
	return float64(wav.IntValue(sample, channel)) / math.Pow(2, float64(wav.Format.BitsPerSample))
}
