// http://stackoverflow.com/questions/369438/smooth-spectrum-for-mandelbrot-set-rendering

// This code was originally contributed by Dr. David Alan Gilbert on
// thr golang ML
//
// http://groups.google.com/group/golang-nuts/browse_thread/thread/e5a65b94b27a1916/838087508caf0873?lnk=gst&q=mandelbrot#838087508caf0873
//
// Dave Gilbert's parallel mandelbrot.  Inspiration and the odd line of code
// drawn from the mandelbrot bench in the go distrobution
// d...@treblig.org

package main

import (
	"time"
	"flag"
	"image"
	"fmt"
)

var size = flag.Int("size", 256, "Size of the image")
var left = flag.Float("left", -2.0, "Left extent of brot")
var right = flag.Float("right", 2.0, "Right extent of brot")
var top = flag.Float("top", -2.0, "Top extent of brot")
var bottom = flag.Float("bottom", 2.0, "Bottom extent of brot")
var maxIter = flag.Int("maxIter", 255, "Maximum number of iteractions")

var im *image.RGBA

func doPoint(x, y int, xval, yval float) {
        var tmpcol image.RGBAColor;
        tmpcol.A = 255;

        var Ti, Tr, Zi, Zr float;
        var Ci, Cr float;
        var iter int;

        Cr = xval;
        Ci = yval;

        for iter = 0; (iter <= *maxIter) && (Ti+Tr < 4.0); iter++ {
                // These 4 lines from the C benchmark version
                Zi = 2.0*Zr*Zi + Ci;
                Zr = Tr - Ti + Cr;
                Tr = Zr * Zr;
                Ti = Zi * Zi;
        }

        // Try and get some basic interesting colour
        tmpcol.R = uint8(iter % 256);
        tmpcol.G = uint8((iter & 15 << 4) | ((iter & 240) / 16));
        tmpcol.B = tmpcol.R;
        im.Set(x, y, tmpcol);

}

func doRow(y int, yval float, ack chan bool) {
        for x := 0; x < *size; x++ {
                doPoint(x, y, valInRange(*left, *right, *size, x), yval)
        }

        ack <- true;

}

// Returns a value representing point 'offset' within a range between start and end where there are 'size' points
func valInRange(start, end float, size, offset int) float {
        var ran = end - start;

        return start + (ran*float(offset))/float(size);

}

func main() {
        flag.Parse();

        im = image.NewRGBA(*size, *size);

        donechan := make(chan bool, *size);
	
	start := time.Nanoseconds()

        for y := 0; y < *size; y++ {
                go doRow(y, valInRange(*top, *bottom, *size, y), donechan)
        }

        for y := 0; y < *size; y++ {
                <-donechan
        }

	end := time.Nanoseconds()

	fmt.Printf("%f", float(end - start) / float(10E9))
}