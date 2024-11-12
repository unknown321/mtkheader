package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"os"
)

type Header struct {
	Magic   [4]byte
	Length  uint32
	Type    [32]byte
	Padding [472]byte
}

var Magic = [4]byte{'\x88', '\x16', '\x88', '\x58'}

func (h *Header) FromReader(r io.Reader) error {
	return binary.Read(r, binary.LittleEndian, h)
}

func (h *Header) Write(w io.Writer) error {
	return binary.Write(w, binary.LittleEndian, h)
}

func NewHeader() *Header {
	h := &Header{
		Magic:   Magic,
		Length:  0,
		Type:    [32]byte{},
		Padding: [472]byte{},
	}

	h.Padding[0] = '\xff'
	for j := 1; j < len(h.Padding); j *= 2 {
		copy(h.Padding[j:], h.Padding[:j])
	}

	return h
}

func Patch(headerFile io.Reader, out io.Writer, s os.FileInfo) error {
	var err error

	h := &Header{}
	if err = h.FromReader(headerFile); err != nil {
		return err
	}

	if s.IsDir() {
		return fmt.Errorf("%s is a directory", s.Name())
	}

	h.Length = uint32(s.Size())

	if err = h.Write(out); err != nil {
		return fmt.Errorf("cannot write header: %w", err)
	}

	return nil
}

func Info(filename string) error {
	var err error
	var hf *os.File

	hf, err = os.Open(filename)
	if err != nil {
		return fmt.Errorf("cannot open header file: %w", err)
	}

	h := &Header{}
	if err = h.FromReader(hf); err != nil {
		return err
	}

	hType := bytes.Trim(h.Type[:], "\x00")
	fmt.Printf("File:   %s\nType:   %s\nLength: %d\n", filename, hType, h.Length)

	return nil
}

func main() {
	complete()
	var header string
	var content string
	flag.StringVar(&header, "header", "", "mtk header")
	flag.StringVar(&content, "content", "", "content file")
	sysOut := os.Stdout
	flag.CommandLine.SetOutput(sysOut)
	flag.Usage = func() {
		fmt.Fprintf(sysOut, "Usage: %s -header [FILE] -content [FILE]\n", os.Args[0])
		fmt.Fprintf(sysOut, "   or: %s [FILE]\n", os.Args[0])
		fmt.Fprintf(sysOut, "\nPatches MediaTek `header` file with size of `content`.\n")
		fmt.Fprintf(sysOut, "Print header info by providing path to header file without any flags\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	var err error

	if header == "" {
		if len(os.Args) < 2 {
			fmt.Println("no header file provided")
			os.Exit(1)
		}

		if os.Args[1] == "" {
			fmt.Println("no header file provided")
			os.Exit(1)
		}

		header = os.Args[1]
	}

	if content == "" {
		if err = Info(header); err != nil {
			fmt.Printf("cannot get info: %s\n", err.Error())
			os.Exit(1)
		}

		os.Exit(0)
	}

	data, err := os.ReadFile(header)
	if err != nil {
		fmt.Printf("cannot open header file: %s\n", err.Error())
		os.Exit(1)
	}

	dataReader := bytes.NewReader(data)

	s, err := os.Stat(content)
	if err != nil {
		fmt.Printf("cannot stat content file: %s\n", err.Error())
	}

	out, err := os.OpenFile(header, os.O_RDWR, 0)
	if err != nil {
		fmt.Printf("cannot open file for writing: %s\n", err.Error())
		os.Exit(1)
	}

	if err = Patch(dataReader, out, s); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
