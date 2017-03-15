package main

import (
    "github.com/keltia/cipher-machines"
    "fmt"
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
    e := machine.NewEnigma(3).Setup(rotors).AddReflector(RfB)
    e.DumpState()

    e.SetPlugboard(PBS)
    e.DumpState()

    e.Step()
    e.DumpState()

    plain := "AAAAA"
    want := "BDZGO"

    cipher := e.Encrypt(plain)
    fmt.Printf("Plain: %s\nCipher: %s\nWant: %s\n", plain, cipher, want)
}
