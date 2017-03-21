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
		log.Fatalf("Invalid Setup() with %#v-- %v", rotors, err)
	}

	fmt.Println("---- single step")
	err = e.SetRotorSettings([]int{0, 0, 20})
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	// should have [rI/0, rII/0, rIII/20]
	e.DumpState(true)
	e.DumpIndex()
	fmt.Println("----")
	e.Step() // normal step
	e.DumpIndex()
	e.Step() // notch reached
	e.DumpIndex()
	e.Step() // normal step
	e.DumpIndex()

	fmt.Println("----- double step")
	err = e.SetRotorSettings([]int{0, 3, 20})
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	e.DumpIndex()
	fmt.Println("----")
	e.Step() // normal step
	e.DumpIndex()
	e.Step() // right-most step
	e.DumpIndex()
	e.Step() // left-most will step and middle is double step
	e.DumpIndex()
	e.Step() // normal step
	e.DumpIndex()
}

func main() {
	e, _ := machine.NewEnigma(3)
	err := e.Setup(rotors)
	if err != nil {
		log.Fatalf("Invalid Setup() with %#v: %v", rotors, err)
	}

	err = e.AddReflector(RfB)
	if err != nil {
		log.Fatalf("Invalid AddReflector() with %#v", RfB)
	}

	//e.SetPlugboard(PBS)
	e.DumpState(true)

	plain := "AAAAA"
	want := "BDZGO"

	cipher := e.Encrypt(plain)
	fmt.Printf("Plain: %s\nCipher: %s\nWant: %s\n", plain, cipher, want)

	testNewStep()
}
