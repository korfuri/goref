package main

import (
	"github.com/dustin/go-humanize"
	"github.com/korfuri/goref"

	"log"
	"os"
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
	log.SetFlags(0)
	log.SetOutput(os.Stdout)

	reportMemory()

	start := time.Now()

	m := goref.NewPackageGraph()
	m.LoadProgram("github.com/korfuri/goref/main", []string{"main.go"})

	log.Printf("Loading took %s\n", time.Since(start))
	reportMemory()
	loadingDone := time.Now()

	m.ComputeInterfaceImplementationMatrix()

	log.Printf("Type matrix took %s (total runtime: %s)\n", time.Since(loadingDone), time.Since(start))
	reportMemory()
	computeMatrixDone := time.Now()

	log.Printf("%d packages in the graph\n", len(m.Packages))
	log.Printf("%d files in the graph\n", len(m.Files))

	log.Printf("Packages that depend on `fmt`:\n")
	for d, _ := range m.Packages["fmt"].Dependents {
		log.Printf("   - %s\n", d)
	}

	log.Printf("Packages that `goref` depends on:\n")
	for d, _ := range m.Packages["github.com/korfuri/goref"].Dependencies {
		log.Printf("   - %s\n", d)
	}

	log.Printf("Package `goref` has these files:\n")
	for d, _ := range m.Packages["github.com/korfuri/goref"].Files {
		log.Printf("   - %s\n", d)
	}

	log.Printf("Package `fmt` has these files:\n")
	for d, _ := range m.Packages["fmt"].Files {
		log.Printf("   - %s\n", d)
	}

	log.Printf("Here are the uses of objects in `goref`:\n")
	for pos, ref := range m.Packages["github.com/korfuri/goref"].InRefs {
		log.Printf("   - %s %s\n", pos, ref)
	}

	log.Printf("Here is where `goref`.`InRefs` is used:\n")
	for pos, ref := range m.Packages["github.com/korfuri/goref"].InRefs {
		if ref.Ident == "InRefs" {
			log.Printf("   - %s\n", pos)
		}
	}

	log.Printf("Here are the uses of objects in `log` by `main`:\n")
	for pos, ref := range m.Packages["log"].InRefs {
		if ref.FromPackage == m.Packages["github.com/korfuri/goref/main"] {
			log.Printf("   - %s %s\n", pos, ref)
		}
	}

	log.Printf("Who implements `log.Stringer`?\n")
	for pos, ref := range m.Packages["fmt"].InRefs {
		if ref.Ident == "Stringer" && ref.RefType == goref.Implementation {
			log.Printf("   - implemented at %s by %s\n", pos, ref)
		}
	}

	log.Printf("Displaying took %s (total runtime: %s)\n", time.Since(computeMatrixDone), time.Since(start))
}

func unused() interface{} {
	b := log.Logger{}
	log.Print(b)
	return log.Fatalf
}

type UnusedI interface {
	blah() string
}

type UnusedT int

func (u UnusedT) blah() string {
	return ""
}

type EmptyI interface{}