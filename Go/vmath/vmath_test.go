// Emacs tricks for utf-8:
//
//    C-x RET f utf-8 RET
//    M-x ucs-insert GREEK <tab>
//    M-x ucs-insert DEVANGARI <tab>

package vmath_test

import (
    "math"
    "testing"
    "fmt"
    . "vmath"
)

var ᴨ float32 = float32(math.Atan(1) * 4)
var ε float32 = 1e-4

func TestVectors(t *testing.T) {

    i := V3New(1, 0, 0)
    j := V3New(0, 1, 0)
    k := V3New(0, 0, 1)

    // Right-hand rule:
    îĵ := i.Cross(j)
    if !îĵ.Equivalent(k, ε) {
        t.Fail()
    }

    // Rotation about Z
    M := M3RotateZ(ᴨ / 2)
    v := M.MulV3(i)
    if !v.Equivalent(j, ε) {
        t.Fail()
    }
    v = M.MulV3(j)
    if !v.Equivalent(V3New(-1, 0, 0), ε) {
        t.Fail()
    }

    // Rotation about Y
    M = M3RotateY(ᴨ / 2)
    v = M.MulV3(i)
	fmt.Println(v)
    if !v.Equivalent(k, ε) {
        t.Fail()
    }
}

func BenchmarkVectors(b *testing.B) {
	fmt.Println("No benchmarks yet.")
}
