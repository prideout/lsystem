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

func Option(name string, varargs ...interface{}) {

	if len(varargs)%2 != 0 {
		fmt.Println("odd number of arguments")
		return
	}

	nargs := len(varargs) / 2
	namePointers := make([]C.RtToken, nargs)
	valuePointers := make([]C.RtPointer, nargs)

	var pname string
	for i, v := range varargs {

		// Even-numbered arguments are parameter names
		if i%2 == 0 {
			if stringified, ok := v.(string); ok {
				pname = stringified
				token := C.RtToken(C.CString(pname))
				namePointers[i/2] = token
				defer C.free(unsafe.Pointer(token))
			} else {
				fmt.Printf("argument %d is not a string\n", i)
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
			token := C.RtToken(C.CString(stringified))
			defer C.free(unsafe.Pointer(token))
			valuePointers[i/2] = C.RtPointer(&token)
		default:
			fmt.Printf("'%s' has unknown type\n", pname)
		}
	}

	token := C.RtToken(C.CString(name))
	defer C.free(unsafe.Pointer(token))
	C.RiOptionV(token, C.RtInt(nargs), &namePointers[0], &valuePointers[0])
}

// void RiOptionV(RtToken name, RtInt nParameters, RtToken nms[], RtPointer vals[]);
func Option1(name string, param string, value int) {
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
