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
	"bytes"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/russross/blackfriday"
)

const (
	DEFAULT_NUM_COLUMNS = 80
)

type list struct {
	index int
}

type header struct {
	text     []byte
	level    int
	children []*header
}

type vimDoc struct {
	rootHead *header
	lastHead *header
	title    string
	tocPos   int
	lists    []*list
}

func VimDocRenderer() blackfriday.Renderer {
	return &vimDoc{title: "test", tocPos: -1}
}

func (v *vimDoc) pushl() {
	v.lists = append(v.lists, &list{1})
}

func (v *vimDoc) popl() {
	if len(v.lists) == 0 {
		log.Fatal("invalid list operation")
	}

	v.lists = v.lists[:len(v.lists)-1]
}

func (v *vimDoc) getl() *list {
	if len(v.lists) == 0 {
		log.Fatal("invalid list operation")
	}

	return v.lists[len(v.lists)-1]
}

func (v *vimDoc) fixupCode(input []byte) []byte {
	r := regexp.MustCompile(`(?m)^\s*([<>])$`)
	return r.ReplaceAll(input, []byte("$1"))
}

func (v *vimDoc) fixupHeader(header []byte) []byte {
	return bytes.ToUpper(header)
}

func (v *vimDoc) buildTag(header []byte) []byte {
	return []byte(fmt.Sprintf("%s-%s", v.title, string(bytes.ToLower(header))))
}

func (*vimDoc) writeRule(out *bytes.Buffer, repeat string) {
	out.WriteString(strings.Repeat(repeat, DEFAULT_NUM_COLUMNS))
	out.WriteString("\n")
}

func (v *vimDoc) writeToc(out *bytes.Buffer, h *header, depth int) {
	out.WriteString(fmt.Sprintf(
		"%s%s: |%s|\n",
		strings.Repeat(" ", depth*blackfriday.TAB_SIZE_DEFAULT),
		string(h.text),
		v.buildTag(h.text)))
	for _, c := range h.children {
		v.writeToc(out, c, depth+1)
	}
}

func (v *vimDoc) format(out *bytes.Buffer, text string, trim int) {
	lines := strings.Split(text, "\n")

	for index, line := range lines {
		width := blackfriday.TAB_SIZE_DEFAULT
		if width >= trim && index == 0 {
			width -= trim
		}

		if len(line) == 0 {
			continue
		}

		out.WriteString(strings.Repeat(" ", width))
		out.WriteString(line)
		out.WriteString("\n")
	}
}

// Block-level callbacks
func (v *vimDoc) BlockCode(out *bytes.Buffer, text []byte, lang string) {
	out.WriteString(">\n")
	v.format(out, string(text), 0)
	out.WriteString("<\n\n")
}

func (v *vimDoc) BlockQuote(out *bytes.Buffer, text []byte) {
	out.WriteString(">\n")
	v.format(out, string(text), 0)
	out.WriteString("<\n\n")
}

func (v *vimDoc) BlockHtml(out *bytes.Buffer, text []byte) {
	out.WriteString(">\n")
	v.format(out, string(text), 0)
	out.WriteString("<\n\n")
}

func (v *vimDoc) Header(out *bytes.Buffer, text func() bool, level int, id string) {
	initPos := out.Len()

	switch level {
	case 1:
		v.writeRule(out, "=")
	case 2:
		v.writeRule(out, "-")
	}

	headerPos := out.Len()

	if !text() {
		out.Truncate(initPos)
		return
	}

	if v.tocPos == -1 && v.rootHead != nil {
		v.tocPos = initPos
	}

	var value []byte
	value = append(value, out.Bytes()[headerPos:]...)
	header := &header{value, level, nil}

	if v.lastHead == nil {
		if header.level != 1 {
			log.Println("warning: top-level header in document is not a level 1 header")
		}

		v.rootHead = header
		v.lastHead = header
	} else {
		if v.rootHead.level >= header.level {
			log.Println("warning: found header of higher or equal level to the root header")
		}

		if header.level <= v.lastHead.level {
			v.lastHead = header
		} else {
			v.lastHead.children = append(v.lastHead.children, header)
		}
	}

	out.Truncate(headerPos)
	out.WriteString(fmt.Sprintf("%s\t\t*%s*\n\n", v.fixupHeader(header.text), v.buildTag(header.text)))
}

func (v *vimDoc) HRule(out *bytes.Buffer) {
	v.writeRule(out, "-")
}

func (v *vimDoc) List(out *bytes.Buffer, text func() bool, flags int) {
	v.pushl()
	text()
	v.popl()
}

