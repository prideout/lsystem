package vmath

//import "math"

// https://bitbucket.org/prideout/pez-viewer/src/11899f6b6f02/vmath.h

// Implements a 4x4 matrix type for 3D graphics.
// Much like go's string type, M4 is generally immutable.
// Nullary functions like transpose and derivative are methods:
//    m = m.Transpose()
//    f := m.Derivative()
// All other operations are functions:
//    var x M4 := M4Mul(m, m)
//    var y V4 = M4MulV3(m, v)
//    var z M4 := M4MulT3(m, t)
//    scale := M4Scale(1.5)
type M4 struct {
    matrix [4 * 4]float64
}

// Create a 4x4 from the identity
func M4Identity() *M4 {
    m := new(M4)
    m.matrix = [4 * 4]float64{
        1, 0, 0, 0,
        0, 1, 0, 0,
        0, 0, 1, 0,
        0, 0, 0, 1}
    return m
}

// Create a 4x4 translation matrix
func M4Translate(x, y, z float32) *M4 {
    m := new(M4)
    m.matrix = [4 * 4]float64{
        1, 0, 0, 0,
        0, 1, 0, 0,
        0, 0, 1, 0,
        x, y, z, 1}
    return m
}

// Create a 4x4 scale matrix
func M4Scale(x, y, z float32) *M4 {
    m := new(M4)
    m.matrix = [4 * 4]float64{
        x, 0, 0, 0,
        0, y, 0, 0,
        0, 0, z, 0,
        0, 0, 0, 1}
    return m
}

// Create the product of two 4x4 matrices
func M4Mul(a *M4, b *M4) *M4 {
    m := new(M4)
    m.matrix = new([4 * 4]float64)
    /*
    m.x.x = a.x.x * b.x.x + x.y * b.y.x + a.x.z * b.z.x + a.x.w * b.w.x
    m.x.y = a.x.x * b.x.y + x.y * b.y.y + a.x.z * b.z.y + a.x.w * b.w.y
    m.x.z = a.x.x * b.x.z + x.y * b.y.z + a.x.z * b.z.z + a.x.w * b.w.z
    m.x.w = a.x.x * b.x.w + x.y * b.y.w + a.x.z * b.z.w + a.x.w * b.w.w
    m.y.x = a.y.x * b.x.x + y.y * b.y.x + a.y.z * b.z.x + a.y.w * b.w.x
    m.y.y = a.y.x * b.x.y + y.y * b.y.y + a.y.z * b.z.y + a.y.w * b.w.y
    m.y.z = a.y.x * b.x.z + y.y * b.y.z + a.y.z * b.z.z + a.y.w * b.w.z
    m.y.w = a.y.x * b.x.w + y.y * b.y.w + a.y.z * b.z.w + a.y.w * b.w.w
    m.z.x = a.z.x * b.x.x + z.y * b.y.x + a.z.z * b.z.x + a.z.w * b.w.x
    m.z.y = a.z.x * b.x.y + z.y * b.y.y + a.z.z * b.z.y + a.z.w * b.w.y
    m.z.z = a.z.x * b.x.z + z.y * b.y.z + a.z.z * b.z.z + a.z.w * b.w.z
    m.z.w = a.z.x * b.x.w + z.y * b.y.w + a.z.z * b.z.w + a.z.w * b.w.w
    m.w.x = a.w.x * b.x.x + w.y * b.y.x + a.w.z * b.z.x + a.w.w * b.w.x
    m.w.y = a.w.x * b.x.y + w.y * b.y.y + a.w.z * b.z.y + a.w.w * b.w.y
    m.w.z = a.w.x * b.x.z + w.y * b.y.z + a.w.z * b.z.z + a.w.w * b.w.z
    m.w.w = a.w.x * b.x.w + w.y * b.y.w + a.w.z * b.z.w + a.w.w * b.w.w
    */
    return m
}

// Create a 4x4 for rotation about the X-axis
func M4RotateX(radians float32) *M4 {
    m := new(M4)
    // TODO
    return m
}

// Create a 4x4 for rotation about the Y-axis
func M4RotateY(radians float32) *M4 {
    m := new(M4)
    // TODO
    return m
}

// Create a 4x4 for rotation about the Z-axis
func M4RotateZ(radians float32) *M4 {
    m := new(M4)
    // TODO
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
