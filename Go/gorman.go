package gorman

// #include "ri_wrap.c"
import "C"

func DrawCurve() {
    C.curve(0.1, 0.04)
}