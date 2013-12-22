# go-wav [![Build Status](https://travis-ci.org/youpy/go-wav.png?branch=master)](https://travis-ci.org/youpy/go-wav)

A Go library to read WAVE(RIFF waveform Audio) Format

## Usage

```go
package main

import (
	"flag"
	"fmt"
	wav "github.com/youpy/go-wav"
	"io"
	"os"
)

func main() {
	infile_epath := flag.String("infile", "", "wav file to read")
	flag.Parse()

	file, _ := os.Open(*infile_path)
	reader := wav.NewReader(file)
	wav, _ := reader.Read()

	for {
		samples, err := wav.ReadSamples()
		if err == io.EOF {
			break
		}

		for _, sample := range samples {
			fmt.Printf("L/R: %d/%d\n", wav.IntValue(sample, 0), wav.IntValue(sample, 1))
		}
	}
}
```

## Supported format

Number of channels

- 1(mono)
- 2(stereo)

Bits per sample

- 16bit
- 8bit

## Documentation

- https://godoc.org/github.com/youpy/go-wav
