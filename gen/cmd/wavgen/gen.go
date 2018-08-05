package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/irifrance/snd"
	"github.com/irifrance/snd/encoding/wav"
	"github.com/irifrance/snd/freq"
	"github.com/irifrance/snd/gen"
	"github.com/irifrance/snd/ops"
	"github.com/irifrance/snd/sample"
)

var hz = flag.Int("sin", 0, "sine wave frequency in hertz")
var dur = flag.Duration("dur", time.Second, "duration")
var chirp = flag.Int("chirp", 0, "linear chirp delta in hz/sec")
var chirpStart = flag.Int("chirpstart", 100, "linear chirp start in Hz.")
var note = flag.Int("note", 0, "note fundamental in hertz")

func main() {
	flag.Parse()
	var src snd.Source
	if *hz != 0 {
		src = gen.Sin(freq.T(*hz) * freq.Hertz)
	} else if *chirp != 0 {
		step := gen.Default().SampleRate() / freq.Hertz
		src = gen.Chirp(freq.T(*chirpStart)*freq.Hertz, (freq.T(*chirp)*freq.Hertz)/step)
	} else if *note != 0 {
		src = gen.Note(freq.T(*note) * freq.Hertz)
	}
	src = ops.LimitDur(src, *dur)
	fmt := wav.NewFormatForm(snd.MonoCd(), sample.SFloat32L)
	enc, err := wav.NewEncoder(fmt, os.Stdout)
	defer func() {
		enc.Close()
		os.Stdout.Close()
	}()
	if err != nil {
		log.Fatal(err)
	}
	if err := ops.Copy(enc, src); err != nil {
		log.Fatal(err)
	}
}
