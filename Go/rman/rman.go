package rman

// #include <stdlib.h>
// #include "ri.h"
import "C"

import (
    "fmt"
    "unsafe"
)

func Display(name string, dtype string, mode string, varargs ...interface{}) {
    pName := C.CString(name)
    defer C.free(unsafe.Pointer(pName))
    pDtype := C.CString(dtype)
    defer C.free(unsafe.Pointer(pDtype))
    pMode := C.CString(mode)
    defer C.free(unsafe.Pointer(pMode))
    names, values, ownership := unzipArgs(varargs...)
    defer freeArgs(names, ownership)
    nArgs := C.RtInt(len(varargs) / 2)
    pNames, pVals := safeArgs(names, values)
    C.RiDisplayV(pName, pDtype, pMode, nArgs, pNames, pVals)
}

func Disk(height float32, radius float32, tmax float32, varargs ...interface{}) {
    h := C.RtFloat(height)
    r := C.RtFloat(radius)
    t := C.RtFloat(tmax)
    names, values, ownership := unzipArgs(varargs...)
    defer freeArgs(names, ownership)
    nArgs := C.RtInt(len(varargs) / 2)
    pNames, pVals := safeArgs(names, values)
    C.RiDiskV(h, r, t, nArgs, pNames, pVals)
}

func Torus(majrad float32, minrad float32,
    phimin float32, phimax float32,
    tmax float32, varargs ...interface{}) {
    a := C.RtFloat(majrad)
    i := C.RtFloat(minrad)
    p := C.RtFloat(phimin)
    x := C.RtFloat(phimax)
    t := C.RtFloat(tmax)
    names, values, ownership := unzipArgs(varargs...)
    defer freeArgs(names, ownership)
    nArgs := C.RtInt(len(varargs) / 2)
    pNames, pVals := safeArgs(names, values)
    C.RiTorusV(a, i, p, x, t, nArgs, pNames, pVals)
}

func Curves(basis string, spans []int,
    wrap string, varargs ...interface{}) {
    pBasis := C.CString(basis)
    defer C.free(unsafe.Pointer(pBasis))
    pWrap := C.CString(wrap)
    defer C.free(unsafe.Pointer(pWrap))
    names, values, ownership := unzipArgs(varargs...)
    defer freeArgs(names, ownership)
    nArgs := C.RtInt(len(varargs) / 2)
    nSpans := C.RtInt(len(spans))
    pSpans := (*C.RtInt)(C.RtPointer(&spans[0]))
    pNames, pVals := safeArgs(names, values)
    C.RiCurvesV(pBasis, nSpans, pSpans,
        pWrap, nArgs, pNames, pVals)
}

func Attribute(name string, varargs ...interface{}) {
    pName := C.CString(name)
    defer C.free(unsafe.Pointer(pName))
    names, values, ownership := unzipArgs(varargs...)
    defer freeArgs(names, ownership)
    nArgs := C.RtInt(len(varargs) / 2)
    pNames, pVals := safeArgs(names, values)
    C.RiAttributeV(pName, nArgs, pNames, pVals)
}

func Option(name string, varargs ...interface{}) {
    pName := C.CString(name)
    defer C.free(unsafe.Pointer(pName))
    names, values, ownership := unzipArgs(varargs...)
    defer freeArgs(names, ownership)
    nArgs := C.RtInt(len(varargs) / 2)
    pNames, pVals := safeArgs(names, values)
    C.RiOptionV(pName, nArgs, pNames, pVals)
}

func Projection(name string, varargs ...interface{}) {
    pName := C.CString(name)
    defer C.free(unsafe.Pointer(pName))
    names, values, ownership := unzipArgs(varargs...)
    defer freeArgs(names, ownership)
    nArgs := C.RtInt(len(varargs) / 2)
    pNames, pVals := safeArgs(names, values)
    C.RiProjectionV(pName, nArgs, pNames, pVals)
}

func Imager(name string, varargs ...interface{}) {
    pName := C.CString(name)
    defer C.free(unsafe.Pointer(pName))
    names, values, ownership := unzipArgs(varargs...)
    defer freeArgs(names, ownership)
    nArgs := C.RtInt(len(varargs) / 2)
    pNames, pVals := safeArgs(names, values)
    C.RiImagerV(pName, nArgs, pNames, pVals)
}

func Surface(name string, varargs ...interface{}) {
    pName := C.CString(name)
    defer C.free(unsafe.Pointer(pName))
    names, values, ownership := unzipArgs(varargs...)
    defer freeArgs(names, ownership)
    nArgs := C.RtInt(len(varargs) / 2)
    pNames, pVals := safeArgs(names, values)
    C.RiSurfaceV(pName, nArgs, pNames, pVals)
}

