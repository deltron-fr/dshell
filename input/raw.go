package input

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"
)

func RawModeHandler() string {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	r := bufio.NewReader(os.Stdin)
	var buffer []byte
	var cursorPos int

	for {
		b, err := r.ReadByte()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v", err)
			break
		}

		if b == 0x0D {
			fmt.Fprintf(os.Stdout, "\r\n")
			return string(buffer)
		}

		if b <= 0x1f || b == 0x7f {
			switch b {
			case 0x1b:
				key := handleKeys(r)
				switch key {
				case "":
					continue
				case "Left":
					if cursorPos > 0 {
						fmt.Fprintf(os.Stdout, "\x1b[D")
						cursorPos--
					}
				case "Right":
					if cursorPos < len(buffer) {
						fmt.Fprintf(os.Stdout, "\x1b[C")
						cursorPos++
					}
				}
			case 0x0A, 0x0C:
				fmt.Fprintf(os.Stdout, "\r\n")
				return string(buffer)

			case 0x7f, 0x08:
				if len(buffer) == 0 {
					continue
				}
				fmt.Fprintf(os.Stdout, "\x1b[D")
				fmt.Fprintf(os.Stdout, " ")
				fmt.Fprintf(os.Stdout, "\x1b[D")
				cursorPos--
				buffer = buffer[:len(buffer)-1]

			case 0x09:
				if len(buffer) == 0 {
					continue
				}

				parts := strings.Split(string(buffer), " ")
				targetInput := parts[len(parts)-1]
				restOfInput := autoCompletion(targetInput)
				if len(restOfInput) == 0 {
					continue
				}

				for _, b := range restOfInput {
					fmt.Fprintf(os.Stdout, "%c", b)
					buffer = append(buffer, b)
					cursorPos++
				}
				fmt.Fprintf(os.Stdout, " ")
				buffer = append(buffer, ' ')
				cursorPos++
			}
		} else {
			if cursorPos == len(buffer) {
				fmt.Fprintf(os.Stdout, "%c", b)
				buffer = append(buffer, b)
				cursorPos++
			} else {
				buffer = append(buffer, 0)
				copy(buffer[cursorPos+1:], buffer[cursorPos:len(buffer)-1])
				buffer[cursorPos] = b

				for i := cursorPos; i < len(buffer); i++ {
					fmt.Fprintf(os.Stdout, "%c", buffer[i])
				}
				for i := 0; i < len(buffer[cursorPos:])-1; i++ {
					fmt.Fprintf(os.Stdout, "\x1b[D")
				}
			}
		}
	}
	return ""
}
