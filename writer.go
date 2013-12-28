package wav

import (
	"encoding/binary"
	"github.com/youpy/go-riff"
	"io"
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

	var i uint16
	for _, sample := range samples {
		if bitsPerSample == 16 {
			for i = 0; i < numChannels; i++ {
				err = binary.Write(w, binary.LittleEndian, int16(sample.Values[i]))
				if err != nil {
					return
				}
			}
		} else {
			for i = 0; i < numChannels; i++ {
				binary.Write(w, binary.LittleEndian, int8(sample.Values[i]))
				if err != nil {
					return
				}
			}
		}
	}

	return
}
