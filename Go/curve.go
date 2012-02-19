package main

import (
    "bufio"
    "fmt"
    "os"
    "io/ioutil"
    "os/exec"
    ri "rman"
    "strings"
)

func initCamera() {
    ri.Projection("perspective", "fov", 30.)
    ri.Translate(0, -0.25, 10)
    ri.Rotate(-25, 1, 0, 0)
    ri.Rotate(180, 1, 0, 0)
    ri.Imager("Vignette")
}

func drawCurve(curve *Curve) {
    var normals, points []float32
    var vertsPerCurve []int

    var count int = 0
    for i, c := range *curve {
        marker := (c.N.X == 0 && c.N.Y == 0 && c.N.Z == 0)
        if i == len(*curve)-1 || marker {
            if count == 1 {
                n := len(points)
                points = points[:n-3]
                normals = normals[:n-3]
            } else if count > 1 {
                vertsPerCurve = append(vertsPerCurve, count)
            }
            count = 0
            continue
        }
        add := func(s *[]float32, x, y, z float32) {
            *s = append(append(append(*s, x), y), z)
        }
        add(&points, c.P.X, c.P.Y, c.P.Z)
        add(&normals, c.N.X, c.N.Y, c.N.Z)
        count++
    }

    // Validate that everything is kosher:
    var total int = 0
    for _, c := range vertsPerCurve {
        total += c
    }
    pointCount := len(points) / 3
    if total != pointCount {
        fmt.Printf("Incorrect curve topology: %d != %d\n",
            total, pointCount)
    }

    // Submit the gprim to prman:
    var curveWidth float32 = 0.2
    ri.Curves("linear", vertsPerCurve, "nonperiodic",
        "P", points,
        "N", normals,
        "constantwidth", curveWidth)
}

func drawWorld(curve *Curve) {
    ri.WorldBegin()

    ri.Declare("samples", "float")
    ri.Declare("em", "color")
    ri.Attribute("cull",
        "int backfacing", false,
        "int hidden", false)
    ri.Attribute("visibility",
        "int diffuse", true,
        "int specular", true)
    ri.Attribute("dice",
        "int rasterorient", false)

    // Floor
    ri.Attribute("identifier", "string name", "Floor")
    ri.Surface("Occlusion",
        "em", color(0.12, 0.12, 0.83),
        "samples", 16.)
    ri.TransformBegin()
    ri.Rotate(90, 1, 0, 0)
    ri.Disk(-0.7, 300, 360)
    ri.TransformEnd()

    // Sculpture
    ri.Attribute("identifier", "string name", "Sculpture")
    ri.Surface("Occlusion",
        "em", gray(1.0),
        "samples", 1024.)
    ri.TransformBegin()
    ri.Rotate(90, 1, 0, 0)
    ri.Translate(0, 0, -0.55)

    drawCurve(curve)

    ri.TransformEnd()

    ri.WorldEnd()
}

func usage() {
    fmt.Println("Usage: lsystem [file.xml] [single|multi] [render|ribgen|lsys]")
    os.Exit(1)
}

func main() {
        
    args := os.Args[1:]
    if len(args) != 3 {
        usage()
    }
    filename, threading, mode := args[0], args[1], args[2]

    var isThreading bool
    switch threading {
    case "single":
        isThreading = false
    case "multi":
        isThreading = true
    default:
        usage()
    }
    _ = isThreading // ignore for now

    var isSilent bool = false
    var isRendering bool
    switch mode {
    case "render":
        isRendering = true
    case "ribgen":
        isRendering = false
    case "lsys":
        isSilent = true
    }

    contents, _ := ioutil.ReadFile(filename)
    xml := strings.NewReader(string(contents))
    curve := Evaluate(xml)

    if isSilent {
        return
    }

    compileShader("Occlusion")
    compileShader("Vignette")

    var launch string
    if isRendering {
        launch = "launch:prman? -t -ctrl $ctrlin $ctrlout -capture debug.rib"
    } else {
        launch = ""
    }

    ri.Begin(launch)
    ri.Format(1920/2, 800/2, 1)
    ri.Display("grasshopper", "framebuffer", "rgba")
    ri.ShadingRate(1)
    ri.Option("limits", "int threads", 2)
    ri.Option("statistics",
        "xmlfilename", "stats.xml",
        "endofframe", true)
    ri.PixelSamples(4, 4)
    initCamera()
    drawWorld(&curve)
    ri.End()
}

func compileShader(name string) {
    rslFile := name + ".sl"
    cmd := exec.Command("shader", rslFile)
    b := bufio.NewWriter(os.Stdout)
    cmd.Stdout = b
    cmd.Stderr = b
    if nil != cmd.Run() {
        fmt.Print(ANSI_RED)
        b.Flush()
        fmt.Print(ANSI_RESET)
        os.Exit(1)
    }
}

func color(r, g, b float32) [3]float32 {
    return [3]float32{r, g, b}
}

func gray(x float32) [3]float32 {
    return [3]float32{x, x, x}
}

const (
    ANSI_BLACK   string = "\x1b[1;30m"
    ANSI_RED     string = "\x1b[1;31m"
    ANSI_GREEN   string = "\x1b[1;32m"
    ANSI_YELLOW  string = "\x1b[1;33m"
    ANSI_BLUE    string = "\x1b[1;34m"
    ANSI_MAGENTA string = "\x1b[1;35m"
    ANSI_CYAN    string = "\x1b[1;36m"
    ANSI_WHITE   string = "\x1b[1;37m"
    ANSI_RESET   string = "\x1b[0m"
)
