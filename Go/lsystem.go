package main

import (
    "bufio"
    "encoding/xml"
    "fmt"
    "io"
    "math/rand"
    "strconv"
    "strings"
    "vmath"
)

type CurvePoint struct {
    P   vmath.P3
    N   vmath.V3
}

type Curve []CurvePoint

// Evaluates the rules in the given XML stream and returns a list of curves
func Evaluate(stream io.Reader) Curve {

    var curve Curve
    var lsys LSystem
    if err := xml.Unmarshal(stream, &lsys); err != nil {
        fmt.Println("Error parsing XML file:", err)
        return curve
    }

    // Parse the transform strings
    lsys.Matrices = make(MatrixCache)
    for _, rule := range lsys.Rules {
        for _, call := range rule.Calls {
            lsys.Matrices.ParseString(call.Transforms)
        }
        for _, inst := range rule.Instances {
            lsys.Matrices.ParseString(inst.Transforms)
        }
    }
    lsys.Matrices[""] = *vmath.M4Identity()

    // Provide defaults
    for i := range lsys.Rules {
        r := &lsys.Rules[i]
        if r.Weight == 0 {
            r.Weight = 1
        }
    }

    // Seed the random-number generator and pick the entry rule
    random := rand.New(rand.NewSource(42))
    start := StackNode{
        RuleIndex: lsys.PickRule("entry", random),
        Transform: vmath.M4Identity(),
    }

    lsys.ProcessRule(start, &curve, random)
    return curve
}

type LSystem struct {
    MaxDepth int    `xml:"max_depth,attr"`
    Rules    []Rule `xml:"rule"`
    Matrices MatrixCache
}

type Rule struct {
    Name      string     `xml:"name,attr"`
    Calls     []Call     `xml:"call"`
    Instances []Instance `xml:"instance"`
    MaxDepth  int        `xml:"max_depth,attr"`
    Successor string     `xml:"successor,attr"`
    Weight    int        `xml:"weight,attr"`
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

func radians(degrees float32) float32 {
    return degrees * 3.1415926535 / 180.0
}

func (self *LSystem) ProcessRule(start StackNode, curve *Curve, random *rand.Rand) {

    stack := new(Stack)
    stack.Push(start)

    for stack.Len() > 0 {

        e := stack.Pop()

        localMax := self.MaxDepth
        rule := self.Rules[e.RuleIndex]

        if rule.MaxDepth > 0 {
            localMax = rule.MaxDepth
        }

        if stack.Len() >= self.MaxDepth {
            *curve = append(*curve, CurvePoint{})
            continue
        }

        matrix := e.Transform
        if e.Depth >= localMax {
            // Switch to a different rule is one is specified
            if rule.Successor != "" {
                next := StackNode{
                    RuleIndex: self.PickRule(rule.Successor, random),
                    Transform: matrix.Clone(),
                }
                stack.Push(next)
            }
            *curve = append(*curve, CurvePoint{})
            continue
        }

        for _, call := range rule.Calls {
            t, ok := self.Matrices[call.Transforms]
            if !ok {
                fmt.Println("Unable to compute matrix for: ", call.Transforms)
            }
            count := call.Count
            if count == 0 {
                count = 1
            }
            for ; count > 0; count-- {
                matrix = t.Compose(matrix)
                newRule := self.PickRule(call.Rule, random)
                next := StackNode{
                    RuleIndex: newRule,
                    Depth:     e.Depth + 1,
                    Transform: matrix.Clone(),
                }
                stack.Push(next)
            }
        }

        for _, instance := range rule.Instances {
            t := self.Matrices[instance.Transforms]
            matrix = t.Compose(matrix)
            v := vmath.V4{0, 0, 0, 1}
            n := vmath.V3{0, 0, 1}
            p := vmath.P3FromV4(matrix.Mul(v))
            n = matrix.GetUpperLeft().Mul(n)
            c := CurvePoint{P: p, N: n}
            *curve = append(*curve, c)
            if len(*curve)%10000 == 0 {
                fmt.Printf("Instanced %d nodes (stack = %d)\n", len(*curve), stack.Len())
            }
        }
    }
}

func (self *LSystem) PickRule(name string, random *rand.Rand) int {

    // Sum up the weights of all rules with this name:
    var sum int = 0
    for _, rule := range self.Rules {
        if rule.Name != name {
            continue
        }
        weight := rule.Weight
        sum += weight
    }

    // Choose a rule at random:
    n := random.Intn(sum)
    for i, rule := range self.Rules {
        if rule.Name != name {
            continue
        }
        weight := rule.Weight
        if n < weight {
            return i
        }
        n -= weight
    }

    fmt.Println("Error.")
    return -1
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

    xform := vmath.M4Identity()

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
            m := vmath.M4Scale(x, y, z)
            xform = m.Compose(xform)
        case "sa":
            a := readFloat()
            m := vmath.M4Scale(a, a, a)
            xform = m.Compose(xform)
        case "t":
            x := readFloat()
            y := readFloat()
            z := readFloat()
            m := vmath.M4Translate(x, y, z)
            xform = m.Compose(xform)
        case "tx":
            x := readFloat()
            m := vmath.M4Translate(x, 0, 0)
            xform = m.Compose(xform)
        case "ty":
            y := readFloat()
            m := vmath.M4Translate(0, y, 0)
            xform = m.Compose(xform)
        case "tz":
            z := readFloat()
            m := vmath.M4Translate(0, 0, z)
            xform = m.Compose(xform)
        case "rx":
            x := radians(readFloat())
            m := vmath.M4RotateX(-x)
            xform = m.Compose(xform)
        case "ry":
            y := radians(readFloat())
            m := vmath.M4RotateY(-y)
            xform = m.Compose(xform)
        case "rz":
            z := radians(readFloat())
            m := vmath.M4RotateZ(-z)
            xform = m.Compose(xform)
        case "":
        default:
            fmt.Println("Unknown token: ", token)
        }
        if err != nil {
            break
        }
    }
    (*cache)[s] = *xform
    return
}
