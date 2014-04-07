package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-hep/rio"
)

func main() {
	var fname string

	flag.Parse()

	if flag.NArg() > 0 {
		fname = flag.Arg(0)
	}

	fmt.Printf("::: inspecting file [%s]...\n", fname)
	if fname == "" {
		flag.Usage()
		os.Exit(1)
	}

	f, err := rio.Open(fname)
	if err != nil {
		fmt.Printf("*** error: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	_, _ = f.ReadRecord()

	_, err = f.Seek(0, 0)
	if err != nil {
		fmt.Printf("*** error: %v\n", err)
	}

	for _, rec := range f.Records() {
		fmt.Printf(" -> %v\n", rec.Name())
	}

	fmt.Printf("::: inspecting file [%s]... [done]\n", fname)
}
