package machine

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestNewEnigma(t *testing.T) {
    e, err := NewEnigma(3)
    assert.EqualValues(t, e.Size, 3, "should be equal")
    assert.Nil(t, err, "no error")

    e, err = NewEnigma(4)
    assert.Nil(t, err, "no error")
    assert.EqualValues(t, e.Size, 4, "should be equal")

    e, err = NewEnigma(666)
    assert.Error(t, err, "should protest")
    assert.Panics(t, func() {
        var _ = e.Size
    }, "should panic")
}

func TestEnigma_Setup(t *testing.T) {
    var rotors = []string{
        rI,
        rII,
        rIII,
    }

    e, _ := NewEnigma(3)
    e.Setup(rotors)

    assert.EqualValues(t, e.RotorSet[0], NewRotor(rI, false), "should be equal")
    assert.EqualValues(t, e.RotorSet[1], NewRotor(rII, false), "should be equal")
    assert.EqualValues(t, e.RotorSet[2], NewRotor(rIII, false), "should be equal")
    assert.Nil(t, e.PlugBoard, "should be nil")
    assert.Nil(t, e.Reflector, "should be nil")

    rotors[0] = "JHJHSJDHJSHDKHDHKSHDKJSHDKJSHDKJHSKJDH"
    err := e.Setup(rotors)
    assert.Error(t, err, "should be in error")

    rotors[0] = rVI
    err = e.Setup(rotors)
    assert.EqualValues(t, e.RotorSet[0], NewRotor(rVI, false), "should be equal")
    assert.NoError(t, err, "should not be in error")

    rotors = append(rotors, "JHJHSJDHJSHDKHDHKSHDKJSHDKJSHDKJHSKJDH")
    err = e.Setup(rotors)
    assert.Error(t, err, "should be in error")

}

func TestEnigma_SetRotorSettings(t *testing.T) {
    var rotors = []string{
        rI,
        rII,
        rIII,
    }

    var set = []int{ 1, 4, 2}

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

    assert.EqualValues(t, e.RotorSet[0].index, 1, "should be equal")
    assert.EqualValues(t, e.RotorSet[1].index, 4, "should be equal")
    assert.EqualValues(t, e.RotorSet[2].index, 2, "should be equal")
}

func TestEnigma_SetPlugboard(t *testing.T) {

    e, _ := NewEnigma(3)
    e.SetPlugboard(PBS)

    realPBS := map[int]int{12:19, 11:20, 0:16, 10:22, 5:23, 15:18, 1:3, 4:9, 14:24, 8:21}

    assert.EqualValues(t, e.PlugBoard, realPBS, "should be equal")
}

func TestEnigma_AddReflector(t *testing.T) {
    e, _ := NewEnigma(3)
    e.AddReflector(RfB)

    realPfA := NewRotor(RfB, true)
    assert.Equal(t, e.Reflector, realPfA, "should be equal")

    arrayPfA := []int{24, 17, 20, 7, 16, 18, 11, 3, 15, 23, 13, 6, 14, 10, 12, 8, 4, 1, 5, 25, 2, 22, 21, 9, 0, 19}
    assert.EqualValues(t, e.Reflector.rotor, arrayPfA, "should be equal")
}

func TestEnigma_Encrypt(t *testing.T) {

}

func TestEnigma_Decrypt(t *testing.T) {

}

func TestEnigma_DumpState(t *testing.T) {

}
