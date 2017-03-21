package machine

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewEnigma(t *testing.T) {
	e, err := NewEnigma(3)
	assert.EqualValues(t, e.Size, 3, "should be equal")
	assert.Nil(t, err, "no error")

	e, err = NewEnigma(4)
	assert.Nil(t, err, "no error")
	assert.EqualValues(t, 4, e.Size, "should be equal")

	e, err = NewEnigma(666)
	assert.Error(t, err, "should protest")
	assert.Panics(t, func() {
		var _ = e.Size
	}, "should panic")
}

func TestEnigma_Setup(t *testing.T) {
	var rotors = []string{
		rIII,
		rII,
		rI,
	}

	e, _ := NewEnigma(3)
	err := e.Setup(rotors)
	if err != nil {
		t.Errorf("error running .Setup(), err=%v", err)
	}

	rrI, err := NewRotor(rI, false)
	rrII, _ := NewRotor(rII, false)
	rrIII, _ := NewRotor(rIII, false)

	assert.NoError(t, err, "should be ok")
	assert.EqualValues(t, rrI, e.RS.R[2], "should be equal")
	assert.EqualValues(t, rrII, e.RS.R[1], "should be equal")
	assert.EqualValues(t, rrIII, e.RS.R[0], "should be equal")

	assert.Nil(t, e.PlugBoard, "should be nil")
	assert.Nil(t, e.Reflector, "should be nil")
}

func TestEnigma_Setup_Badlen(t *testing.T) {
	var rotors = []string{
		"JHJHSJDHJSHDKHDHKSHDKJSHDKJSHDKJHSKJDH",
		rII,
		rI,
	}

	e, _ := NewEnigma(3)
	err := e.Setup(rotors)

	err = e.Setup(rotors)
	assert.Error(t, err, "should be in error")
}

func TestEnigma_Setup_Chgrotor(t *testing.T) {
	var rotors = []string{
		rVI,
		rII,
		rI,
	}

	e, _ := NewEnigma(3)
	// Now we have [rVI, rII, rI]
	err := e.Setup(rotors)

	rrVI, _ := NewRotor(rVI, false)

	assert.EqualValues(t, e.RS.R[0], rrVI, "should be equal")
	assert.NoError(t, err, "should not be in error")

	rotors = append(rotors, "JHJHSJDHJSHDKHDHKSHDKJSHDKJSHDKJHSKJDH")
	err = e.Setup(rotors)
	assert.Error(t, err, "should be in error")

}

func TestEnigma_Setup_Fourrotor(t *testing.T) {
	var rotors = []string{
		rVI,
		rII,
		rI,
	}

	e, _ := NewEnigma(3)

	rotors = append(rotors, "JHJHSJDHJSHDKHDHKSHDKJSHDKJSHDKJHSKJDH")
	err := e.Setup(rotors)
	assert.Error(t, err, "should be in error")

}

func TestEnigma_SetRotorSettings(t *testing.T) {
	// [rIII, rII, rI]
	var rotors = []string{
		rIII,
		rII,
		rI,
	}

	var set = []int{1, 4, 2}

	e, _ := NewEnigma(3)
	assert.Panics(t, func() {
		e.SetRotorSettings(set)
	}, "should panic")

	err := e.Setup(rotors)
	if err != nil {
		t.Error(err)
	}
	e.SetPlugboard(PBS)
	e.AddReflector(RfB)
	e.SetRotorSettings(set)

	// should have [rIII/1, rII/4, rI/2]
	assert.EqualValues(t, e.RS.R[0].index, 1, "should be equal")
	assert.EqualValues(t, e.RS.R[1].index, 4, "should be equal")
	assert.EqualValues(t, e.RS.R[2].index, 2, "should be equal")
}

func TestEnigma_SetPlugboard(t *testing.T) {

	e, _ := NewEnigma(3)
	e.SetPlugboard(PBS)

	realPBS := map[int]int{12: 19, 11: 20, 0: 16, 10: 22, 5: 23, 15: 18, 1: 3, 4: 9, 14: 24, 8: 21}

	assert.EqualValues(t, e.PlugBoard, realPBS, "should be equal")
}

func TestEnigma_AddReflector(t *testing.T) {
	var rotors = []string{
		rI,
		rII,
		rIII,
	}

	e, _ := NewEnigma(3)
	err := e.Setup(rotors)
	assert.NoError(t, err, "no error")
	assert.NotNil(t, e.RS, "rs not null")

	err = e.AddReflector(RfB)
	assert.NoError(t, err, "no error")
	assert.NotNil(t, e.Reflector, "refl not null")

	realPfB, _ := NewRotor(RfB, true)
	assert.Equal(t, e.Reflector, realPfB, "should be equal")

	arrayPfB := []int{24, 17, 20, 7, 16, 18, 11, 3, 15, 23, 13, 6, 14, 10, 12, 8, 4, 1, 5, 25, 2, 22, 21, 9, 0, 19}
	assert.EqualValues(t, e.Reflector.rotor, arrayPfB, "should be equal")
}

func TestEnigma_Encrypt(t *testing.T) {

}

func TestEnigma_Decrypt(t *testing.T) {

}

func TestEnigma_DumpState(t *testing.T) {

}
