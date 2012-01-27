package main

import (
    "bufio"
    "fmt"
    "os"
    "os/exec"
    "strings"
    //    ri "gorman"
)

func main() {
    xml := strings.NewReader(RIBBON)
    foo := Evaluate(xml)
    fmt.Print(foo)
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

const RIBBON string = `<rules max_depth="3000">
    <rule name="entry">
        <!--call count="72" transforms="rz 10" rule="hbox"/-->
        <call count="144" transforms="rz 5" rule="hbox"/>
    </rule>
    <rule name="hbox"><call rule="r"/></rule>
    <rule name="r"><call rule="forward"/></rule>
    <rule name="r"><call rule="turn"/></rule>
    <rule name="r"><call rule="turn2"/></rule>
    <rule name="r"><call rule="turn4"/></rule>
    <rule name="r"><call rule="turn3"/></rule>
    <rule name="forward" max_depth="90" successor="r">
        <call rule="dbox"/>
        <call transforms="rz 5.6 tx 0.1 sa 0.996" rule="forward"/>
    </rule>
    <rule name="turn" max_depth="90" successor="r">
        <call rule="dbox"/>
        <call transforms="rz 5.6 tx 0.1 sa 0.996" rule="turn"/>
    </rule>
    <rule name="turn2" max_depth="90" successor="r">
        <call rule="dbox"/>
        <call transforms="rz -5.6 tx 0.1 sa 0.996" rule="turn2"/>
    </rule>
    <rule name="turn3" max_depth="90" successor="r">
        <call rule="dbox"/>
        <call transforms="ry -5.6 tx 0.1 sa 0.996" rule="turn3"/>
    </rule>
    <rule name="turn4" max_depth="90" successor="r">
        <call rule="dbox"/>
        <call transforms="ry -5.6 tx 0.1 sa 0.996" rule="turn4"/>
    </rule>
    <rule name="turn5" max_depth="90" successor="r">
        <call rule="dbox"/>
        <call transforms="rx -5.6 tx 0.1 sa 0.996" rule="turn5"/>
    </rule>
    <rule name="turn6" max_depth="90" successor="r">
        <call rule="dbox"/>
        <call transforms="rx -5.6 tx 0.1 sa 0.996" rule="turn6"/>
    </rule>
    <rule name="dbox">
        <instance transforms="s 0.55 2.0 1.25" shape="curve"/>
    </rule>
</rules>
`
