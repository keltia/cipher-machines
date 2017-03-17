package main

import (
    "github.com/keltia/cipher-machines"
    "fmt"
    "log"
)

const (
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

func main() {
    var rotors = []string{
	    rI,
	    rII,
	    rIII,
    }
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
}
