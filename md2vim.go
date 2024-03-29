package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/russross/blackfriday"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options] input output\n", path.Base(os.Args[0]))
	fmt.Fprintf(os.Stderr, "https://foosoft.net/projects/md2vim/\n\n")
	fmt.Fprintf(os.Stderr, "Parameters:\n")
	flag.PrintDefaults()
}

func main() {
	cols := flag.Int("cols", defNumCols, "number of columns to use for layout")
	tabs := flag.Int("tabs", defTabSize, "tab width specified in number of spaces")
	notoc := flag.Bool("notoc", false, "do not generate table of contents for headings")
	norules := flag.Bool("norules", false, "do not generate horizontal rules above headings")
	pascal := flag.Bool("pascal", false, "use PascalCase for abbreviating tags")
	desc := flag.String("desc", "", "short description of the help file")
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		flag.Usage()
		os.Exit(2)
	}

	input, err := ioutil.ReadFile(args[0])
	if err != nil {
		log.Fatalf("unable to read from file: %s", args[0])
	}

	flags := 0
	if *notoc {
		flags |= flagNoToc
	}
	if *norules {
		flags |= flagNoRules
	}
	if *pascal {
		flags |= flagPascal
	}

	renderer := VimDocRenderer(args[1], *desc, *cols, *tabs, flags)
	extensions := blackfriday.EXTENSION_FENCED_CODE | blackfriday.EXTENSION_NO_INTRA_EMPHASIS | blackfriday.EXTENSION_SPACE_HEADERS
	output := blackfriday.Markdown(input, renderer, extensions)

	file, err := os.Create(args[1])
	if err != nil {
		log.Fatalf("unable to write to file: %s", args[1])
	}
	defer file.Close()

	if _, err := file.Write(output); err != nil {
		log.Fatal("unable to write output")
	}
}
