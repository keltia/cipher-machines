package machine

import (
	"bytes"
	"fmt"
	"strings"
	"text/scanner"
)

const (
	// Official rotors
	// cf. https://en.wikipedia.org/wiki/Enigma_rotor_details
	//      ABCDEFGHIJKLMNOPQRSTUVWXYZ
	rI   = "EKMFLGDQVZNTOWYH/XUSPAIBRCJ"
	rII  = "AJDK/SIRUXBLHWTMCQGZNPYFVOE"
	rIII = "BDFHJLCPRTXVZNYEIWGAK/MUSQO"
	rIV  = "ESOVPZJAYQ/UIRHXLNFTGKDCMWB"
	rV   = "VZBRGITYUPSDNHLXAWMJQOFEC/K"
	rVI  = "JPGVOUMFYQBENHZRDKASXLICT/W" // special, has notch in M too

	// Reflectors
	//     ABCDEFGHIJKLMNOPQRSTUVWXYZ
	RfA = "EJMZALYXVBWFCRQUONTSPIKHGD"
	RfB = "YRUHQSLDPXNGOKMIEBFZCWVJAT"
	RfC = "FVPJIAOYEDRZXWGCTKUQSBNMHL"

	// Plugboard settings
	PBS = "EJOYIVAQKWFXMTPSLUBD"

	// Enigma types
	EnigmaStd    = 3
	EnigmaMarine = 4

	// Misc
	RotorSize = 26
)

// Enigma
type Enigma struct {
	PlugBoard map[int]int // plugboard settings
	RS        *RotorSet   // rotor set (3-4)
	Reflector *Rotor      // Enigma specific
	Size      int         // number of rotors
}

// RS is an array starting at 0 BUT the rotors are numbered the reverse way:
// ([r3]) [r2] [r1] [r0]

func NewEnigma(size int) (m *Enigma, err error) {
	if size != EnigmaStd && size != EnigmaMarine {
		err = fmt.Errorf("wrong size %d, should 3 or 4", size)
		return
	}

	m = &Enigma{
		Size: size,
	}
	return
}

func (m *Enigma) Setup(rotors []string) (err error) {
	// Only plain rotors, no reflector here
	if len(rotors) != m.Size {
		return fmt.Errorf("Bad size: %d", len(rotors))
	}

	m.RS, err = NewRotorSet(rotors)
	return
}

func (m *Enigma) SetRotorSettings(set []int) (err error) {
	if len(set) != m.RS.Len {
		err = fmt.Errorf("Mismatch in rotors number: %d vs %d", len(set), m.RS.Len)
	}
	err = m.RS.Set(set)
	return
}

func (m *Enigma) Settings() (set []int) {
	return m.RS.Settings()
}

func (m *Enigma) AddReflector(ref string) (err error) {
	m.Reflector, err = NewRotor(ref, true)
	return
}

func (m *Enigma) SetPlugboard(plug string) error {
	var s scanner.Scanner

	s.Init(strings.NewReader(plug))

	sa := make(map[int]int)

	var tok rune
	for tok != scanner.EOF {
		tok = s.Next()
		pa := textToInt[string(tok)]

		// Check next one
		nb := s.Next()
		if nb == scanner.EOF {
			break
		}

		pb := string(nb)
		sa[pa] = textToInt[pb]
	}

	m.PlugBoard = sa
	return nil
}

// Step makes the rotors turn.  At some point, the 2nd can step as well, which can trigger
// the 3rd one.  In the Kriegsmarine Enigma, the 4th rotor did not step.
// XXX assume single notch rotors
func (m *Enigma) Step() {
	// New mode, take the notches into account
	m.RS.Step()
}

func (m *Enigma) Out(i int) int {
	var next int

	fmt.Printf("in: %s ", intToText[i])
	// Go through plugboard if any
	if m.PlugBoard != nil {
		if pbc, ok := m.PlugBoard[i]; ok {
			next = pbc
		}
	} else {
		// 1st phase
		next = i
	}
	fmt.Printf("after PB(%s) ", intToText[next])

	// 1st phase back through the rotors
	next = m.RS.left(next)

	fmt.Printf("after rtr(%s) ", intToText[next])

	// Reflector
	next = m.Reflector.Out(next)

	fmt.Printf("after refl(%s) ", intToText[next])

	next = m.RS.right(next)

	fmt.Printf("after rtr.back(%s) ", intToText[next])

	// Finally go through plugboard again if any
	if m.PlugBoard != nil {
		if pbc, ok := m.PlugBoard[next]; ok {
			next = pbc
		}
	}
	fmt.Printf("after PBout(%s) ", intToText[next])

	return next
}

func (m *Enigma) Index() (state []int) {
	state = make([]int, m.Size)
	return
}

func (m *Enigma) Encrypt(text string) (cipher string) {
	var (
		s   scanner.Scanner
		str bytes.Buffer
	)

	s.Init(strings.NewReader(text))

	var tok rune

	for tok != scanner.EOF {
		tok = s.Next()
		p := textToInt[string(tok)]
		fmt.Printf("plain: %s", intToText[p])

		// Stepping is done before current goes through
		m.Step()

		// Dive into the rotors
		c := intToText[m.Out(p)]

		fmt.Printf(" - cipher: %s\n", c)
		str.WriteString(c)
	}
	cipher = str.String()
	return
}

func (m *Enigma) Decrypt(text string) (clear string) {
	return
}

func (m *Enigma) DumpState(t bool) {
	if t {
		fmt.Printf("PB: %#v\nRefl: %#v\n", m.PlugBoard, m.Reflector)
	}
	if m.RS != nil {
		fmt.Printf("rs: %#v\n%#v\n%#v\n-----\n", m.RS.R[0], m.RS.R[1], m.RS.R[2])
	}
}

func (m *Enigma) DumpIndex() {
	fmt.Printf("%s\n", m.RS.Index())
}
