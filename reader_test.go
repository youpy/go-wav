package wav

import (
	"io/ioutil"
	"math"
	"testing"
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

	if fmt.AudioFormat != AudioFormatPCM {
		t.Fatalf("Audio format is invalid: %d", fmt.AudioFormat)
	}

	if fmt.NumChannels != 2 {
		t.Fatalf("Number of channels is invalid: %d", fmt.NumChannels)
	}

	if fmt.SampleRate != 44100 {
		t.Fatalf("Sample rate is invalid: %d", fmt.SampleRate)
	}

	if fmt.ByteRate != 44100*4 {
		t.Fatalf("Byte rate is invalid: %d", fmt.ByteRate)
	}

	if int(fmt.BlockAlign) != blockAlign {
		t.Fatalf("Block align is invalid: %d", fmt.BlockAlign)
	}

	if fmt.BitsPerSample != 16 {
		t.Fatalf("Bits per sample is invalid: %d", fmt.BitsPerSample)
	}

	samples, err := reader.ReadSamples(1)

	if len(samples) != 1 {
		t.Fatalf("Length of samples is invalid: %d", len(samples))
	}

	sample := samples[0]

	if reader.IntValue(sample, 0) != 318 {
		t.Fatalf("Value is invalid: %d", reader.IntValue(sample, 0))
	}

	if reader.IntValue(sample, 1) != 289 {
		t.Fatalf("Value is invalid: %d", reader.IntValue(sample, 1))
	}

	if math.Abs(reader.FloatValue(sample, 0)-0.004852) > 0.0001 {
		t.Fatalf("Value is invalid: %f", reader.FloatValue(sample, 0))
	}

	if math.Abs(reader.FloatValue(sample, 1)-0.004409) > 0.0001 {
		t.Fatalf("Value is invalid: %d", reader.FloatValue(sample, 1))
	}

	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}

	if len(bytes) != int(reader.WavData.Size)-(1*blockAlign) {
		t.Fatalf("Data size is invalid: %d", len(bytes))
	}

	t.Logf("Data size: %d", len(bytes))
}
