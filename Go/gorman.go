package gorman

// #include "ri_wrap.c"
import "C"

import (
	"unsafe"
	"fmt"
)

func DrawCurve() {
	C.curve(0.1, 0.04)
}

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

func unzipArgs(varargs ...interface{}) (names rtTokens, vals rtPointers, owned rawPointers) {

	if len(varargs)%2 != 0 {
		fmt.Println("odd number of arguments")
		return
	}

	nargs := len(varargs) / 2
	names = make(rtTokens, nargs)
	vals = make(rtPointers, nargs)
	owned = make(rawPointers, nargs)

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
		case int, int32:
			intified := v.(int)
			vals[i/2] = C.RtPointer(&intified)
		case float32:
			floatified := v.(float32)
			vals[i/2] = C.RtPointer(&floatified)
		case string:
			stringified := v.(string)
			token := C.CString(stringified)
			vals[i/2] = C.RtPointer(&token)
			owned[i/2] = unsafe.Pointer(token)
		default:
			fmt.Printf("'%s' has unknown type\n", pname)
			return
		}
	}
	return
}

func Display(name string, dtype string, mode string, varargs ...interface{}) {

	pName := C.CString(name)
	defer C.free(unsafe.Pointer(pName))

	pDtype := C.CString(dtype)
	defer C.free(unsafe.Pointer(pDtype))

	pMode := C.CString(mode)
	defer C.free(unsafe.Pointer(pMode))

	names, values, ownership := unzipArgs(varargs...)
	defer freeArgs(names, ownership)

	nargs := C.RtInt(len(varargs) / 2)
	pNames, pVals := safeArgs(names, values)
    C.RiDisplayV(pName, pDtype, pMode, nargs, pNames, pVals)
}

func Option(name string, varargs ...interface{}) {

	pName := C.CString(name)
	defer C.free(unsafe.Pointer(pName))

	names, values, ownership := unzipArgs(varargs...)
	defer freeArgs(names, ownership)

	nargs := C.RtInt(len(varargs) / 2)
    pNames, pVals := safeArgs(names, values)
	C.RiOptionV(pName, nargs, pNames, pVals)
}

func Format(width int32, height int32, aspectRatio float32) {
	w := C.RtInt(width)
	h := C.RtInt(height)
	r := C.RtFloat(aspectRatio)
	C.RiFormat(w, h, r)
}

func Begin(name string) {
	pName := C.CString(name)
	defer C.free(unsafe.Pointer(pName))
	C.RiBegin(pName)
}

func End() {
	C.RiEnd()
}

func Temp() {
	C.curve2(0.1, 0.04)
}
