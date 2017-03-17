package machine

import (
	"fmt"
	"strings"
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

func (r *Rotor) turn(index int) (ret int, wrap bool) {
	prev := index
	index = (index + 1) % r.size
	if index < prev {
		wrap = true
	}
	return index, wrap
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
	r.index, r.wrap = r.turn(r.index)
	r.Rotate()
	r.notch, _ = r.turn(r.notch)
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

func NewRotor(str string, refl bool) (r *Rotor, err error) {
	// Plain rotor descriptions have one or two more bytes because of the notches
	// Reflectors' are just RotorSize
	if refl {
		if len(str) != RotorSize ||
			strings.ContainsRune(str, '/') {
			return &Rotor{}, fmt.Errorf("bad description for reflector: %s", str)
		}
	} else {
		// Check for notches
		if !strings.ContainsRune(str, '/') &&
			len(str) != RotorSize + 1 {
			return &Rotor{}, fmt.Errorf("bad description for rtor: %s", str)
		}
	}
    r = &Rotor{
		size: RotorSize,
		rotor: make([]int, RotorSize),
		refl: refl,
	}

	// Transform data
	var i = 0
    for _, c := range str {
		// Identify the notch that make the next rotor advance
		if c == '/' {
			r.notch = i
		} else {
			r.rotor[i] = textToInt[string(c)]
			i++
		}
    }
	err = nil
    return
}


