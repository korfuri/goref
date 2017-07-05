package main

import (
	"github.com/dustin/go-humanize"
	"github.com/korfuri/goref"

	"fmt"
	"log"
	"runtime"
	"time"
)

func reportMemory() {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	log.Printf("Allocated memory: %s\n", humanize.Bytes(mem.Alloc))
	log.Printf("Total allocated memory: %s\n", humanize.Bytes(mem.TotalAlloc))
	log.Printf("Heap allocated memory: %s\n", humanize.Bytes(mem.HeapAlloc))
	log.Printf("System heap allocated memory: %s\n", humanize.Bytes(mem.HeapSys))
}

func main() {
	reportMemory()

	start := time.Now()

	m := goref.NewPackageGraph()
	m.LoadProgram("github.com/korfuri/goref/main", "main.go")

	fmt.Printf("Loading took %s\n", time.Since(start))

	reportMemory()

	fmt.Printf("%d packages in the graph\n", len(m.Packages))
	fmt.Printf("%d files in the graph\n", len(m.Files))

	loadingDone := time.Now()

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

	fmt.Printf("Here are the uses of objects in `goref`:\n")
	for pos, ref := range m.Packages["github.com/korfuri/goref"].InRefs {
		fmt.Printf("   - %s %s\n", pos, ref)
	}

	fmt.Printf("Here is where `goref`.`InRefs` is used:\n")
	for pos, ref := range m.Packages["github.com/korfuri/goref"].InRefs {
		if ref.Ident == "InRefs" {
			fmt.Printf("   - %s\n", pos)
		}
	}

	fmt.Printf("Displaying took %s (total runtime: %s)\n", time.Since(loadingDone), time.Since(start))
}

func unused() interface{} {
	b := log.Logger{}
	log.Print(b)
	return log.Fatalf
}
