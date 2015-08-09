/*
 * Copyright (c) 2015 Alex Yatskov <alex@foosoft.net>
 * Author: Alex Yatskov <alex@foosoft.net>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of
 * this software and associated documentation files (the "Software"), to deal in
 * the Software without restriction, including without limitation the rights to
 * use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
 * the Software, and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
 * FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
 * IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/russross/blackfriday"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options] input output\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "http://foosoft.net/projects/vimdown/\n\n")
	fmt.Fprintf(os.Stderr, "Parameters:\n")
	flag.PrintDefaults()
}

func main() {
	cols := flag.Int("cols", DEFAULT_NUM_COLUMNS, "number of columns to use for alignment and rules")
	tabs := flag.Int("tabs", DEFAULT_TAB_SIZE, "size of the tab to use specified as number of spaces")
	notoc := flag.Bool("notoc", false, "do not generate table of contents for headings")
	norules := flag.Bool("norules", false, "do not generate horizontal rules above headings")
	pascal := flag.Bool("pascal", false, "use PascalCase for abbreviating tags")
	desc := flag.String("desc", "", "short description of extension")
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		flag.Usage()
		os.Exit(-1)
	}

	input, err := ioutil.ReadFile(args[0])
	if err != nil {
		log.Fatalf("error: unable to read from file %s", args[0])
	}

	flags := 0
	if *notoc {
		flags |= FLAG_NO_TOC
	}
	if *norules {
		flags |= FLAG_NO_RULES
	}
	if *pascal {
		flags |= FLAG_PASCAL
	}

	renderer := VimDocRenderer(args[1], *desc, *cols, *tabs, flags)
	extensions := blackfriday.EXTENSION_FENCED_CODE | blackfriday.EXTENSION_NO_INTRA_EMPHASIS | blackfriday.EXTENSION_SPACE_HEADERS
	output := blackfriday.Markdown(input, renderer, extensions)

	file, err := os.Create(args[1])
	if err != nil {
		log.Fatalf("error: unable to write to file %s", args[1])
	}
	defer file.Close()

	if _, err := file.Write(output); err != nil {
		log.Fatal("error: unable to write output")
	}
}
