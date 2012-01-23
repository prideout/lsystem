package vmath

//import "math"

type M4 struct {
    matrix [4 * 4]float64
}

func NewM4() *M4 {
    m := new(M4)
    m.SetIdentity()
    return m
}

func (m *M4) Copy() *M4 {
    n := NewM4()
    for i := 0; i < 4*4; i += 1 {
        n.matrix[i] = m.matrix[i]
    }
    return n
}

func (m *M4) SetIdentity() {
    m.matrix = [4 * 4]float64{
        1, 0, 0, 0,
        0, 1, 0, 0,
        0, 0, 1, 0,
        0, 0, 0, 1}
}
