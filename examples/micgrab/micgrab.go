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
	numChannels := 1
	sampleRate := 48000

	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, func(message string) {
		fmt.Printf("LOG <%v>\n", message)
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer func() {
		_ = ctx.Uninit()
		ctx.Free()
	}()

	deviceConfig := malgo.DefaultDeviceConfig(malgo.Capture)
	deviceConfig.Capture.Format = malgo.FormatS16
	deviceConfig.Capture.Channels = uint32(numChannels)

	var capturedSampleCount uint32
	pCapturedSamples := make([]byte, 0)

	sizeInBytes := uint32(malgo.SampleSizeInBytes(deviceConfig.Capture.Format))
	onRecvFrames := func(pSample2, pSample []byte, framecount uint32) {
		sampleCount := framecount * deviceConfig.Capture.Channels * sizeInBytes
		newCapturedSampleCount := capturedSampleCount + sampleCount
		pCapturedSamples = append(pCapturedSamples, pSample...)
		capturedSampleCount = newCapturedSampleCount
	}

	fmt.Println("Recording...")
	captureCallbacks := malgo.DeviceCallbacks{
		Data: onRecvFrames,
	}
	device, err := malgo.InitDevice(ctx.Context, deviceConfig, captureCallbacks)
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

	os.Exit(0)
}