func Shader(name string, handle string, varargs ...interface{}) {
    pName := C.CString(name)
    defer C.free(unsafe.Pointer(pName))
    pHandle := C.CString(handle)
    defer C.free(unsafe.Pointer(pHandle))
    names, values, ownership := unzipArgs(varargs...)
    defer freeArgs(names, ownership)
    nArgs := C.RtInt(len(varargs) / 2)
    pNames, pVals := safeArgs(names, values)
    C.RiShaderV(pName, pHandle, nArgs, pNames, pVals)
}

func Format(width int32, height int32, aspectRatio float32) {
    w := C.RtInt(width)
    h := C.RtInt(height)
    r := C.RtFloat(aspectRatio)
    C.RiFormat(w, h, r)
}

func Begin(name string) {
    if name == "" {
        C.RiBegin(nil)
        return
    }
    pName := C.CString(name)
    defer C.free(unsafe.Pointer(pName))
    C.RiBegin(pName)
}

func Declare(name string, decl string) {
    pName := C.CString(name)
    defer C.free(unsafe.Pointer(pName))
    pDecl := C.CString(decl)
    defer C.free(unsafe.Pointer(pDecl))
    C.RiDeclare(pName, pDecl)
}

func End() {
    C.RiEnd()
}

func WorldBegin() {
    C.RiWorldBegin()
}

func WorldEnd() {
    C.RiWorldEnd()
}

func TransformBegin() {
    C.RiTransformBegin()
}

func TransformEnd() {
    C.RiTransformEnd()
}

func ShadingRate(val float32) {
    C.RiShadingRate(C.RtFloat(val))
}

func PixelSamples(x, y float32) {
    C.RiPixelSamples(C.RtFloat(x), C.RtFloat(y))
}

func Translate(x, y, z float32) {
    C.RiTranslate(C.RtFloat(x), C.RtFloat(y), C.RtFloat(z))
}

func Rotate(angle, x, y, z float32) {
    C.RiRotate(C.RtFloat(angle), C.RtFloat(x), C.RtFloat(y), C.RtFloat(z))
}

// Private utility functions for managing name-value pairs \\

type rtTokens []C.RtToken
type rtPointers []C.RtPointer
type rawPointers []unsafe.Pointer

func freeArgs(names rtTokens, owned rawPointers) {
    for _, v := range names {
        C.free(unsafe.Pointer(v))
    }
    for _, v := range owned {
        C.free(v)
    }
}

func safeArgs(names rtTokens, values rtPointers) (*C.RtToken, *C.RtPointer) {
    if len(names) == 0 {
        return new(C.RtToken), new(C.RtPointer)
    }
    return &names[0], &values[0]
}

func unzipArgs(varargs ...interface{}) (names rtTokens,
    vals rtPointers, owned rawPointers) {

    if len(varargs)%2 != 0 {
        fmt.Println("odd number of arguments")
        return
    }

    nArgs := len(varargs) / 2
    names = make(rtTokens, nArgs)
    vals = make(rtPointers, nArgs)
    owned = make(rawPointers, nArgs)

    var pname string
    for i, v := range varargs {

        // Even-numbered arguments are parameter names
        if i%2 == 0 {
            if stringified, ok := v.(string); ok {
                pname = stringified
                token := C.CString(pname)
                names[i/2] = token
            } else {
                fmt.Printf("argument %d is not a string\n", i)
                return
            }
            continue
        }

        // Odd-numbered arguments are values
        switch v.(type) {
        case bool:
            boolified := v.(bool)
            var intified int32 = 0
            if boolified {
                intified = 1
            }
            vals[i/2] = C.RtPointer(&intified)
        case int, int32:
            intified := v.(int)
            vals[i/2] = C.RtPointer(&intified)
        case [3]float32:
            floatified := v.([3]float32)
            vals[i/2] = C.RtPointer(&floatified[0])
        case [4]float32:
            floatified := v.([4]float32)
            vals[i/2] = C.RtPointer(&floatified[0])
        case []float32:
            floatified := v.([]float32)
            vals[i/2] = C.RtPointer(&floatified[0])
        case float32:
            floatified := v.(float32)
            vals[i/2] = C.RtPointer(&floatified)
        case float64:
            floatified := float32(v.(float64))
            vals[i/2] = C.RtPointer(&floatified)
        case string:
            stringified := v.(string)
            token := C.CString(stringified)
            vals[i/2] = C.RtPointer(&token)
            owned[i/2] = unsafe.Pointer(token)
        default:
            m := fmt.Sprintf("'%s' has unknown type %T\n", pname, v)
            panic(m)
        }
    }
    return
}
