package main

import (
	"fmt"
	"github.com/keltia/cipher-machines"
	"log"
)

const (
	//      ABCDEFGHIJKLMNOPQRSTUVWXYZ
	rI   = "EKMFLGDQVZNTOWYH/XUSPAIBRCJ"
	rII  = "AJDK/SIRUXBLHWTMCQGZNPYFVOE"
	rIII = "BDFHJLCPRTXVZNYEIWGAK/MUSQO"
	rIV  = "ESOVPZJAY/QUIRHXLNFTGKDCMWB"
	rV   = "VZBRGITYUPSDNHLXAWMJQOFECK/"
	rVI  = "JPGVOUMFYQBENHZRDKASXLICTW/" // special, has notch in M too

	RfA = "EJMZALYXVBWFCRQUONTSPIKHGD"
	RfB = "YRUHQSLDPXNGOKMIEBFZCWVJAT"
	RfC = "FVPJIAOYEDRZXWGCTKUQSBNMHL"

	PBS = "EJOYIVAQKWFXMTPSLUBD"
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
	e.NewStep()		// normal step
	e.DumpIndex()
	e.NewStep()		// notch reached
	e.DumpIndex()
	e.NewStep()		// normal step
	e.DumpIndex()

	fmt.Println("----- double step")
	e.SetRotorSettings([]int{0, 3, 20})
	e.DumpIndex()
	fmt.Println("----")
	e.NewStep()		// normal step
	e.DumpIndex()
	e.NewStep()		// right-most step
	e.DumpIndex()
	e.NewStep()		// left-most will step and middle is double step
	e.DumpIndex()
	e.NewStep()		// normal step
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
