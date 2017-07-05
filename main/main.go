package main

import (
	"github.com/korfuri/goref"

	"fmt"
	"log"
)

func main() {
	m := goref.NewPackageGraph()
	m.LoadProgram("github.com/korfuri/codesearch", "main.go")

	fmt.Printf("%d packages in the graph\n", len(m.Packages))
	fmt.Printf("%d files in the graph\n", len(m.Files))

	fmt.Printf("Packages that depend on `fmt`:\n")
	for d, _ := range m.Packages["fmt"].Dependents {
		fmt.Printf("   - %s\n", d)
	}
	
	fmt.Printf("Packages that `goref` depends on:\n")
	for d, _ := range m.Packages["github.com/korfuri/goref"].Dependencies {
		fmt.Printf("   - %s\n", d)
	}

	fmt.Printf("Package `goref` has these files:\n")
	for d, _ := range m.Packages["github.com/korfuri/goref"].Files {
		fmt.Printf("   - %s\n", d)
	}

	fmt.Printf("Package `fmt` has these files:\n")
	for d, _ := range m.Packages["fmt"].Files {
		fmt.Printf("   - %s\n", d)
	}

	fmt.Printf("Here are the uses of objects in `log`:\n")
	for pos, ref := range m.Packages["log"].InRefs {
		fmt.Printf("   - %s %s\n", pos, ref)
	}
}

func unused() interface{} {
	b := log.Logger{}
	log.Print(b)
	return log.Fatalf
}
