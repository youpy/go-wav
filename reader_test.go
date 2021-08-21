package wav

import (
	"io/ioutil"
	"math"
	"testing"

	"gotest.tools/assert"
)

func TestRead(t *testing.T) {
	blockAlign := 4

	file, err := fixtureFile("a.wav")
	if err != nil {
		t.Fatalf("Failed to open fixture file")
	}

	reader := NewReader(file)
	fmt, err := reader.Format()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, AudioFormatPCM, int(fmt.AudioFormat))
	assert.Equal(t, 2, int(fmt.NumChannels))
	assert.Equal(t, 44100, int(fmt.SampleRate))
	assert.Equal(t, 44100*4, int(fmt.ByteRate))
	assert.Equal(t, blockAlign, int(fmt.BlockAlign))
	assert.Equal(t, 16, int(fmt.BitsPerSample))

	duration, err := reader.Duration()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "1.381496598s", duration.String())

	samples, err := reader.ReadSamples(1)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 1, len(samples))

	sample := samples[0]

	assert.Equal(t, 318, reader.IntValue(sample, 0))
	assert.Equal(t, 289, reader.IntValue(sample, 1))
	assert.Assert(t, math.Abs(reader.FloatValue(sample, 0)-0.009705) <= 0.0001)
	assert.Assert(t, math.Abs(reader.FloatValue(sample, 1)-0.008820) <= 0.0001)

	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(bytes), int(reader.WavData.Size)-(1*blockAlign))

	t.Logf("Data size: %d", len(bytes))
}

func TestReadMulaw(t *testing.T) {
	blockAlign := 1

	file, err := fixtureFile("mulaw.wav")
	if err != nil {
		t.Fatalf("Failed to open fixture file")
	}

	reader := NewReader(file)
	fmt, err := reader.Format()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, AudioFormatMULaw, int(fmt.AudioFormat))
	assert.Equal(t, 1, int(fmt.NumChannels))
	assert.Equal(t, 8000, int(fmt.SampleRate))
	assert.Equal(t, 8000, int(fmt.ByteRate))
	assert.Equal(t, blockAlign, int(fmt.BlockAlign))
	assert.Equal(t, 8, int(fmt.BitsPerSample))

	duration, err := reader.Duration()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "4.59125s", duration.String())

	samples, err := reader.ReadSamples(1)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 1, len(samples))

	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(bytes), int(reader.WavData.Size)-(1*blockAlign))

	t.Logf("Data size: %d", len(bytes))
}

func TestReadAlaw(t *testing.T) {
	blockAlign := 1

	file, err := fixtureFile("alaw.wav")
	if err != nil {
		t.Fatalf("Failed to open fixture file")
	}

	reader := NewReader(file)
	fmt, err := reader.Format()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, AudioFormatALaw, int(fmt.AudioFormat))
	assert.Equal(t, 1, int(fmt.NumChannels))
	assert.Equal(t, 8000, int(fmt.SampleRate))
	assert.Equal(t, 8000, int(fmt.ByteRate))
	assert.Equal(t, blockAlign, int(fmt.BlockAlign))
	assert.Equal(t, 8, int(fmt.BitsPerSample))

	duration, err := reader.Duration()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "4.59125s", duration.String())

	samples, err := reader.ReadSamples(1)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 1, len(samples))

	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(bytes), int(reader.WavData.Size)-(1*blockAlign))

	t.Logf("Data size: %d", len(bytes))
}
