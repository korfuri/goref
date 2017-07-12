package main

import (
	"github.com/dustin/go-humanize"
	"github.com/korfuri/goref"
	"github.com/korfuri/goref/json"

	"log"
	"os"
	"runtime"
	"time"
)

func reportMemory() {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	log.Printf("Memory stats: Alloc:%s, TotalAlloc:%s, HeapAlloc:%s, HeapSys:%s\n",
		humanize.Bytes(mem.Alloc), humanize.Bytes(mem.TotalAlloc),
		humanize.Bytes(mem.HeapAlloc), humanize.Bytes(mem.HeapSys))
}

func main() {
	log.SetFlags(0)
	log.SetOutput(os.Stdout)

	reportMemory()

	start := time.Now()

	m := goref.NewPackageGraph(goref.FileMTimeVersion)
	m.LoadPrograms([]string{"github.com/korfuri/goref/main"}, true)

	log.Printf("Loading took %s\n", time.Since(start))
	reportMemory()
	loadingDone := time.Now()

	m.ComputeInterfaceImplementationMatrix()

	log.Printf("Type matrix took %s (total runtime: %s)\n", time.Since(loadingDone), time.Since(start))
	reportMemory()
	computeMatrixDone := time.Now()

	log.Printf("%d packages in the graph\n", len(m.Packages))

	log.Printf("Here are the uses of objects in `goref`:\n")
	for _, ref := range m.Packages["github.com/korfuri/goref"].InRefs {
		log.Printf("   - %s\n", ref)
	}

	log.Printf("Here is where `goref`.`InRefs` is used:\n")
	for _, ref := range m.Packages["github.com/korfuri/goref"].InRefs {
		if ref.ToIdent == "InRefs" {
			log.Printf("   - %s\n", ref)
		}
	}

	log.Printf("Here are the uses of objects in `log` by `main`:\n")
	for _, ref := range m.Packages["log"].InRefs {
		if ref.FromPackage == m.Packages["github.com/korfuri/goref/main"] {
			log.Printf("   - %s\n", ref)
		}
	}

	log.Printf("Who implements `log.Stringer`?\n")
	for _, ref := range m.Packages["fmt"].InRefs {
		if ref.ToIdent == "Stringer" && ref.RefType == goref.Implementation {
			log.Printf("   - %s\n", ref)
		}
	}

	log.Printf("Displaying took %s (total runtime: %s)\n", time.Since(computeMatrixDone), time.Since(start))

	jsonch := make(chan []byte)
	errch := make(chan error)
	done := make(chan struct{})
	go json.GraphAsJSON(*m, jsonch, errch, done)
	for {
		select {
		case j := <-jsonch:
			log.Printf("%s\n", string(j))
		case err := <-errch:
			log.Fatal(err)
		case <-done:
			return
		}
	}
}

func unused() interface{} {
	b := log.Logger{}
	b.Println("")
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
