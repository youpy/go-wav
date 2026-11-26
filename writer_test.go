package wav

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrite(t *testing.T) {
	outfile, err := ioutil.TempFile("/tmp", "outfile")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		outfile.Close()
		os.Remove(outfile.Name())
	}()

	var numSamples uint32 = 2
	var numChannels uint16 = 2
	var sampleRate uint32 = 44100
	var bitsPerSample uint16 = 16

	writer := NewWriter(outfile, numSamples, numChannels, sampleRate, bitsPerSample)
	samples := make([]Sample, numSamples)

	samples[0].Values[0] = 32767
	samples[0].Values[1] = -32768
	samples[1].Values[0] = 123
	samples[1].Values[1] = -123

	err = writer.WriteSamples(samples)
	if err != nil {
		t.Fatal(err)
	}

	outfile.Close()
	file, err := os.Open(outfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		file.Close()
		os.Remove(outfile.Name())
	}()

	reader := NewReader(file)
	if err != nil {
		t.Fatal(err)
	}

	fmt, err := reader.Format()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, int(fmt.AudioFormat), AudioFormatPCM)
	assert.Equal(t, fmt.NumChannels, numChannels)
	assert.Equal(t, fmt.SampleRate, sampleRate)
	assert.Equal(t, fmt.ByteRate, sampleRate*4)
	assert.Equal(t, fmt.BlockAlign, numChannels*(bitsPerSample/8))
	assert.Equal(t, fmt.BitsPerSample, bitsPerSample)

	samples, err = reader.ReadSamples()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(samples), 2)
	assert.Equal(t, samples[0].Values[0], 32767)
	assert.Equal(t, samples[0].Values[1], -32768)
	assert.Equal(t, samples[1].Values[0], 123)
	assert.Equal(t, samples[1].Values[1], -123)
}

func TestWrite8BitStereo(t *testing.T) {
	outfile, err := ioutil.TempFile("/tmp", "outfile")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		outfile.Close()
		os.Remove(outfile.Name())
	}()

	var numSamples uint32 = 2
	var numChannels uint16 = 2
	var sampleRate uint32 = 44100
	var bitsPerSample uint16 = 8

	writer := NewWriter(outfile, numSamples, numChannels, sampleRate, bitsPerSample)
	samples := make([]Sample, numSamples)

	samples[0].Values[0] = 255
	samples[0].Values[1] = 0
	samples[1].Values[0] = 123
	samples[1].Values[1] = 234

	err = writer.WriteSamples(samples)
	if err != nil {
		t.Fatal(err)
	}

	outfile.Close()
	file, err := os.Open(outfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		file.Close()
		os.Remove(outfile.Name())
	}()

	reader := NewReader(file)
	if err != nil {
		t.Fatal(err)
	}

	fmt, err := reader.Format()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, int(fmt.AudioFormat), AudioFormatPCM)
	assert.Equal(t, fmt.NumChannels, numChannels)
	assert.Equal(t, fmt.SampleRate, sampleRate)
	assert.Equal(t, fmt.ByteRate, sampleRate*2)
	assert.Equal(t, fmt.BlockAlign, numChannels*(bitsPerSample/8))
	assert.Equal(t, fmt.BitsPerSample, bitsPerSample)

	samples, err = reader.ReadSamples()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(samples), 2)
	assert.Equal(t, samples[0].Values[0], 255)
	assert.Equal(t, samples[0].Values[1], 0)
	assert.Equal(t, samples[1].Values[0], 123)
	assert.Equal(t, samples[1].Values[1], 234)
}

func TestWrite24BitStereo(t *testing.T) {
	outfile, err := ioutil.TempFile("/tmp", "outfile")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		outfile.Close()
		os.Remove(outfile.Name())
	}()

	var numSamples uint32 = 2
	var numChannels uint16 = 2
	var sampleRate uint32 = 44100
	var bitsPerSample uint16 = 24

	writer := NewWriter(outfile, numSamples, numChannels, sampleRate, bitsPerSample)
	samples := make([]Sample, numSamples)

	samples[0].Values[0] = 32767
	samples[0].Values[1] = -32768
	samples[1].Values[0] = 123
	samples[1].Values[1] = -123

	err = writer.WriteSamples(samples)
	if err != nil {
		t.Fatal(err)
	}

	outfile.Close()
	file, err := os.Open(outfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		file.Close()
		os.Remove(outfile.Name())
	}()

	reader := NewReader(file)
	if err != nil {
		t.Fatal(err)
	}

	fmt, err := reader.Format()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, int(fmt.AudioFormat), AudioFormatPCM)
	assert.Equal(t, fmt.NumChannels, numChannels)
	assert.Equal(t, fmt.SampleRate, sampleRate)
	assert.Equal(t, fmt.ByteRate, sampleRate*6)
	assert.Equal(t, fmt.BlockAlign, numChannels*(bitsPerSample/8))
	assert.Equal(t, fmt.BitsPerSample, bitsPerSample)

	samples, err = reader.ReadSamples()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(samples), 2)
	assert.Equal(t, samples[0].Values[0], 32767)
	assert.Equal(t, samples[0].Values[1], -32768)
	assert.Equal(t, samples[1].Values[0], 123)
	assert.Equal(t, samples[1].Values[1], -123)
}

