package main

import (
	"fmt"
	"os"
)

type chanWriter struct {
	ch chan byte
}

func newChanWriter() *chanWriter {
	return &chanWriter{make(chan byte, 1024)}
}

func (w *chanWriter) Chan() <-chan byte {
	return w.ch
}

func (w *chanWriter) Write(p []byte) (int, error) {
	n := 0
	for _, b := range p {
		w.ch <- b
		n++
	}
	return n, nil
}

func (w *chanWriter) Close() error {
	close(w.ch)
	return nil
}

func producer(w *chanWriter) {

	defer w.Close()
	i := 0
	for {
		w.Write([]byte(fmt.Sprint(i)))
		w.Write([]byte(" - Stream "))
		w.Write([]byte("me "))
		w.Write([]byte("PLEAAASE!\n"))
		w.Write([]byte("PLEAAASE!\n"))
		//time.Sleep(2 * time.Second)
		i++
	}
}

func main() {

	// Creates file
	file, err := os.Create("./streaming")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	writer := newChanWriter()
	go producer(writer)

	for c := range writer.Chan() {
		val := []byte{c}
		n, err := file.Write(val)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if n != len(val) {
			fmt.Println("failed to write data")
			os.Exit(1)
		}
	}
	fmt.Println("file write done")

	fmt.Println()
}
