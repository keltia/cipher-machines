package machine

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	//      ABCDEFGHIJKLMNOPQRSTUVWXYZ
	rIr1 = "KMFLGDQVZNTOWYHXU/SPAIBRCJE"
	rIr2 = "MFLGDQVZNTOWYHXUS/PAIBRCJEK"
)

func TestNewRotor(t *testing.T) {
	var notch int

	r, err := NewRotor(rI, false)

	assert.NoError(t, err, "should be no error")
	assert.EqualValues(t, r.size, len(rI)-1, "they should be the same")
	assert.EqualValues(t, len(r.rotor), len(rI)-1, "they should be of the same length")
	assert.EqualValues(t, r.refl, false, "should be false")

	// What we want
	arI := make([]int, len(rI)-1)
	var i = 0
	for _, v := range rI {
		if string(v) == "/" {
			notch = i
		} else {
			arI[i] = textToInt[string(v)]
			i++
		}
	}

	// Check content
	assert.EqualValues(t, len(r.rotor), RotorSize, "should be equal")
	assert.EqualValues(t, r.rotor, arI, "should be equal")
	assert.EqualValues(t, r.notch, notch, "should be equal")

	r, err = NewRotor("KJKJK", false)
	assert.Error(t, err, "should be error")
	assert.NotNil(t, r, "R not null")

	r, err = NewRotor("KJKJK", true)
	assert.Error(t, err, "should be error")
	assert.NotNil(t, r, "R not null")
}

func TestNewRotor_Refl(t *testing.T) {
	r, err := NewRotor(RfA, true)

	assert.NoError(t, err, "should be no error")
	assert.NotNil(t, r, "R not null")
	assert.EqualValues(t, r.size, len(RfA), "they should be the same")
	assert.EqualValues(t, len(r.rotor), len(RfA), "they should be of the same length")
	assert.EqualValues(t, r.refl, true, "should be true")
}

func TestRotor_Start(t *testing.T) {
	r1, err := NewRotor(rI, false)
	assert.NoError(t, err, "should be no error")
	assert.EqualValues(t, r1.index, 0, "should be equal")

	err = r1.Start(13)
	assert.EqualValues(t, r1.index, 13, "should be equal")
	assert.Nil(t, err, "no error")

	err = r1.Start(42)
	assert.Error(t, err, "should be an error")
	assert.EqualValues(t, r1.index, 13, "should be equal")
}

func TestRotor_Rotate1(t *testing.T) {
	r, err := NewRotor(rI, false)
	assert.NoError(t, err, "should be no error")

	r1, err := NewRotor(rIr1, false)
	assert.NoError(t, err, "should be no error")

	r.Rotate()
	assert.EqualValues(t, r, r1, "should be equal")
	assert.EqualValues(t, r.rotor, r1.rotor, "should be equal")
}

func TestRotor_Rotate2(t *testing.T) {
	r, err := NewRotor(rI, false)
	assert.NoError(t, err, "should be no error")

	r1, err := NewRotor(rIr1, false)
	assert.NoError(t, err, "should be no error")

	r2, err := NewRotor(rIr2, false)
	assert.NoError(t, err, "should be no error")

	r.Rotate()
	assert.EqualValues(t, r, r1, "should be equal")
	assert.EqualValues(t, r.rotor, r1.rotor, "should be equal")

	r.Rotate()
	assert.EqualValues(t, r.rotor, r2.rotor, "should be equal")

	// Cross-check
	r1.Rotate()
	assert.EqualValues(t, r1.rotor, r2.rotor, "should be equal")
}

func TestRotor_Rotate_Refl(t *testing.T) {
	r, err := NewRotor(RfA, true)
	assert.NoError(t, err, "should be no error")

	r1, err := NewRotor(RfA, true)
	r.Rotate()
	assert.EqualValues(t, r, r1, "should be equal for a reflector")
}

func TestRotor_Start1(t *testing.T) {
	r, err := NewRotor(rI, false)
	assert.NoError(t, err, "should be no error")

	r.Start(21)
	assert.Equal(t, r.index, 21, "should be equal")
}