func (v *vimDoc) ListItem(out *bytes.Buffer, text []byte, flags int) {
	marker := out.Len()

	list := v.getl()
	if flags&blackfriday.LIST_TYPE_ORDERED == blackfriday.LIST_TYPE_ORDERED {
		out.WriteString(fmt.Sprintf("%d. ", list.index))
		list.index++
	} else {
		out.WriteString("* ")
	}

	v.format(out, string(text), out.Len()-marker)

	if flags&blackfriday.LIST_ITEM_END_OF_LIST == blackfriday.LIST_ITEM_END_OF_LIST {
		out.WriteString("\n")
	}
}

func (*vimDoc) Paragraph(out *bytes.Buffer, text func() bool) {
	marker := out.Len()
	if !text() {
		out.Truncate(marker)
		return
	}
	out.WriteString("\n\n")
}

func (*vimDoc) Table(out *bytes.Buffer, header []byte, body []byte, columnData []int) {
	// unimplemented
	log.Println("Table is a stub")
}

func (*vimDoc) TableRow(out *bytes.Buffer, text []byte) {
	// unimplemented
	log.Println("TableRow is a stub")
}

func (*vimDoc) TableHeaderCell(out *bytes.Buffer, text []byte, flags int) {
	// unimplemented
	log.Println("TableHeaderCell is a stub")
}

func (*vimDoc) TableCell(out *bytes.Buffer, text []byte, flags int) {
	// unimplemented
	log.Println("TableCell is a stub")
}

func (*vimDoc) Footnotes(out *bytes.Buffer, text func() bool) {
	// unimplemented
	log.Println("Footnotes is a stub")
}

func (*vimDoc) FootnoteItem(out *bytes.Buffer, name, text []byte, flags int) {
	// unimplemented
	log.Println("FootnoteItem is a stub")
}

func (*vimDoc) TitleBlock(out *bytes.Buffer, text []byte) {
	// unimplemented
	log.Println("TitleBlock is a stub")
}

// Span-level callbacks
func (*vimDoc) AutoLink(out *bytes.Buffer, link []byte, kind int) {
	out.Write(link)
}

func (*vimDoc) CodeSpan(out *bytes.Buffer, text []byte) {
	r := regexp.MustCompile(`\s`)

	// vim does not properly highlight spaces in code spans
	if !r.Match(text) {
		out.WriteString("`")
		out.Write(text)
		out.WriteString("`")
	}
}

func (*vimDoc) DoubleEmphasis(out *bytes.Buffer, text []byte) {
	out.Write(text)
}

func (*vimDoc) Emphasis(out *bytes.Buffer, text []byte) {
	out.Write(text)
}

func (*vimDoc) Image(out *bytes.Buffer, link []byte, title []byte, alt []byte) {
	// cannot view images in vim
}

func (*vimDoc) LineBreak(out *bytes.Buffer) {
	out.WriteString("\n")
}

func (*vimDoc) Link(out *bytes.Buffer, link []byte, title []byte, content []byte) {
	out.WriteString(fmt.Sprintf("%s (%s)", content, link))
}

func (*vimDoc) RawHtmlTag(out *bytes.Buffer, tag []byte) {
	// unimplemented
	log.Println("StrikeThrough is a stub")
}

func (*vimDoc) TripleEmphasis(out *bytes.Buffer, text []byte) {
	out.Write(text)
}

func (*vimDoc) StrikeThrough(out *bytes.Buffer, text []byte) {
	// unimplemented
	log.Println("StrikeThrough is a stub")
}

func (*vimDoc) FootnoteRef(out *bytes.Buffer, ref []byte, id int) {
	// unimplemented
	log.Println("FootnoteRef is a stub")
}

// Low-level callbacks
func (v *vimDoc) Entity(out *bytes.Buffer, entity []byte) {
	out.Write(entity)
}

func (v *vimDoc) NormalText(out *bytes.Buffer, text []byte) {
	out.Write(text)
}

// Header and footer
func (*vimDoc) DocumentHeader(out *bytes.Buffer) {
	// unimplemented
	log.Println("DocumentHeader is a stub")
}

func (v *vimDoc) DocumentFooter(out *bytes.Buffer) {
	var temp bytes.Buffer

	temp.Write(out.Bytes()[:v.tocPos])
	v.writeToc(&temp, v.rootHead, 0)
	temp.WriteString("\n")
	temp.Write(out.Bytes()[v.tocPos:])

	out.Reset()
	out.Write(v.fixupCode(temp.Bytes()))
}

func (*vimDoc) GetFlags() int {
	return 0
}
