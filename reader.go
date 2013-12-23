package wav

import (
	"bufio"
	"encoding/binary"
	"errors"
	riff "github.com/youpy/go-riff"
	"math"
)

type Reader struct {
	r         *riff.Reader
	riffChunk *riff.RIFFChunk
	format    *WavFormat
	*WavData
}

func NewReader(r riff.RIFFReader) *Reader {
	riffReader := riff.NewReader(r)
	return &Reader{r: riffReader}
}

func (r *Reader) Format() (format *WavFormat, err error) {
	if r.format == nil {
		format, err = r.readFormat()
		if err != nil {
			return
		}
		r.format = format
	} else {
		format = r.format
	}

	return
}

func (r *Reader) Read(p []byte) (n int, err error) {
	if r.WavData == nil {
		data, err := r.readData()
		if err != nil {
			return n, err
		}
		r.WavData = data
	}

	return r.WavData.Read(p)
}

func (r *Reader) ReadSamples(params ...uint32) (samples []Sample, err error) {
	var bytes []byte
	var numSamples, n int

	if len(params) > 0 {
		numSamples = int(params[0])
	} else {
		numSamples = 2048
	}

	format, err := r.Format()
	if err != nil {
		return
	}

	numChannels := int(format.NumChannels)
	blockAlign := int(format.BlockAlign)
	bitsPerSample := int(format.BitsPerSample)

	bytes = make([]byte, numSamples*blockAlign)
	n, err = r.Read(bytes)

	if err != nil {
		return
	}

	numSamples = n / blockAlign
	r.WavData.pos += uint32(numSamples * blockAlign)
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

func (r *Reader) IntValue(sample Sample, channel uint) int {
	return sample.Values[channel]
}

func (r *Reader) FloatValue(sample Sample, channel uint) float64 {
	// XXX
	return float64(r.IntValue(sample, channel)) / math.Pow(2, float64(r.format.BitsPerSample))
}

func (r *Reader) readFormat() (fmt *WavFormat, err error) {
	var riffChunk *riff.RIFFChunk

	fmt = new(WavFormat)

	if r.riffChunk == nil {
		riffChunk, err = r.r.Read()
		if err != nil {
			return
		}

		r.riffChunk = riffChunk
	} else {
		riffChunk = r.riffChunk
	}

	fmtChunk := findChunk(riffChunk, "fmt ")

	if fmtChunk == nil {
		err = errors.New("Format chunk is not found")
		return
	}

	err = binary.Read(fmtChunk, binary.LittleEndian, fmt)
	if err != nil {
		return
	}

	return
}

func (r *Reader) readData() (data *WavData, err error) {
	var riffChunk *riff.RIFFChunk

	if r.riffChunk == nil {
		riffChunk, err = r.r.Read()
		if err != nil {
			return
		}

		r.riffChunk = riffChunk
	} else {
		riffChunk = r.riffChunk
	}

	dataChunk := findChunk(riffChunk, "data")
	if dataChunk == nil {
		err = errors.New("Data chunk is not found")
		return
	}

	data = &WavData{bufio.NewReader(dataChunk), dataChunk.ChunkSize, 0}

	return
}

func findChunk(riffChunk *riff.RIFFChunk, id string) (chunk *riff.Chunk) {
	for _, ch := range riffChunk.Chunks {
		if string(ch.ChunkID[:]) == id {
			chunk = ch
			break
		}
	}

	return
}
