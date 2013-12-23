package wav

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
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
