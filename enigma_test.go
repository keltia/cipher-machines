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

func TestEnigma_SetPlugboard(t *testing.T) {

    e, _ := NewEnigma(3)
    e.SetPlugboard(PBS)

    realPBS := map[int]int{12:19, 11:20, 0:16, 10:22, 5:23, 15:18, 1:3, 4:9, 14:24, 8:21}

    assert.EqualValues(t, e.PlugBoard, realPBS, "should be equal")
}

func TestEnigma_AddReflector(t *testing.T) {
    e, _ := NewEnigma(3)
    e.AddReflector(RfB)

    realPfA := *NewRotor(RfB, true)
    assert.Equal(t, e.Reflector, realPfA, "should be equal")

    arrayPfA := []int{24, 17, 20, 7, 16, 18, 11, 3, 15, 23, 13, 6, 14, 10, 12, 8, 4, 1, 5, 25, 2, 22, 21, 9, 0, 19}
    assert.EqualValues(t, e.Reflector.rotor, arrayPfA, "should be equal")
}
