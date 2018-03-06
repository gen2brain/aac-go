package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gen2brain/aac-go"
	"github.com/gen2brain/malgo"
)

func main() {
	device := mal.NewDevice()

	numChannels := 2
	sampleRate := 48000

	var capturedSampleCount uint32
	pCapturedSamples := make([]byte, 0)

	onRecvFrames := func(framecount uint32, pSamples []byte) {
		sizeInBytes := device.SampleSizeInBytes(device.Format())
		sampleCount := framecount * device.Channels() * sizeInBytes

		newCapturedSampleCount := capturedSampleCount + sampleCount
		pCapturedSamples = append(pCapturedSamples, pSamples...)
		capturedSampleCount = newCapturedSampleCount
	}

	err := device.ContextInit(nil, mal.ContextConfig{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer device.ContextUninit()

	config := device.ConfigInit(mal.FormatS16, uint32(numChannels), uint32(sampleRate), onRecvFrames, nil)

	fmt.Println("Recording...")
	err = device.Init(mal.Capture, nil, &config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = device.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Press Enter to stop recording...")
	fmt.Scanln()

	device.Uninit()

	fmt.Println("Encoding...")
	buf := bytes.NewBuffer(make([]byte, 0))

	opts := &aac.Options{}
	opts.SampleRate = sampleRate
	opts.NumChannels = numChannels

	enc, err := aac.NewEncoder(buf, opts)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	reader := bytes.NewReader(pCapturedSamples)

	err = enc.Encode(reader)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = enc.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = ioutil.WriteFile("capture.aac", buf.Bytes(), 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Press Enter to quit...")
	fmt.Scanln()

	os.Exit(0)
}
