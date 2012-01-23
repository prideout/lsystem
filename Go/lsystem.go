package lsystem

import (
    "bufio"
    "encoding/xml"
    "fmt"
    "os"
    "strconv"
    "strings"
    "vmath"
)

type LSystem struct {
    MaxDepth int
    Rules    []Rule `xml:"rule"`
}

type Rule struct {
    Name      string     `xml:"name,attr"`
    Calls     []Call     `xml:"call"`
    Instances []Instance `xml:"instance"`
    MaxDepth  int        `xml:"max_depth,attr"`
    Successor string     `xml:"successor,attr"`
}

type Call struct {
    Transforms string `xml:"transforms,attr"`
    Rule       string `xml:"rule,attr"`
    Count      int    `xml:"count,attr"`
}

type Instance struct {
    Transforms string `xml:"transforms,attr"`
    Shape      string `xml:"string"`
}

type MatrixCache map[string]vmath.M4

// Parse a string in the xform language and add the resulting matrix to the lookup table.
// Examples:
//   "rx -2 tx 0.1 sa 0.996"
//   "s 0.55 2.0 1.25"
func (cache *MatrixCache) ParseString(s string) {
    if len(s) == 0 {
        return
    }
    reader := bufio.NewReader(strings.NewReader(s))

    readFloat := func() float32 {
        sx, _ := reader.ReadString(' ')
        fx, _ := strconv.ParseFloat(strings.TrimSpace(sx), 32)
        return float32(fx)
    }

    for {
        token, err := reader.ReadString(' ')
        token = strings.TrimSpace(token)
        switch token {
        case "s":
            x := readFloat()
            y := readFloat()
            z := readFloat()
            fmt.Printf("s %f %f %f\n", x, y, z)
        case "t":
            x := readFloat()
            y := readFloat()
            z := readFloat()
            fmt.Printf("t %f %f %f\n", x, y, z)
        case "sa":
            a := readFloat()
            fmt.Printf("sa %f\n", a)
        case "tx":
            x := readFloat()
            fmt.Printf("tx %f\n", x)
        case "ty":
            y := readFloat()
            fmt.Printf("ty %f\n", y)
        case "tz":
            z := readFloat()
            fmt.Printf("tz %f\n", z)
        case "rx":
            x := readFloat()
            fmt.Printf("rx %f\n", x)
        case "ry":
            y := readFloat()
            fmt.Printf("ry %f\n", y)
        case "rz":
            z := readFloat()
            fmt.Printf("rz %f\n", z)
        case "":
        default:
            fmt.Println("Unknown token: ", token)
        }
        if err != nil {
            break
        }
    }
    return
}

// evaluates the rules in the given XML file and returns a list of curves
func Evaluate(filename string) int {

    xmlFile, err := os.Open(filename)
    if err != nil {
        fmt.Println("Error opening XML file:", err)
        return 0
    }
    defer xmlFile.Close()

    var lsys LSystem
    if err = xml.Unmarshal(xmlFile, &lsys); err != nil {
        fmt.Println("Error parsing XML file:", err)
        return 0
    }

    var xforms MatrixCache

    for _, rule := range lsys.Rules {
        //fmt.Printf("Rule %d. %s -> %s\n", r, rule.Name, rule.Successor)
        for _, call := range rule.Calls {
            //fmt.Printf("%d. count=%d\ttransforms=%s\trule=%s\n", 
            //    i, call.Count, call.Transforms, call.Rule)
            xforms.ParseString(call.Transforms)
        }
        for _, inst := range rule.Instances {
            xforms.ParseString(inst.Transforms)
        }
    }

    return len(lsys.Rules)

}
