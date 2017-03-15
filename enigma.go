package machine

import (
    "bytes"
    "strings"
    "log"
    "fmt"
    "text/scanner"
)

const (
    // Official rotors
    // cf. https://en.wikipedia.org/wiki/Enigma_rotor_details

    rI   = "EKMFLGDQVZNTOWYHXUSPAIBRCJ"
    rII  = "AJDKSIRUXBLHWTMCQGZNPYFVOE"
    rIII = "BDFHJLCPRTXVZNYEIWGAKMUSQO"
    rIV  = "ESOVPZJAYQUIRHXLNFTGKDCMWB"
    rV   = "VZBRGITYUPSDNHLXAWMJQOFECK"
    rVI  = "JPGVOUMFYQBENHZRDKASXLICTW"

    // Reflectors
    RfA   = "EJMZALYXVBWFCRQUONTSPIKHGD"
    RfB   = "YRUHQSLDPXNGOKMIEBFZCWVJAT"
    RfC   = "FVPJIAOYEDRZXWGCTKUQSBNMHL"

    // Plugboard settings
    PBS   = "EJOYIVAQKWFXMTPSLUBD"

    // Enigma types
    EnigmaStd    = 3
    EnigmaMarine = 4
)

// -- Enigma
type Enigma struct {
    PlugBoard map[int]int   // plugboard settings
    RotorSet  []*Rotor // rotor set (3-4)
    Reflector *Rotor   // Enigma specific
    Size      int     // number of rotors
}

func (m *Enigma) Step() (*Enigma){
    return m
}

func (m *Enigma) Setup(rotors []string) (*Enigma){
    // Only plain rotors, no reflector here
    if len(rotors) != m.Size{
		log.Fatalf("Bad size: %d", len(rotors))
    }

    m.RotorSet = make([]*Rotor, m.Size)

    for i, r := range rotors {
	    m.RotorSet[i] = NewRotor(r, false)
		//log.Printf("%v\n", m.RotorSet[i])
    }
	return m
}

func (m *Enigma) AddReflector(ref string) (*Enigma) {
    m.Reflector = NewRotor(ref, true)
    return m
}

func (m *Enigma) SetPlugboard(plug string) (*Enigma) {
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
    return m
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

func (m *Enigma) DumpState() {
    fmt.Printf("%#v\n-----\n", m)
}