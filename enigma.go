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
	RS        *RotorSet    // rotor set (3-4)
	Reflector *Rotor      // Enigma specific
	Size      int         // number of rotors
}

// RS is an array starting at 0 BUT the rotors are numbered the reverse way:
// ([r3]) [r2] [r1] [r0]

// Step makes the rotors turn.  At some point, the 2nd can step as well, which can trigger
// the 3rd one.  In the Kriegsmarine Enigma, the 4th rotor did not step.
// XXX assume single notch rotors
func (m *Enigma) Step() *Enigma {
	// New mode, take the notches into account

	// if this is a 4-wheel machine, the foremost one (aka the 4th) does not move.
	r0 := m.RotorSet[2]
	r1 := m.RotorSet[1]
	r2 := m.RotorSet[0]

	// [r2, r1, r0]
    n0 := r0.NotchHit()
	r0.Step()
	if n0 {
		// Check for double step
        r1.Step()
        n1 := r1.NotchHit()
		if n1 {
            r1.Step()
			r2.Step()
		}
	}
	return m
}

func (m *Enigma) SetRotorSettings(set []int) (err error) {
	if len(set) != len(m.RotorSet) {
		err = fmt.Errorf("Mismatch in rotors number: %d vs %d", len(set), len(m.RotorSet))
	}
	for i, v := range set {
		m.RotorSet[i].Start(v)
	}
	err = nil
	return
}

func (m *Enigma) Setup(rotors []string) (err error) {
	// Only plain rotors, no reflector here
	if len(rotors) != m.Size {
		return fmt.Errorf("Bad size: %d", len(rotors))
	}

	m.RotorSet = make([]*Rotor, m.Size)

	// Reverse insert rotors
	for i, r := range rotors {
		if len(r) != RotorSize+1 {
			return fmt.Errorf("bad length %d should be 26", len(r))
		}
		m.RotorSet[i], err = NewRotor(r, false)
		//log.Printf("%v\n", m.RotorSet[i])
	}
	return
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
    fmt.Printf("in: after PB %s ", intToText[next])

    // 1st phase back through the rotors
    for i := len(m.RotorSet); i <= 0; i-- {
        r := m.RotorSet[i]
        next = r.In(next)
    }

    fmt.Printf("in: after rtr %s ", intToText[next])

	// Reflector
	next = m.Reflector.Out(next)

    fmt.Printf("in: after refl %s ", intToText[next])

    // Go round
    for _, r := range m.RotorSet {
        next = r.Out(next)
    }


    fmt.Printf("in: after rtr.back %s ", intToText[next])

    // Finally go through plugboard again if any
	if m.PlugBoard != nil {
		if pbc, ok := m.PlugBoard[next]; ok {
			next = pbc
		}
	}
    fmt.Printf("in: after PBout %s ", intToText[next])

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

func (m *Enigma) DumpState(t bool) {
	if t {
		fmt.Printf("PB: %#v\nRefl: %#v\n", m.PlugBoard, m.Reflector)
	}
	if m.RotorSet != nil {
		fmt.Printf("r0: %#v\nr1: %#v\nr2: %##v\n-----\n", m.RotorSet[0], m.RotorSet[1], m.RotorSet[2])
	}
}

func (m *Enigma) DumpIndex() {
	var (
		off int
		ri3 string
	)

	if m.Size == EnigmaMarine {
		ri3 = intToText[m.RotorSet[0].index]
		off = 1
	} else {
		ri3 = "-"
		off = 0
	}
	ri2 := intToText[m.RotorSet[off].index]
	ri1 := intToText[m.RotorSet[off+1].index]
	ri0 := intToText[m.RotorSet[off+2].index]
	fmt.Printf("%s%s%s%s\n", ri3, ri2, ri1, ri0)
}
