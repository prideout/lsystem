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

// evaluates the rules in the given XML stream and returns a list of curves
func Evaluate(stream io.Reader) Curve {

    var curve Curve
    var lsys LSystem
    if err := xml.Unmarshal(stream, &lsys); err != nil {
        fmt.Println("Error parsing XML file:", err)
        return curve
    }

    xforms := make(MatrixCache)

    for _, rule := range lsys.Rules {
        for _, call := range rule.Calls {
            xforms.ParseString(call.Transforms)
        }
        for _, inst := range rule.Instances {
            xforms.ParseString(inst.Transforms)
        }
    }

    lsys.WeightSum = 0
    for _, rule := range lsys.Rules {
        if rule.Weight != 0 {
            lsys.WeightSum += rule.Weight
        } else {
            lsys.WeightSum++
        }
    }

    random := rand.New(rand.NewSource(42))
    start := StackNode{
        RuleIndex: lsys.PickRule("entry", random),
        Transform: vmath.M4Identity(),
    }

    lsys.ProcessRule(start, &curve, random)
    return curve
}

type LSystem struct {
    MaxDepth  int
    Rules     []Rule `xml:"rule"`
    WeightSum int
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
            if  rule.Successor != "" {
                next := StackNode{
                    RuleIndex: self.PickRule(rule.Successor, random),
                    Transform: matrix,
                }
                stack.Push(next)
            }
            *curve = append(*curve, CurvePoint{})
            continue
        }
        
    /*
        
        // Iterate through the operations in the rule
        pugi::xml_node op = entry.Node.first_child();
        for (; op; op = op.next_sibling()) {
            string xformString = op.attribute("transforms").value();
            Matrix4 xform = _ParseTransform(xformString);
            int count = op.attribute("count").as_int();
            while (count--) {
                matrix *= xform;
                string tag = op.name();
                if (tag == "call") {
                    const char* rule = op.attribute("rule").value();
                    pugi::xml_node newRule = _PickRule(_doc, rule, &seed);
                    StackEntry newEntry = { newRule, entry.Depth + 1, matrix };
                    stack.push_back(newEntry);
                } else if (tag == "instance") {
                    CurvePoint* cp = new CurvePoint;
                    Vector4 v = matrix * Vector4(0,0,0,1);
                    cp->P = Point3(v.getXYZ());
                    cp->N = matrix.getUpper3x3() * Vector3(0, 0, 1);
                    result->push_back(cp);
                    
                    // Give users a crude progress indicator by maintaining a count
                    // Note that size() on std::list can have linear complexity (!)
                    if (++_curveLength >= progressCount + 10000) {
                        diag::Print("%d curve segments so far...\n", _curveLength);
                        progressCount = _curveLength;
                    }

                } else {
                    op.print(std::cout);
                    diag::Print("Unrecognized operation: '%s'\n", tag.c_str());
                }
            }
        }
        */
    }
}

func (self *LSystem) PickRule(name string, random *rand.Rand) int {
    n := random.Intn(self.WeightSum)
    for i, rule := range self.Rules {
        weight := rule.Weight
        if weight == 0 {
            weight = 1
        }
        if n < weight {
            return i
        }
        n -= weight
    }
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
            xform = vmath.M4Mul(m, xform)
        case "t":
            x := readFloat()
            y := readFloat()
            z := readFloat()
            m := vmath.M4Translate(x, y, z)
            xform = vmath.M4Mul(m, xform)
        case "sa":
            a := readFloat()
            m := vmath.M4Scale(a, a, a)
            xform = vmath.M4Mul(m, xform)
        case "tx":
            x := readFloat()
            m := vmath.M4Translate(x, 0, 0)
            xform = vmath.M4Mul(m, xform)
        case "ty":
            y := readFloat()
            m := vmath.M4Translate(0, y, 0)
            xform = vmath.M4Mul(m, xform)
        case "tz":
            z := readFloat()
            m := vmath.M4Translate(0, 0, z)
            xform = vmath.M4Mul(m, xform)
        case "rx":
            x := readFloat()
            m := vmath.M4RotateX(x)
            xform = vmath.M4Mul(m, xform)
        case "ry":
            y := readFloat()
            m := vmath.M4RotateY(y)
            xform = vmath.M4Mul(m, xform)
        case "rz":
            z := readFloat()
            m := vmath.M4RotateZ(z)
            xform = vmath.M4Mul(m, xform)
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
