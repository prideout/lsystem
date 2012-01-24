package vmath

import "fmt"

// https://bitbucket.org/prideout/pez-viewer/src/11899f6b6f02/vmath.h

type P3 struct {
    X, Y, Z float32
}

func P3New(x, y, z float32) P3 {
    v := new(P3)
    v.X = x
    v.Y = y
    v.Z = z
    return *v
}

func P3Distance(a P3, b P3) float32 {
    return P3Sub(a, b).Length()
}

func P3Add(a P3, b V3) P3 {
    return P3New(
        a.X+b.X,
        a.Y+b.Y,
        a.Z+b.Z)
}

func P3Sub(a P3, b P3) V3 {
    return V3New(
        a.X-b.X,
        a.Y-b.Y,
        a.Z-b.Z)
}

func (v P3) Clone() P3 {
    return P3New(v.X, v.Y, v.Z)
}

func (v P3) String() string {
    return fmt.Sprint(v.X, ", ", v.Y, ", ", v.Z)
}

