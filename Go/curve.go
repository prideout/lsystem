package main

import (
    "bufio"
    "fmt"
    "os"
    "os/exec"
    //    ri "gorman"
    "./lsystem"
)

func main() {

    foo := lsystem.Evaluate("Ribbon.xml")
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
