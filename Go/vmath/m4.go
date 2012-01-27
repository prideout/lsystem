package vmath

import "fmt"

// Implements a 4x4 matrix type for 3D graphics.
// Much like go's string type, M4 is generally immutable.
// Unlike the V3 (et al) type, matrices use pass-by-pointer semantics.
// Unary operations are methods:
//    m = m.Transpose()
//    f := m.Derivative()
// Nullary and binary operations are functions:
//    i := M4Identity()
//    var x M4 = M4Mul(m, m)
//    var y V4 = M4MulV3(m, v)
//    var z M4 = M4MulT3(m, t)
//    scale := M4Scale(1.5)
type M4 struct {
    matrix [4 * 4]float32
}

// Create a 4x4 from the identity
func M4Identity() *M4 {
    m := new(M4)
    m.matrix = [4 * 4]float32{
        1, 0, 0, 0,
        0, 1, 0, 0,
        0, 0, 1, 0,
        0, 0, 0, 1}
    return m
}

// Create a 4x4 translation matrix
func M4Translate(x, y, z float32) *M4 {
    m := new(M4)
    m.matrix = [4 * 4]float32{
        1, 0, 0, 0,
        0, 1, 0, 0,
        0, 0, 1, 0,
        x, y, z, 1}
    return m
}

// Create a 4x4 scale matrix
func M4Scale(x, y, z float32) *M4 {
    m := new(M4)
    m.matrix = [4 * 4]float32{
        x, 0, 0, 0,
        0, y, 0, 0,
        0, 0, z, 0,
        0, 0, 0, 1}
    return m
}

// Create the product of two 4x4 matrices
func M4Mul(a *M4, b *M4) *M4 {
    m := new(M4)
    //    m.matrix = new([4 * 4]float32)
    for x := 0; x < 16; x += 4 {
        y, z, w := x+1, x+2, x+3
        m.matrix[x] = a.matrix[x]*b.matrix[0] +
            a.matrix[y]*b.matrix[4] +
            a.matrix[z]*b.matrix[8] +
            a.matrix[w]*b.matrix[12]
        m.matrix[y] = a.matrix[x]*b.matrix[1] +
            a.matrix[y]*b.matrix[5] +
            a.matrix[z]*b.matrix[9] +
            a.matrix[w]*b.matrix[13]
        m.matrix[z] = a.matrix[x]*b.matrix[2] +
            a.matrix[y]*b.matrix[6] +
            a.matrix[z]*b.matrix[10] +
            a.matrix[w]*b.matrix[14]
        m.matrix[w] = a.matrix[x]*b.matrix[3] +
            a.matrix[y]*b.matrix[7] +
            a.matrix[z]*b.matrix[11] +
            a.matrix[w]*b.matrix[15]
    }
    return m
}

// Create a 4x4 for rotation about the X-axis
func M4RotateX(radians float32) *M4 {
    m := new(M4)
    s, c := sin(radians), cos(radians)
    m.matrix = [4 * 4]float32{
        1, 0, 0, 0,
        0, c, s, 0,
        0, -s, c, 0,
        0, 0, 0, 1}
    return m
}

// Create a 4x4 for rotation about the Y-axis
func M4RotateY(radians float32) *M4 {
    m := new(M4)
    s, c := sin(radians), cos(radians)
    m.matrix = [4 * 4]float32{
        1, 0, 0, 0,
        c, 0, -s, 0,
        s, 0, c, 0,
        0, 0, 0, 1}
    return m
}

// Create a 4x4 for rotation about the Z-axis
func M4RotateZ(radians float32) *M4 {
    m := new(M4)
    s, c := sin(radians), cos(radians)
    m.matrix = [4 * 4]float32{
        c, s, 0, 0,
        -s, c, 0, 0,
        0, 0, 1, 0,
        0, 0, 0, 1}
    return m
}

// Create a duplicate of self
func (m *M4) Clone() *M4 {
    n := new(M4)
    for i := 0; i < 4*4; i += 1 {
        n.matrix[i] = m.matrix[i]
    }
    return n
}

func (m *M4) String() string {
    x := m.matrix
    return fmt.Sprintf("%f %f %f %f\n"+
        "%f %f %f %f\n"+
        "%f %f %f %f\n"+
        "%f %f %f %f\n",
        x[0], x[1], x[2], x[3],
        x[4], x[5], x[6], x[7],
        x[8], x[9], x[10], x[11],
        x[12], x[13], x[14], x[15])
}