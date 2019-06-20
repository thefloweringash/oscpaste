package main

import (
	"encoding/base64"
	"log"
	"os"
	"strings"

	"github.com/mattn/go-isatty"
	"golang.org/x/crypto/ssh/terminal"
)

func readPasteBuffer(term *os.File) ([]byte, error) {
	var err error

	fd := int(term.Fd())

	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		return nil, err
	}
	defer terminal.Restore(fd, oldState)

	_, err = term.Write([]byte("\033]52;;?\007"))
	if err != nil {
		return nil, err
	}

	buffer := []byte{}
	chbuf := make([]byte, 1)
	for {
		_, err := term.Read(chbuf)
		if err != nil {
			return nil, err
		}
		ch := chbuf[0]
		if ch == 007 {
			break
		}
		buffer = append(buffer, ch)
	}

	return buffer, nil
}

func main() {
	if !isatty.IsTerminal(os.Stdin.Fd()) {
		log.Fatalln("stdin is not a terminal!")
	}

	buffer, err := readPasteBuffer(os.Stdin)
	if err != nil {
		panic(err)
	}

	parts := strings.Split(string(buffer), ";")
	encoded := parts[2]

	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		panic(err)
	}

	_, err = os.Stdout.Write(decoded)
	if err != nil {
		panic(err)
	}
}
