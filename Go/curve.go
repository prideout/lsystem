package main

import (
    "bufio"
    "fmt"
    "os"
    "os/exec"
    ri "rman"
    "strings"
)

func initCamera() {
    ri.Projection("perspective", "fov", 30.)
    ri.Translate(0, -0.25, 5)
    ri.Rotate(-20, 1, 0, 0)
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
    var curveWidth float32 = 0.025
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
        "em", color(0, 0.65, 0.83),
        "samples", 64.)
    ri.TransformBegin()
    ri.Rotate(90, 1, 0, 0)
    //   ri.Disk(-0.7, 300, 360)
    ri.TransformEnd()

    // Sculpture
    ri.Attribute("identifier", "string name", "Sculpture")
    ri.Surface("Occlusion",
        "em", gray(1.1),
        "samples", 64.)
    ri.TransformBegin()
    ri.Rotate(90, 1, 0, 0)
    ri.Translate(0, 0, -0.55)
    drawCurve(curve)
    ri.TransformEnd()

    ri.WorldEnd()
}

func main() {
    xml := strings.NewReader(RIBBON)
    curve := Evaluate(xml)

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

const RIBBON string = `
<rules max_depth="3000">
    <rule name="entry">
    <call count="10" transforms="rz 36" rule="hbox"/>
    </rule>
    <rule name="hbox"><call rule="r"/></rule>
    <rule name="r"><call rule="turn3"/></rule>
    <rule name="turn3" max_depth="90" successor="r">
    <call rule="dbox"/>
    <call transforms="ry -0.5 tx 0.0125 sa 0.996" rule="turn3"/>
    </rule>
    <rule name="dbox">
    <instance transforms="s 0.55 2.0 1.25" shape="curve"/>
    </rule>
</rules>    
`
