package main

import (
	"fmt"
	sp "github.com/scipipe/scipipe"
	"github.com/scipipe/scipipe/components"
	"math"
	"runtime"
	"strings"
)

const (
	BUFSIZE = 16384
)

// ======= Main =======

func main() {
	numThreads := runtime.NumCPU() - 1
	fmt.Println("Starting ", numThreads, " threads ...")
	runtime.GOMAXPROCS(numThreads)

	// Init processes
	hisay := NewHiSayer()
	split := NewStringSplitter()
	lower := NewLowerCaser()
	upper := NewUpperCaser()
	zippr := NewZipper()
	s2byt := NewStrToByte()
	filew := components.NewFileWriterFromPath("out2.txt")
	// prntr := NewPrinter(pl)

	// Network definition *** This is where to look! ***
	split.In = hisay.Out
	lower.In = split.OutLeft
	upper.In = split.OutRight
	zippr.In1 = lower.Out
	zippr.In2 = upper.Out
	s2byt.In = zippr.Out
	filew.In = s2byt.Out

	// Create and run pipeline
	pl := sp.NewPipelineRunner()
	pl.AddProcesses(hisay, split, lower, upper, zippr, s2byt, filew)
	pl.Run()
}

// ======= HiSayer =======

type hiSayer struct {
	Out chan string
}

func NewHiSayer() *hiSayer {
	t := &hiSayer{Out: make(chan string, BUFSIZE)}
	return t
}

func (proc *hiSayer) Run() {
	defer close(proc.Out)
	for i := 1; i <= 1e6; i++ {
		proc.Out <- fmt.Sprintf("Hi for the %d:th time!", i)
	}
}

func (proc *hiSayer) IsConnected() bool { return true }

// ======= StringSplitter =======

type stringSplitter struct {
	In       chan string
	OutLeft  chan string
	OutRight chan string
}

func NewStringSplitter() *stringSplitter {
	t := &stringSplitter{
		OutLeft:  make(chan string, BUFSIZE),
		OutRight: make(chan string, BUFSIZE),
	}
	return t
}

func (proc *stringSplitter) Run() {
	defer close(proc.OutLeft)
	defer close(proc.OutRight)
	for s := range proc.In {
		halfLen := int(math.Floor(float64(len(s)) / float64(2)))
		proc.OutLeft <- s[0:halfLen]
		proc.OutRight <- s[halfLen:len(s)]
	}
}

func (proc *stringSplitter) IsConnected() bool { return true }

// ======= LowerCaser =======

type lowerCaser struct {
	In  chan string
	Out chan string
}

func NewLowerCaser() *lowerCaser {
	t := &lowerCaser{Out: make(chan string, BUFSIZE)}
	return t
}

func (proc *lowerCaser) Run() {
	defer close(proc.Out)
	for s := range proc.In {
		proc.Out <- strings.ToLower(s)
	}
}

func (proc *lowerCaser) IsConnected() bool { return true }

// ======= UpperCaser =======

type upperCaser struct {
	In  chan string
	Out chan string
}

func NewUpperCaser() *upperCaser {
	t := &upperCaser{Out: make(chan string, BUFSIZE)}
	return t
}

func (proc *upperCaser) Run() {
	defer close(proc.Out)
	for s := range proc.In {
		proc.Out <- strings.ToUpper(s)
	}
}

func (proc *upperCaser) IsConnected() bool { return true }

// ======= Merger =======

type zipper struct {
	In1 chan string
	In2 chan string
	Out chan string
}

func NewZipper() *zipper {
	t := &zipper{Out: make(chan string, BUFSIZE)}
	return t
}

func (proc *zipper) Run() {
	defer close(proc.Out)
	for {
		s1, ok1 := <-proc.In1
		s2, ok2 := <-proc.In2
		if !ok1 && !ok2 {
			break
		}
		proc.Out <- fmt.Sprint(s1, s2)
	}
}

func (proc *zipper) IsConnected() bool { return true }

// ======= Printer =======

type printer struct {
	In chan string
}

func NewPrinter() *printer {
	t := &printer{}
	return t
}

func (proc *printer) Run() {
	for s := range proc.In {
		fmt.Println(s)
	}
}

func (proc *printer) IsConnected() bool { return true }

// ======= StrToByte =======

type strToByte struct {
	In  chan string
	Out chan []byte
}

func NewStrToByte() *strToByte {
	return &strToByte{
		In:  make(chan string, BUFSIZE),
		Out: make(chan []byte, BUFSIZE),
	}
}

func (proc *strToByte) Run() {
	defer close(proc.Out)
	for line := range proc.In {
		proc.Out <- []byte(line)
	}
}

func (proc *strToByte) IsConnected() bool { return true }
