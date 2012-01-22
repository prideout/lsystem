package lsystem

import (
    "os"
    "fmt"
    "xml"
    )
    
type LSystem struct {
    MaxDepth int
    Rules []Rule `xml:"rule>"`
}

type Rule struct {
    Calls []Call `xml:"call>"`
    Instances []Instance `xml:"instance>"`
    MaxDepth int `xml:"attr>max_depth"`
    Successor string `xml:"attr>successor"`
}

type Call struct {
    Transforms string `xml:"attr>transforms"`
    Rule string  `xml:"attr>rule"`
}

type Instance struct {
    Transforms string `xml:"transforms>"`
    Shape string  `xml:"string>"`
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
	
    return len(lsys.Rules)
    
}