func TestRotor_Start2(t *testing.T) {
	r, err := NewRotor(rI, false)
	assert.NoError(t, err, "should be no error")

	_ = r.Start(21)
	assert.Equal(t, r.index, 21, "should be equal")
}

func TestRotor_Start_OOB(t *testing.T) {
	r, err := NewRotor(rI, false)
	assert.NoError(t, err, "should be no error")

	err = r.Start(42)
	assert.Error(t, err, "index out of bounds")
}

func TestRotor_Step(t *testing.T) {
	r, err := NewRotor(rI, false)
	assert.NoError(t, err, "should be no error")

	r1, err := NewRotor(rI, false)
	assert.NoError(t, err, "should be no error")

	r.Step()
	r1i := (r1.index + 1) % r1.size
	assert.EqualValues(t, r.index, r1i, "index is +1 mod size")
}

func TestRotor_Step_Wrap(t *testing.T) {
	r, err := NewRotor(rI, false)
	assert.NoError(t, err, "should be no error")


	r.index = r.notch - 1
	r.Step()
	assert.EqualValues(t, true, r.wrap, "index = notch")
	r.Step()
	assert.EqualValues(t, false, r.wrap, "index = notch")
}

func TestRotor_Step_Refl(t *testing.T) {
	r, err := NewRotor(RfA, true)
	assert.NoError(t, err, "should be no error")

	r1, err := NewRotor(RfA, true)
	assert.NoError(t, err, "should be no error")

	r.Step()
	assert.EqualValues(t, r.index, r1.index, "index is invariant for reflectors")
}

func TestRotor_Step_Refl_RT(t *testing.T) {
	r, err := NewRotor(RfB, true)
	assert.NoError(t, err, "should be no error")

	p := textToInt["J"]
	c := intToText[r.Out(p)]
	assert.EqualValues(t, "X", c, "should be equal")

	d := r.Out(textToInt[c])
	assert.EqualValues(t, p, d, "should be equal")
}

func TestRotor_HasWrapped(t *testing.T) {
	r, err := NewRotor(rI, false)
	assert.NoError(t, err, "should be no error")

	r.Start(3)
	r.Step()
	assert.EqualValues(t, false, r.NotchHit(), "should be false")

	r.Start(r.notch - 1)
	r.Step()
	assert.EqualValues(t, true, r.NotchHit(), "should be true")
}

func TestRotor_In(t *testing.T) {
	// EKMFLGDQVZNTOWYHXUSPAIBRCJ
	r, err := NewRotor(rI, false)
	assert.NoError(t, err, "should be no error")

	v := r.In(textToInt["E"])
	assert.EqualValues(t, textToInt["A"], v, "should be equal")

	v = r.In(textToInt["A"])
	assert.EqualValues(t, textToInt["U"], v, "should be equal")

	r.Step()
	v = r.In(textToInt["K"])
	assert.EqualValues(t, textToInt["A"], v, "should be equal")

	v = r.In(textToInt["A"])
	assert.EqualValues(t, textToInt["T"], v, "should be equal")

	v = r.In(666)
	assert.EqualValues(t, -1, v, "should be equal")
}

func TestRotor_InRefl(t *testing.T) {
	r, err := NewRotor(RfB, true)
	assert.NoError(t, err, "should be no error")

	v := r.Out(textToInt["A"])
	p := intToText[r.Out(v)]
	assert.EqualValues(t, "A", p, "should be equal")

	v = r.In(textToInt["A"])
	p = intToText[r.In(v)]
	assert.EqualValues(t, "A", p, "should be equal")
}

func TestRotor_Out(t *testing.T) {
	// EKMFLGDQVZNTOWYHXUSPAIBRCJ
	r, err := NewRotor(rI, false)
	assert.NoError(t, err, "should be no error")

	v := r.Out(textToInt["A"])
	assert.EqualValues(t, v, textToInt["E"], "should be equal")
	r.Step()
	v = r.Out(textToInt["A"])
	assert.EqualValues(t, v, textToInt["K"], "should be equal")
}
