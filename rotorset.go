package machine

import (
    "fmt"
)

type RotorSet struct {
    R      []*Rotor // rotor set (3-4)
    Double bool
    Len    int
}

func NewRotorSet(rotors []string) (rs *RotorSet, err error) {
    rs = &RotorSet{
        Len:    len(rotors),
        Double: false,
        R:      make([]*Rotor, len(rotors)),
    }

    for i, r := range rotors {
        rotor, err := NewRotor(r, false)
        if err == nil {
            rs.R[i] = rotor
        } else {
            return rs, err
        }
    }
    return rs, err
}

func (rs *RotorSet) Step() {
    var off= 0

    // if this is a 4-wheel machine, the foremost one (aka the 4th) does not move.
    if rs.Len == EnigmaMarine {
        off = 1
    }
    r0 := rs.R[off]     // r[1] if 4-rotor Enigma
    r1 := rs.R[off+1]   // r[2]
    r2 := rs.R[off+2]   // r[3]

    n2 := r2.NotchHit()
    r2.Step()

    if rs.Double {
        r1.Step()
        r0.Step()
        rs.Double = false
    } else {
        if n2 {
            // Check for double step
            r1.Step()
            if r1.NotchHit() {
                rs.Double = true
            }
        }
    }
}

func (rs *RotorSet) Set(init []int) (err error) {

    if rs.Len != len(init) {
        err = fmt.Errorf("incorrect settings for rotors")
    } else {
        for i, s := range init {
            rs.R[i].index = s
        }
        err = nil
    }
    return
}

func (rs *RotorSet) Settings() (set []int) {
    set = []int{}
    for _, s := range rs.R {
        set = append(set, s.index)
    }
    return
}

func (rs *RotorSet) Index() (set string) {

    var (
        off int
        ri0, ri1, ri2, ri3 string
    )

    if rs.Len == EnigmaMarine {
        ri3 = intToText[rs.R[0].index]
        off = 1
    } else {
        ri3 = "-"
        off = 0
    }

    ri2 = intToText[rs.R[off].index]
    ri1 = intToText[rs.R[off+1].index]
    ri0 = intToText[rs.R[off+2].index]
    set = fmt.Sprintf("%s%s%s%s", ri3, ri2, ri1, ri0)
    return
}

// left goes through all the rotors from right to left
func (rs * RotorSet) left(plain int) (cipher int) {
    next := plain
    for i := rs.Len - 1; i >= 0; i-- {
        next = rs.R[i].In(next)
    }
    cipher = next
    return
}

// right goes through all the rotors from left to right
func (rs * RotorSet) right(plain int) (cipher int) {
    next := plain
    for _, r := range rs.R {
        next = r.Out(next)
    }
    cipher = next
    return
}