func TestWrite32BitStereo(t *testing.T) {
	outfile, err := ioutil.TempFile("/tmp", "outfile")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		outfile.Close()
		os.Remove(outfile.Name())
	}()

	var numSamples uint32 = 2
	var numChannels uint16 = 2
	var sampleRate uint32 = 44100
	var bitsPerSample uint16 = 32

	writer := NewWriter(outfile, numSamples, numChannels, sampleRate, bitsPerSample)
	samples := make([]Sample, numSamples)

	samples[0].Values[0] = 32767
	samples[0].Values[1] = -32768
	samples[1].Values[0] = 123
	samples[1].Values[1] = -123

	err = writer.WriteSamples(samples)
	if err != nil {
		t.Fatal(err)
	}

	outfile.Close()
	file, err := os.Open(outfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		file.Close()
		os.Remove(outfile.Name())
	}()

	reader := NewReader(file)
	if err != nil {
		t.Fatal(err)
	}

	fmt, err := reader.Format()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, int(fmt.AudioFormat), AudioFormatPCM)
	assert.Equal(t, fmt.NumChannels, numChannels)
	assert.Equal(t, fmt.SampleRate, sampleRate)
	assert.Equal(t, fmt.ByteRate, sampleRate*8)
	assert.Equal(t, fmt.BlockAlign, numChannels*(bitsPerSample/8))
	assert.Equal(t, fmt.BitsPerSample, bitsPerSample)

	samples, err = reader.ReadSamples()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(samples), 2)
	assert.Equal(t, samples[0].Values[0], 32767)
	assert.Equal(t, samples[0].Values[1], -32768)
	assert.Equal(t, samples[1].Values[0], 123)
	assert.Equal(t, samples[1].Values[1], -123)
}

type failingDestination struct {
	bytes.Buffer
	err   error
	short bool
}

func (w *failingDestination) Write(p []byte) (int, error) {
	if w.err != nil {
		return 0, w.err
	}
	if w.short {
		return len(p) - 1, nil
	}
	return w.Buffer.Write(p)
}

func TestWriteSamplesReportsDestinationErrors(t *testing.T) {
	wantErr := errors.New("destination write failed")
	for _, count := range []int{1, 2048, 2049} {
		t.Run(fmt.Sprintf("write_error_%d_bytes", count*2), func(t *testing.T) {
			dst := &failingDestination{}
			writer := NewWriter(dst, uint32(count), 1, 8000, 16)
			dst.err = wantErr
			if err := writer.WriteSamples(make([]Sample, count)); !errors.Is(err, wantErr) {
				t.Fatalf("WriteSamples error = %v, want %v", err, wantErr)
			}
		})
	}
	for _, count := range []int{1, 2048} {
		t.Run(fmt.Sprintf("short_write_%d_bytes", count*2), func(t *testing.T) {
			dst := &failingDestination{}
			writer := NewWriter(dst, uint32(count), 1, 8000, 16)
			dst.short = true
			if err := writer.WriteSamples(make([]Sample, count)); !errors.Is(err, io.ErrShortWrite) {
				t.Fatalf("WriteSamples error = %v, want %v", err, io.ErrShortWrite)
			}
		})
	}
}

func TestWriteSamplesSuccessfulDestination(t *testing.T) {
	for _, count := range []int{0, 1, 2048, 2049} {
		t.Run(fmt.Sprintf("%d_bytes", count*2), func(t *testing.T) {
			dst := &failingDestination{}
			writer := NewWriter(dst, uint32(count), 1, 8000, 16)
			headerSize := dst.Len()
			samples := make([]Sample, count)
			for i := range samples {
				samples[i].Values[0] = 0x1234
			}
			if err := writer.WriteSamples(samples); err != nil {
				t.Fatal(err)
			}
			if got, want := dst.Bytes()[headerSize:], bytes.Repeat([]byte{0x34, 0x12}, count); !bytes.Equal(got, want) {
				t.Fatalf("encoded sample bytes differ: got %d bytes, want %d", len(got), len(want))
			}
		})
	}
}

func BenchmarkWriteSamples(b *testing.B) {
	n := 4096
	samples := []Sample{}

	file, _ := os.Open("./files/a.wav")
	reader := NewReader(file)

	for {
		spls, err := reader.ReadSamples(uint32(n))
		if err == io.EOF {
			break
		}
		samples = append(samples, spls...)
	}

	b.Run("write samples", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			tmpfile, err := ioutil.TempFile("", "example")
			if err != nil {
				b.Fatal(err)
			}
			defer os.Remove(tmpfile.Name())
			writer := NewWriter(tmpfile, uint32(len(samples)), 2, 44100, 16)
			writer.WriteSamples(samples)
		}
	})
}
