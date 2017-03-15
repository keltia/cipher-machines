package machine

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

const (
	rIr1 = "KMFLGDQVZNTOWYHXUSPAIBRCJE"
	rIr2 = "MFLGDQVZNTOWYHXUSPAIBRCJEK"
)

func TestNewRotor(t *testing.T) {
	r := NewRotor(rI, false)
	
	assert.EqualValues(t, r.size, len(rI), "they should be the same")
	assert.EqualValues(t, len(r.rotor),len(rI), "they should be of the same length")
	assert.EqualValues(t, r.refl, false, "should be false")

	// What we want
	arI := make([]int, len(rI))
	for i, v := range rI {
		arI[i] = textToInt[string(v)]
	}

	// Check content
	assert.EqualValues(t, r.rotor, arI, "should be equal")
}

func TestNewRotor_Refl(t *testing.T) {
	r := NewRotor(RfA, true)

	assert.EqualValues(t, r.size, len(rI), "they should be the same")
	assert.EqualValues(t, len(r.rotor),len(rI), "they should be of the same length")
	assert.EqualValues(t, r.refl, true, "should be true")
}

func TestRotor_Start(t *testing.T) {
	r1 := NewRotor(rI, false)
	assert.EqualValues(t, r1.index, 0, "should be equal")

	err := r1.Start(13)
	assert.EqualValues(t, r1.index, 13, "should be equal")
	assert.Nil(t, err, "no error")

	err = r1.Start(42)
	assert.Error(t, err, "should be an error")
	assert.EqualValues(t, r1.index, 13, "should be equal")
}

func TestRotor_Rotate1(t *testing.T) {
	r := NewRotor(rI, false)
	r1 := NewRotor(rIr1, false)

	r.Rotate()
	assert.EqualValues(t, r, r1, "should be equal")

	for i, v := range r1.rotor {
		if r.rotor[i] != v {
			t.Errorf("Invalid value at %d\nrI: %v\nr1: %v", i, r, r1)
		}
	}
}

func TestRotor_Rotate2(t *testing.T) {
	r := NewRotor(rI, false)
	r1 := NewRotor(rIr1, false)
	r2 := NewRotor(rIr2, false)

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
	r := NewRotor(RfA, true)
	r1 := NewRotor(RfA, true)
	r.Rotate()
	assert.EqualValues(t, r, r1, "should be equal for a reflector")
}

func TestRotor_Start1(t *testing.T) {
	r := NewRotor(rI, false)
	r.Start(21)
	assert.Equal(t, r.index, 21, "should be equal")
}

func TestRotor_Start2(t *testing.T) {
	r := NewRotor(rI, false)
	_ = r.Start(21)
	assert.Equal(t, r.index, 21, "should be equal")
}

func TestRotor_Start_OOB(t *testing.T) {
	r := NewRotor(rI, false)
	err := r.Start(42)
	assert.Error(t, err,"index out of bounds")
}

func TestRotor_Step(t *testing.T) {
	r := NewRotor(RfA, false)
	r1 := NewRotor(RfA, false)

	r.Step()
	r1i := (r1.index + 1) % r1.size
	assert.EqualValues(t, r.index, r1i, "index is +1 mod size")
}

func TestRotor_Step_Refl(t *testing.T) {
	r := NewRotor(RfA, true)
	r1 := NewRotor(RfA, true)

	r.Step()
	assert.EqualValues(t, r.index, r1.index, "index is invariant for reflectors")
}

func TestRotor_Out(t *testing.T) {

}

func TestRotor_HasWrapped(t *testing.T) {
	r := NewRotor(rI, false)

	r.index = 3
	r.Step()
	assert.EqualValues(t, r.HasWrapped(), false, "should be false")

	r.index = 25
	r.Step()
	assert.EqualValues(t, r.HasWrapped(), true, "should be true")
}
