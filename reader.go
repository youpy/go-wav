package wav

import (
	"bytes"
	"encoding/binary"
	"errors"
	riff "github.com/youpy/go-riff"
	"io/ioutil"
)

type Reader struct {
	r         *riff.Reader
	riffChunk *riff.RIFFChunk
}

func NewReader(r riff.RIFFReader) *Reader {
	riffReader := riff.NewReader(r)
	return &Reader{r: riffReader}
}

func (r *Reader) Read() (wav *Wav, err error) {
	fmt, err := r.readFormat()
	if err != nil {
		return
	}

	data, err := r.readData()
	if err != nil {
		return
	}

	wav = &Wav{data, fmt}

	return
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

	b, err := ioutil.ReadAll(fmtChunk)
	if err != nil {
		return
	}

	buf := bytes.NewReader(b)
	err = binary.Read(buf, binary.LittleEndian, &fmt.AudioFormat)
	if err != nil {
		return
	}

	err = binary.Read(buf, binary.LittleEndian, &fmt.NumChannels)
	if err != nil {
		return
	}

	err = binary.Read(buf, binary.LittleEndian, &fmt.SampleRate)
	if err != nil {
		return
	}

	err = binary.Read(buf, binary.LittleEndian, &fmt.ByteRate)
	if err != nil {
		return
	}

	err = binary.Read(buf, binary.LittleEndian, &fmt.BlockAlign)
	if err != nil {
		return
	}

	err = binary.Read(buf, binary.LittleEndian, &fmt.BitsPerSample)
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

	data = &WavData{dataChunk}

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
