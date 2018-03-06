package main

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/gen2brain/aac-go"
	"github.com/youpy/go-wav"
)

func main() {
	file, err := os.Open("test.wav")
	if err != nil {
		panic(err)
	}

	wreader := wav.NewReader(file)
	f, err := wreader.Format()
	if err != nil {
		panic(err)
	}

	buf := bytes.NewBuffer(make([]byte, 0))

	opts := &aac.Options{}
	opts.SampleRate = int(f.SampleRate)
	opts.NumChannels = int(f.NumChannels)

	enc, err := aac.NewEncoder(buf, opts)
	if err != nil {
		panic(err)
	}

	err = enc.Encode(wreader)
	if err != nil {
		panic(err)
	}

	err = enc.Close()
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("test.aac", buf.Bytes(), 0644)
	if err != nil {
		panic(err)
	}
}
