package main

import (
	"flag"
	"fmt"
	wav "github.com/youpy/go-wav"
	"io"
	"os"
	"testing"
)

func main() {
	infile_epath := flag.String("infile", "", "wav file to read")
	flag.Parse()

	n := []uint32{1, 10, 100, 1000, 2000, 3000, 5000, 8000, 10000, 20000, 40000}

	var t int

	for _, numSamples := range n {
		result := testing.Benchmark(func(b *testing.B) {
			file, _ := os.Open(*infile_epath)
			reader := wav.NewReader(file)
			wav, _ := reader.Read()

			for {
				samples, err := wav.ReadSamples(numSamples)
				if err == io.EOF {
					break
				}
				for _, sample := range samples {
					t += wav.IntValue(sample, 0)
					t += wav.IntValue(sample, 1)
				}
			}

		})
		fmt.Printf("ReadSamples(%d): \t%s\n", numSamples, result.String())
	}
}
