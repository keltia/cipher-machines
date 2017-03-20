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
	rI   = "EKMFLGDQVZNTOWYHX/USPAIBRCJ"
	rII  = "AJDKS/IRUXBLHWTMCQGZNPYFVOE"
	rIII = "BDFHJLCPRTXVZNYEIWGAKM/USQO"
	rIV  = "ESOVPZJAYQ/UIRHXLNFTGKDCMWB"
	rV   = "VZBRGITYUPSDNHLXAWMJQOFECK/"
	rVI  = "JPGVOUMFYQBENHZRDKASXLICTW/" // special, has notch in M too

	// Reflectors
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
	RotorSet  []*Rotor    // rotor set (3-4)
	Reflector *Rotor      // Enigma specific
	Size      int         // number of rotors
}

// RotorSet is an array starting at 0 BUT the rotors are numbered the reverse way:
// ([r3]) [r2] [r1] [r0]

// Step makes the rotors turn.  At some point, the 2nd can step as well, which can trigger
// the 3rd one.  In the Kriegsmarine Enigma, the 4th rotor did not step.
func (m *Enigma) Step() *Enigma {
	// Naive mode: do not have step at a specific place

	var last = EnigmaStd - 1

	// Loop unrolled, not worth abstracting the whole process
	if len(m.RotorSet) == EnigmaMarine {
		last = EnigmaMarine - 1
	}

	r0 := m.RotorSet[last]
	r0.Step()
	if r0.NotchHit() {
		r1 := m.RotorSet[last-1]
		r1.Step()
		if r1.NotchHit() {
			r2 := m.RotorSet[last-2]
			r2.Step()
		}
	}
	return m
}

// Step makes the rotors turn.  At some point, the 2nd can step as well, which can trigger
// the 3rd one.  In the Kriegsmarine Enigma, the 4th rotor did not step.
// XXX assume single notch rotors
func (m *Enigma) NewStep() *Enigma {
	// New mode, take the notches into account

	// if this is a 4-wheel machine, the foremost one (aka the 4th) does not move.
	r0 := m.RotorSet[2]
	r1 := m.RotorSet[1]
	r2 := m.RotorSet[0]

	// [r2, r1, r0]
	r0.Step()
	if r0.NotchHit() {
		r1.Step()
		// Check for double step
		if r1.NotchHit() {
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

	// Go through plugboard if any
	if m.PlugBoard != nil {
		if pbc, ok := m.PlugBoard[i]; ok {
			next = pbc
		}
	} else {
		// 1st phase
		next = i
	}

	// Go round
	for _, r := range m.RotorSet {
		next = r.Out(next)
	}

	// Reflector
	next = m.Reflector.Out(next)

	// 2nd phase back through the rotors
	for i := len(m.RotorSet); i <= 0; i-- {
		r := m.RotorSet[i]
		next = r.In(next)
	}

	// Finally go through plugboard again if any
	if m.PlugBoard != nil {
		if pbc, ok := m.PlugBoard[next]; ok {
			next = pbc
		}
	}

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
		fmt.Printf("plain: %d", p)

		// Dive into the rotors
		c := intToText[m.Out(p)]

		fmt.Printf(" - cipher: %s\n", c)
		str.WriteString(c)
		m.Step()
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
