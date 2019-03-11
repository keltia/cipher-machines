package machine

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var rotors = []string{
	rI,
	rII,
	rIII,
}

var rotorsM = []string{
	rIV,
	rI,
	rII,
	rIII,
}

func TestNewRotorSet(t *testing.T) {
	rs, err := NewRotorSet(rotors)
	assert.NoError(t, err, "no error")
	assert.EqualValues(t, rs.Len, len(rotors), "should be equal")

	rrI, _ := NewRotor(rI, false)
	rrII, _ := NewRotor(rII, false)
	rrIII, _ := NewRotor(rIII, false)

	all := []*Rotor{
		rrI,
		rrII,
		rrIII,
	}

	for i, r := range rs.R {
		assert.NotNil(t, r, "should not be nil")
		assert.EqualValues(t, all[i], r, "should be equal")
	}
}

func TestNewRotorSet_M(t *testing.T) {
	rs, err := NewRotorSet(rotorsM)
	assert.NoError(t, err, "no error")
	assert.EqualValues(t, rs.Len, len(rotorsM), "should be equal")

	rrI, _ := NewRotor(rI, false)
	rrII, _ := NewRotor(rII, false)
	rrIII, _ := NewRotor(rIII, false)
	rrIV, _ := NewRotor(rIV, false)

	all := []*Rotor{
		rrIV,
		rrI,
		rrII,
		rrIII,
	}

	for i, r := range rs.R {
		assert.NotNil(t, r, "should not be nil")
		assert.EqualValues(t, all[i], r, "should be equal")
	}
}

func TestRotorSet_Set(t *testing.T) {
	rs, err := NewRotorSet(rotors)
	assert.NoError(t, err, "no error")

	rs.Set([]int{1, 2, 3})
	assert.EqualValues(t, []int{1, 2, 3}, rs.Settings(), "should be equal")
}

func TestRotorSet_Set_M(t *testing.T) {
	rs, err := NewRotorSet(rotorsM)
	assert.NoError(t, err, "no error")

	rs.Set([]int{1, 2, 3, 4})
	assert.EqualValues(t, []int{1, 2, 3, 4}, rs.Settings(), "should be equal")
}

func TestRotorSet_Set_Err(t *testing.T) {
	rs, err := NewRotorSet(rotorsM)
	assert.NoError(t, err, "no error")

	err = rs.Set([]int{1, 2, 3})
	assert.Error(t, err, "should error")
}

func TestRotorSet_Settings(t *testing.T) {
	rs, err := NewRotorSet(rotors)
	assert.NoError(t, err, "no error")

	settings := rs.Settings()
	assert.EqualValues(t, []int{0, 0, 0}, settings, "should be equal")

	rs.Set([]int{3, 13, 20})
	settings = rs.Settings()
	assert.EqualValues(t, []int{3, 13, 20}, settings, "should be equal")
}

func TestRotorSet_Settings_M(t *testing.T) {
	rs, err := NewRotorSet(rotorsM)
	assert.NoError(t, err, "no error")

	settings := rs.Settings()
	assert.EqualValues(t, []int{0, 0, 0, 0}, settings, "should be equal")

	rs.Set([]int{3, 13, 20, 1})
	settings = rs.Settings()
	assert.EqualValues(t, []int{3, 13, 20, 1}, settings, "should be equal")
}

func TestRotorSet_Step_Single(t *testing.T) {
	rs, err := NewRotorSet(rotors)
	assert.NoError(t, err, "no error")

	rs.Set([]int{0, 0, 20})
	index := rs.Index()
	assert.EqualValues(t, "-AAU", index, "should be equal")

	rs.Step()
	index = rs.Index()
	assert.EqualValues(t, "-AAV", index, "should be equal")

	rs.Step()
	index = rs.Index()
	assert.EqualValues(t, "-ABW", index, "should be equal")

	rs.Step()
	index = rs.Index()
	assert.EqualValues(t, "-ABX", index, "should be equal")
}

