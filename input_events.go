package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

type timestamp struct {
	seconds  uint64
	microsec uint64
}

type input_event struct {
	timestamp timestamp
	etype     uint16
	code      uint16
	value     int32
}

func openInputFD(path string) (*os.File, error) {
	input, err := os.OpenFile(path, os.O_RDONLY, 0600)
	if err != nil {
		return nil, err
	}

	return input, nil
}

func processInputEvent(events chan input_event, done chan struct{}, inputFile *os.File) {
	var event input_event
	var buffer = make([]byte, 24)

	for {
		n, err := inputFile.Read(buffer)
		if err != nil {
			return
		}

		if n != 24 {
			log.Println("Wierd Input Event Size: ", n)
		}

		event.timestamp.seconds = binary.LittleEndian.Uint64(buffer[0:8])
		event.timestamp.microsec = binary.LittleEndian.Uint64(buffer[8:16])
		event.etype = binary.LittleEndian.Uint16(buffer[16:18])
		event.code = binary.LittleEndian.Uint16(buffer[18:20])
		event.value = int32(binary.LittleEndian.Uint32(buffer[20:24]))

		select {
		case <-done:
			return
		case events <- event:
		}
	}
}

const (
	EvMake   = 1 // when key is pressed
	EvBreak  = 0 // when key is released
	EvRepeat = 2 // when key switches to repeating
)

func main() {
	events := make(chan input_event, 1)
	done := make(chan struct{})
	input, err := openInputFD("/dev/input/event4")
	if err != nil {
		log.Fatal(err)
	}
	defer input.Close()
	defer close(done)

	go processInputEvent(events, done, input)

	var scanCode, prevCode uint
	var shiftDown, ctrlDown, altDown bool
	var countRepeats int

	var out = os.Stdout

	_ = ctrlDown
	_ = altDown

	for evnt := range events {

		if evnt.etype != 1 { // Keyboard events are always type 1
			continue
		}

		//log.Printf("Event: %v", evnt)
		scanCode = uint(evnt.code)
		if scanCode >= uint(len(charOrFunc)) {
			log.Printf("ScanCode outside of range: ", scanCode)
			continue
		}

		if evnt.value == EvRepeat {
			countRepeats++
		} else if countRepeats > 0 {
			if prevCode == KeyRightShift || prevCode == KeyLeftCtrl || prevCode == KeyRightAlt || prevCode == KeyLeftAlt || prevCode == KeyLeftShift || prevCode == KeyRightCtrl {
			} else {
				fmt.Fprintf(out, "<#+%d>", countRepeats)
			}
			countRepeats = 0
		}

		if evnt.value == EvMake {
			if scanCode == KeyLeftShift || scanCode == KeyRightShift {
				shiftDown = true
			}
			if scanCode == KeyRightAlt {
				altDown = true
			}
			if scanCode == KeyLeftCtrl || scanCode == KeyRightCtrl {
				ctrlDown = true
			}

			var key byte

			if isCharKey(scanCode) {
				if shiftDown == true {
					key = shiftKeys[toCharKeysIndex(int(scanCode))]
					if key == 0 {
						key = charKeys[toCharKeysIndex(int(scanCode))]
					}
				} else {
					key = charKeys[toCharKeysIndex(int(scanCode))]
				}

				if key != 0 {
					fmt.Fprintf(out, "%1c", key)
				}
			} else if isFuncKey(scanCode) {
				if key == KeySpace || key == KeyTab {
					fmt.Fprintf(out, " ")
				} else if key == KeyEnter || key == KeyKPEnter {
					fmt.Fprintf(out, "\n")
				} else {
					fmt.Fprintf(out, "%1s", funcKeys[toFuncKeysIndex(int(scanCode))])
				}
			} else {
				fmt.Fprintf(out, "<E-%x>", scanCode) // unknown scancode
			}
		}

		if evnt.value == EvBreak {
			if scanCode == KeyLeftShift || scanCode == KeyRightShift {
				shiftDown = false
			}
			if scanCode == KeyRightAlt {
				altDown = false
			}
			if scanCode == KeyLeftCtrl || scanCode == KeyRightCtrl {
				ctrlDown = false
			}
		}
	}
}
