package wav

import (
	"bufio"
	"encoding/binary"
	"io"
	"math"

	"github.com/youpy/go-riff"
)

type Writer struct {
	io.Writer
	Format *WavFormat
}

func NewWriter(w io.Writer, numSamples uint32, numChannels uint16, sampleRate uint32, bitsPerSample uint16) (writer *Writer) {
	blockAlign := numChannels * bitsPerSample / 8
	byteRate := sampleRate * uint32(blockAlign)
	format := &WavFormat{AudioFormatPCM, numChannels, sampleRate, byteRate, blockAlign, bitsPerSample}
	dataSize := numSamples * uint32(format.BlockAlign)
	riffSize := 4 + 8 + 16 + 8 + dataSize
	riffWriter := riff.NewWriter(w, []byte("WAVE"), riffSize)

	writer = &Writer{riffWriter, format}
	riffWriter.WriteChunk([]byte("fmt "), 16, func(w io.Writer) {
		binary.Write(w, binary.LittleEndian, format)
	})
	riffWriter.WriteChunk([]byte("data"), dataSize, func(w io.Writer) {})

	return writer
}

func (w *Writer) WriteSamples(samples []Sample) (err error) {
	bitsPerSample := w.Format.BitsPerSample
	numChannels := w.Format.NumChannels
	bytesPerSample := int(bitsPerSample/8) * int(numChannels)
	by := make([]byte, 0, len(samples)*bytesPerSample)

	for _, sample := range samples {
		for i := uint16(0); i < numChannels; i++ {
			value := toUint(sample.Values[i], int(bitsPerSample))
			for b := uint16(0); b < bitsPerSample; b += 8 {
				by = append(by, uint8((value>>b)&math.MaxUint8))
			}
		}
	}

	bufWriter := bufio.NewWriter(w.Writer)
	_, err = bufWriter.Write(by)
	bufWriter.Flush()
	return err
}

func toUint(value int, bits int) uint {
	var result uint

	switch bits {
	case 32:
		result = uint(uint32(value))
	case 16:
		result = uint(uint16(value))
	case 8:
		result = uint(value)
	default:
		if value < 0 {
			result = uint((1 << uint(bits)) + value)
		} else {
			result = uint(value)
		}
	}

	return result
}
