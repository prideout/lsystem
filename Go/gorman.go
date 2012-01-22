package gorman

// #include "ri_wrap.c"
import "C"

import (
    "unsafe"
//    "fmt"
)

func DrawCurve() {
    C.curve(0.1, 0.04)
}

// void RiOptionV(RtToken name, RtInt nParameters, RtToken nms[], RtPointer vals[]);

func Option(name string, param string, value int) {
    pName := C.RtToken(C.CString(name))
    param1 := C.RtToken(C.CString(param))
    value1 := C.RtInt(value)
    pValue1 := C.RtPointer(&value1)
    C.RiOptionV(pName, 1, &param1, &pValue1)
    C.free(unsafe.Pointer(pName))
    C.free(unsafe.Pointer(param1))
}

func Format(width int32, height int32, aspectRatio float32) {
    
    w := C.RtInt(width)
    h := C.RtInt(height)
    r := C.RtFloat(aspectRatio)
    C.RiFormat(w, h, r)
}

func Begin(name string) {
    pName := C.RtToken(C.CString(name))
    C.RiBegin(pName)
    C.free(unsafe.Pointer(pName))
}

func End() {
    C.RiEnd()
}

func Temp() {
    C.curve2(0.1, 0.04)
}
