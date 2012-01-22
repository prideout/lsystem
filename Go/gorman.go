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

func freeArguments(namePointers []C.RtToken, ownedValues []unsafe.Pointer) {
	for _, v := range namePointers {
		C.free(unsafe.Pointer(v))
	}
	for _, v := range ownedValues {
		C.free(v)
	}
}

func unzipArguments(varargs ...interface{}) (namePointers []C.RtToken,
valuePointers []C.RtPointer,
ownedValues []unsafe.Pointer) {

	if len(varargs)%2 != 0 {
		fmt.Println("odd number of arguments")
		return
	}

	nargs := len(varargs) / 2
	namePointers = make([]C.RtToken, nargs)
	valuePointers = make([]C.RtPointer, nargs)
	ownedValues = make([]unsafe.Pointer, nargs)

	var pname string
	for i, v := range varargs {

		// Even-numbered arguments are parameter names
		if i%2 == 0 {
			if stringified, ok := v.(string); ok {
				pname = stringified
				token := C.CString(pname)
				namePointers[i/2] = token
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
			valuePointers[i/2] = C.RtPointer(&intified)
		case float32:
			floatified := v.(float32)
			valuePointers[i/2] = C.RtPointer(&floatified)
		case string:
			stringified := v.(string)
			token := C.CString(stringified)
			valuePointers[i/2] = C.RtPointer(&token)
			ownedValues[i/2] = unsafe.Pointer(token)
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

	names, values, ownership := unzipArguments(varargs...)
	defer freeArguments(names, ownership)

	nargs := C.RtInt(len(varargs) / 2)
	C.RiDisplayV(pName, pDtype, pMode, nargs, &names[0], &values[0])
}

func Option(name string, varargs ...interface{}) {

	pName := C.CString(name)
	defer C.free(unsafe.Pointer(pName))

	names, values, ownership := unzipArguments(varargs...)
	defer freeArguments(names, ownership)

	nargs := C.RtInt(len(varargs) / 2)
	C.RiOptionV(pName, nargs, &names[0], &values[0])
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
