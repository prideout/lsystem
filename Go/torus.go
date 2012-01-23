package main

import (
    "bufio"
    "exec"
    "fmt"
    ri "gorman"
    "os"
)

func main() {

    compileShader("Occlusion")
    compileShader("Vignette")

    launch := "launch:prman? -t -ctrl $ctrlin $ctrlout -capture debug.rib"
    ri.Begin(launch)
    ri.Format(512, 320, 1)
    ri.Display("grasshopper", "framebuffer", "rgba")
    ri.ShadingRate(4)
    ri.Option("limits", "int threads", 2)
    ri.Option("statistics",
        "xmlfilename", "stats.xml",
        "endofframe", true)
    ri.PixelSamples(4, 4)
    ri.Projection("perspective", "fov", 30.)
    ri.Translate(0, 0.25, 4)
    ri.Rotate(-20, 1, 0, 0)
    ri.Rotate(180, 1, 0, 0)
    ri.Imager("Vignette")
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
        "em", color(0, 0.65, 0.83),
        "samples", 64.)
    ri.TransformBegin()
    ri.Rotate(90, 1, 0, 0)
    ri.Disk(-0.7, 300, 360)
    ri.TransformEnd()

    // Sculpture
    ri.Attribute("identifier", "string name", "Sculpture")
    ri.Surface("Occlusion",
        "em", gray(1.1),
        "samples", 64.)
    ri.TransformBegin()
    ri.Translate(0, 0.75, 0)
    ri.Torus(1, 0.2, 0, 360, 360)
    ri.TransformEnd()
    ri.WorldEnd()

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
    return [3]float32{r, b, g}
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
