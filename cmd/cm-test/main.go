package main

import (
    "github.com/keltia/cipher-machines"
    "fmt"
    "log"
)

const (
	//      ABCDEFGHIJKLMNOPQRSTUVWXYZ
    rI   = "EKMFLGDQVZNTOWYHX/USPAIBRCJ"
    rII  = "AJDKS/IRUXBLHWTMCQGZNPYFVOE"
    rIII = "BDFHJLCPRTXVZNYEIWGAKM/USQO"
    rIV  = "ESOVPZJAYQ/UIRHXLNFTGKDCMWB"
    rV   = "VZBRGITYUPSDNHLXAWMJQOFECK/"
    rVI  = "JPGVOUMFYQBENHZRDKASXLICTW/"    // special, has notch in M too

    RfA   = "EJMZALYXVBWFCRQUONTSPIKHGD"
    RfB   = "YRUHQSLDPXNGOKMIEBFZCWVJAT"
    RfC   = "FVPJIAOYEDRZXWGCTKUQSBNMHL"

    PBS   = "EJOYIVAQKWFXMTPSLUBD"
)

var rotors = []string{
    rI,
	rII,
	rIII,
}

func testNewStep() {
    e, _ := machine.NewEnigma(3)
    err := e.Setup(rotors)
    if err != nil {
        log.Fatalf("Invalid Setup() with %#v", rotors)
    }

	fmt.Println("---- single step")
    e.SetRotorSettings([]int{0, 0, 20})
	// should have [rI/0, rII/0, rIII/20]
    e.DumpState(false)
    e.DumpIndex()
	fmt.Println("----")
    e.NewStep()
    e.DumpIndex()
    e.NewStep()
    e.DumpIndex()
    e.NewStep()
    e.DumpIndex()

	fmt.Println("----- double step")
	e.SetRotorSettings([]int{0, 3, 20})
	e.DumpIndex()
	fmt.Println("----")
	e.NewStep()
	e.DumpIndex()
	e.NewStep()
	e.DumpIndex()
	e.NewStep()
	e.DumpIndex()
	e.NewStep()
	e.DumpIndex()
}

func main() {
    e, _ := machine.NewEnigma(3)
    err := e.Setup(rotors)
    if err != nil {
        log.Fatalf("Invalid Setup() with %#v", rotors)
    }

    err = e.AddReflector(RfB)
    if err != nil {
        log.Fatalf("Invalid AddReflector() with %#v", RfB)
    }

    e.SetPlugboard(PBS)
    e.DumpState(true)

    plain := "AAAAA"
    want := "BDZGO"

    cipher := e.Encrypt(plain)
    fmt.Printf("Plain: %s\nCipher: %s\nWant: %s\n", plain, cipher, want)

    testNewStep()
}
