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

			for {
				samples, err := reader.ReadSamples(numSamples)
				if err == io.EOF {
					break
				}
				for _, sample := range samples {
					t += reader.IntValue(sample, 0)
					t += reader.IntValue(sample, 1)
				}
			}
		})
		fmt.Printf("ReadSamples(%d): \t%s\n", numSamples, result.String())
	}
}
