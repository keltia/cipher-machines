package main

import (
    "github.com/keltia/cipher-machines"
    "fmt"
    "log"
)

const (
    rI   = "EKMFLGDQVZNTOWYHXUSPAIBRCJ"
    rII  = "AJDKSIRUXBLHWTMCQGZNPYFVOE"
    rIII = "BDFHJLCPRTXVZNYEIWGAKMUSQO"
    rIV  = "ESOVPZJAYQUIRHXLNFTGKDCMWB"
    rV   = "VZBRGITYUPSDNHLXAWMJQOFECK"
    rVI  = "JPGVOUMFYQBENHZRDKASXLICTW"

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
