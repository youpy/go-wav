package wav

import (
	"io"
	"io/ioutil"
	"testing"
)

func TestRead(t *testing.T) {
	blockAlign := 4

	file, err := fixtureFile("a.wav")

	if err != nil {
		t.Fatalf("Failed to open fixture file")
	}

	reader := NewReader(file)

	wav, err := reader.Read()
	if err != nil {
		t.Fatal(err)
	}

	fmt := wav.Format

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

	var bytes []byte

	sample, err := wav.ReadSample()

	if sample.IntValue(0) != 318 {
		t.Fatalf("Value is invalid: %d", sample.IntValue(0))
	}

	if sample.IntValue(1) != 289 {
		t.Fatalf("Value is invalid: %d", sample.IntValue(1))
	}

	bytes, err = ioutil.ReadAll(wav)
	if err != nil {
		t.Fatal(err)
	}

	if len(bytes) != int(wav.WavData.Size)-(1*blockAlign) {
		t.Fatalf("Data size is invalid: %d", len(bytes))
	}

	t.Logf("Data size: %d", len(bytes))
}

func BenchmarkRead(b *testing.B) {
	file, err := fixtureFile("a.wav")

	if err != nil {
		b.Fatalf("Failed to open fixture file")
	}

	reader := NewReader(file)

	wav, err := reader.Read()
	if err != nil {
		b.Fatal(err)
	}

	for {
		_, err := wav.ReadSample()
		if err == io.EOF {
			break
		}
	}
}