func TestRotorSet_Step_Single_M(t *testing.T) {
	rs, err := NewRotorSet(rotorsM)
	assert.NoError(t, err, "no error")

	rs.Set([]int{0, 0, 0, 20})
	index := rs.Index()
	assert.EqualValues(t, "AAAU", index, "should be equal")

	rs.Step()
	index = rs.Index()
	assert.EqualValues(t, "AAAV", index, "should be equal")

	rs.Step()
	index = rs.Index()
	assert.EqualValues(t, "AABW", index, "should be equal")

	rs.Step()
	index = rs.Index()
	assert.EqualValues(t, "AABX", index, "should be equal")
}

func TestRotorSet_Step_Double(t *testing.T) {
	rs, err := NewRotorSet(rotors)
	assert.NoError(t, err, "no error")

	rs.Set([]int{0, 3, 20})
	index := rs.Index()
	assert.EqualValues(t, "-ADU", index, "should be equal")

	rs.Step()
	index = rs.Index()
	assert.EqualValues(t, "-ADV", index, "should be equal")

	rs.Step()
	index = rs.Index()
	assert.EqualValues(t, "-AEW", index, "should be equal")

	rs.Step()
	index = rs.Index()
	assert.EqualValues(t, "-BFX", index, "should be equal")

	rs.Step()
	index = rs.Index()
	assert.EqualValues(t, "-BFY", index, "should be equal")
}

func TestRotorSet_Step_Double_M(t *testing.T) {
	rs, err := NewRotorSet(rotorsM)
	assert.NoError(t, err, "no error")

	rs.Set([]int{0, 0, 3, 20})
	index := rs.Index()
	assert.EqualValues(t, "AADU", index, "should be equal")

	rs.Step()
	index = rs.Index()
	assert.EqualValues(t, "AADV", index, "should be equal")

	rs.Step()
	index = rs.Index()
	assert.EqualValues(t, "AAEW", index, "should be equal")

	rs.Step()
	index = rs.Index()
	assert.EqualValues(t, "ABFX", index, "should be equal")

	rs.Step()
	index = rs.Index()
	assert.EqualValues(t, "ABFY", index, "should be equal")
}

func TestRotorSet_Index(t *testing.T) {
	rs, err := NewRotorSet(rotors)
	assert.NoError(t, err, "no error")

	index := rs.Index()
	assert.EqualValues(t, "-AAA", index, "should be equal")

	rs.Step()
	index = rs.Index()
	assert.EqualValues(t, "-AAB", index, "should be equal")

}

func TestRotorSet_Index_M(t *testing.T) {
	rs, err := NewRotorSet(rotorsM)
	assert.NoError(t, err, "no error")

	index := rs.Index()
	assert.EqualValues(t, "AAAA", index, "should be equal")

	rs.Step()
	index = rs.Index()
	assert.EqualValues(t, "AAAB", index, "should be equal")

}

func TestRotorSet_left(t *testing.T) {
	rs, err := NewRotorSet(rotors)
	assert.NoError(t, err, "no error")

	p := textToInt["A"]
	c := intToText[rs.left(p)]

	r2 := rs.R[2].In(p)
	r1 := rs.R[1].In(r2)
	r0 := rs.R[0].In(r1)

	assert.EqualValues(t, intToText[r0], c, "should be equal")
	assert.EqualValues(t, "K", c, "should be equal")
}

func TestRotorSet_right(t *testing.T) {
	rs, err := NewRotorSet(rotors)
	assert.NoError(t, err, "no error")

	p := textToInt["A"]
	c := intToText[rs.right(p)]

	r0 := rs.R[0].Out(p)
	r1 := rs.R[1].Out(r0)
	r2 := rs.R[2].Out(r1)

	assert.EqualValues(t, intToText[r2], c, "should be equal")
	assert.EqualValues(t, "G", c, "should be equal")
}
