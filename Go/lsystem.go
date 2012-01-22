package lsystem

import (
    "os"
    "fmt"
    "encoding/xml"
    )
    
type LSystem struct {
    MaxDepth int
    Rules []Rule `xml:"rule"`
}

type Rule struct {
    Name string  `xml:"name,attr"`
    Calls []Call `xml:"call"`
    Instances []Instance `xml:"instance"`
    MaxDepth int `xml:"max_depth,attr"`
    Successor string `xml:"successor,attr"`
}

type Call struct {
    Transforms string `xml:"transforms,attr"`
    Rule string  `xml:"rule,attr"`
    Count int   `xml:"count,attr"`
}

type Instance struct {
    Transforms string `xml:"transforms"`
    Shape string  `xml:"string"`
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
	
	for r, rule := range lsys.Rules {
        fmt.Printf("Rule %d. %s -> %s\n", r, rule.Name, rule.Successor)
	    for i, call := range rule.Calls {
	        fmt.Printf("%d. count=%d\ttransforms=%s\trule=%s\n", 
	            i, call.Count, call.Transforms, call.Rule)
	    }
	}
	
    return len(lsys.Rules)
    
}