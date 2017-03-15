package machine

import (
	"fmt"
)

// -- Rotor
type Rotor struct {
    rotor []int
    index int
    size  int
	notch int
	refl  bool
	wrap  bool
}

func (r *Rotor) testIndex(index int) bool {
	if index < 0 || index > r.size {
		return false
	}
	return true
}

func (r *Rotor) Start(index int) (err error) {
	if r.testIndex(index) {
    	r.index = index
		err = nil
	} else {
		err = fmt.Errorf("index out of bounds")
	}
	return
}

func (r *Rotor) Step() int {
	// If previous step wrapped, then it should not do so now
	if r.wrap {
		r.wrap = false
	}
	// We do not move reflectors
	if r.refl {
		return r.index
	}
	prev := r.index
    r.index = (r.index + 1) % r.size
	if r.index < prev {
		r.wrap = true
	}
	r.Rotate()
    return r.index
}

func (r *Rotor) Out(p int) (c int) {
    return r.rotor[p]
}

func (r *Rotor) In(p int) (c int) {
    for i, v := range r.rotor {
		if v == p {
	    	return i
		}
    }
    return -1
}

func (r *Rotor) Rotate() {
	// We do not move reflectors
	if r.refl {
		return
	}
	first := r.rotor[0]
	length := len(r.rotor)
	for i := 0; i < length - 1; i++ {
		r.rotor[i] = r.rotor[i + 1]
	}
	r.rotor[length - 1] = first
}

func (r *Rotor) HasWrapped() bool {
	return r.wrap
}

func NewRotor(str string, refl bool) (r *Rotor) {
	if len(str) != RotorSize {
		return nil
	}
    r = &Rotor{
		size: len(str),
		rotor: make([]int, len(str)),
		refl: refl,
	}

    for i, c := range str {
		r.rotor[i] = textToInt[string(c)]
    }
    return
}


