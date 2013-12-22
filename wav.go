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
	Pos  uint32
	Size uint32
}

type Wav struct {
	*WavData
	Format *WavFormat
}

type Sample struct {
	values8  []int8
	values16 []int16
}

func (wav *Wav) ReadSamples(params ...uint32) (samples []Sample, err error) {
	var values8 []int8
	var values16 []int16
	var n uint32

	if len(params) > 0 {
		n = params[0]
	} else {
		n = 2048
	}

	format := wav.Format
	numChannels := uint32(format.NumChannels)
	blockAlign := uint32(format.BlockAlign)
	bitsPerSample := format.BitsPerSample

	if wav.WavData.Size < wav.WavData.Pos+(n*blockAlign) {
		n = (wav.WavData.Size - wav.WavData.Pos) / blockAlign
	}

	if n == 0 {
		err = io.EOF
		return
	}

	if bitsPerSample == 16 {
		values16 = make([]int16, numChannels*n)
		err = binary.Read(wav.WavData, binary.LittleEndian, &values16)
	} else { // assumes 8bit
		values8 = make([]int8, numChannels*n)
		err = binary.Read(wav.WavData, binary.LittleEndian, &values8)
	}

	if err != nil {
		return
	}

	wav.WavData.Pos += n * blockAlign

	samples = make([]Sample, n)

	var i uint32
	for i = 0; i < n; i++ {
		if bitsPerSample == 16 {
			samples[i].values16 = values16[i*numChannels : i*numChannels+numChannels]
		} else {
			samples[i].values8 = values8[i*numChannels : i*numChannels+numChannels]
		}
	}

	return
}

func (wav *Wav) IntValue(sample Sample, channel uint) (value int) {
	if wav.Format.BitsPerSample == 16 {
		value = int(sample.values16[channel])
	} else {
		value = int(sample.values8[channel])
	}

	return
}

func (wav *Wav) FloatValue(sample Sample, channel uint) float64 {
	return float64(wav.IntValue(sample, channel)) / math.Pow(2, float64(wav.Format.BitsPerSample))
}
