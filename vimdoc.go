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
	"path"
	"regexp"
	"strings"

	"github.com/russross/blackfriday"
)

const (
	defNumCols = 80
	defTabSize = 4
)

const (
	flagNoToc = 1 << iota
	flagNoRules
	flagPascal
)

type list struct {
	index int
}

type heading struct {
	text     []byte
	level    int
	children []*heading
}

type vimDoc struct {
	filename string
	title    string
	desc     string
	cols     int
	tabs     int
	flags    int
	tocPos   int
	lists    []*list
	rootHead *heading
	lastHead *heading
}

func VimDocRenderer(filename, desc string, cols, tabs, flags int) blackfriday.Renderer {
	filename = path.Base(filename)
	title := filename

	if index := strings.LastIndex(filename, "."); index > -1 {
		title = filename[:index]
		if flags&flagPascal == 0 {
			title = strings.ToLower(title)
		}
	}

	return &vimDoc{
		filename: filename,
		title:    title,
		desc:     desc,
		cols:     cols,
		tabs:     tabs,
		flags:    flags,
		tocPos:   -1}
}

func (v *vimDoc) pushl() {
	v.lists = append(v.lists, &list{1})
}

func (v *vimDoc) popl() {
	v.lists = v.lists[:len(v.lists)-1]
}

func (v *vimDoc) getl() *list {
	return v.lists[len(v.lists)-1]
}

func (v *vimDoc) fixupCode(input []byte) []byte {
	r := regexp.MustCompile(`(?m)^\s*([<>])$`)
	return r.ReplaceAll(input, []byte("$1"))
}

func (v *vimDoc) fixupHeader(text []byte) []byte {
	return bytes.ToUpper(text)
}

func (v *vimDoc) buildTag(text []byte) []byte {
	if v.flags&flagPascal == 0 {
		text = bytes.ToLower(text)
		text = bytes.Replace(text, []byte{' '}, []byte{'_'}, -1)
	} else {
		text = bytes.Title(text)
		text = bytes.Replace(text, []byte{' '}, []byte{}, -1)
	}

	return []byte(fmt.Sprintf("%s-%s", v.title, text))
}

func (v *vimDoc) writeStraddle(out *bytes.Buffer, left, right []byte, trim int) {
	padding := v.cols - (len(left) + len(right)) + trim
	if padding <= 0 {
		padding = 1
	}

	out.Write(left)
	out.WriteString(strings.Repeat(" ", padding))
	out.Write(right)
	out.WriteString("\n")
}

func (v *vimDoc) writeRule(out *bytes.Buffer, repeat string) {
	out.WriteString(strings.Repeat(repeat, v.cols))
	out.WriteString("\n")
}

func (v *vimDoc) writeToc(out *bytes.Buffer, head *heading, depth int) {
	title := fmt.Sprintf("%s%s", strings.Repeat(" ", depth*v.tabs), head.text)
	link := fmt.Sprintf("|%s|", v.buildTag(head.text))
	v.writeStraddle(out, []byte(title), []byte(link), 2)

	for _, child := range head.children {
		v.writeToc(out, child, depth+1)
	}
}

func (v *vimDoc) format(out *bytes.Buffer, text string, trim int) {
	lines := strings.Split(text, "\n")

	for index, line := range lines {
		width := v.tabs
		if width >= trim && index == 0 {
			width -= trim
		}

		if len(line) > 0 {
			out.WriteString(strings.Repeat(" ", width))
			out.WriteString(line)
			out.WriteString("\n")
		}
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

	if v.flags&flagNoRules == 0 {
		switch level {
		case 1:
			v.writeRule(out, "=")
		case 2:
			v.writeRule(out, "-")
		}
	}

	headingPos := out.Len()

	if !text() {
		out.Truncate(initPos)
		return
	}

	if v.tocPos == -1 && v.rootHead != nil {
		v.tocPos = initPos
	}

	var value []byte
	value = append(value, out.Bytes()[headingPos:]...)
	heading := &heading{value, level, nil}

	if v.lastHead == nil {
		if heading.level != 1 {
			log.Println("warning: top-level heading in document is not a level 1 heading")
		}

		v.rootHead = heading
		v.lastHead = heading
	} else {
		if v.rootHead.level >= heading.level {
			log.Println("warning: found heading of higher or equal level to the root heading")
		}

		if heading.level <= v.lastHead.level {
			v.lastHead = heading
		} else {
			v.lastHead.children = append(v.lastHead.children, heading)
		}
	}

	out.Truncate(headingPos)
	tag := fmt.Sprintf("*%s*", v.buildTag(heading.text))
	v.writeStraddle(out, v.fixupHeader(heading.text), []byte(tag), 2)
	out.WriteString("\n")
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

	if flags&blackfriday.LIST_ITEM_END_OF_LIST != 0 {
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

func (*vimDoc) Table(out *bytes.Buffer, heading []byte, body []byte, columnData []int) {
	// unimplemented
	log.Println("warning: Table is a stub")
}

func (*vimDoc) TableRow(out *bytes.Buffer, text []byte) {
	// unimplemented
	log.Println("warning: TableRow is a stub")
}

func (*vimDoc) TableHeaderCell(out *bytes.Buffer, text []byte, flags int) {
	// unimplemented
	log.Println("warning: TableHeaderCell is a stub")
}

func (*vimDoc) TableCell(out *bytes.Buffer, text []byte, flags int) {
	// unimplemented
	log.Println("warning: TableCell is a stub")
}

func (*vimDoc) Footnotes(out *bytes.Buffer, text func() bool) {
	// unimplemented
	log.Println("warning: Footnotes is a stub")
}

func (*vimDoc) FootnoteItem(out *bytes.Buffer, name, text []byte, flags int) {
	// unimplemented
	log.Println("warning: FootnoteItem is a stub")
}

func (*vimDoc) TitleBlock(out *bytes.Buffer, text []byte) {
	// unimplemented
	log.Println("warning: TitleBlock is a stub")
}

// Span-level callbacks
func (*vimDoc) AutoLink(out *bytes.Buffer, link []byte, kind int) {
	out.Write(link)
}

func (*vimDoc) CodeSpan(out *bytes.Buffer, text []byte) {
	r := regexp.MustCompile(`\s`)

	// vim does not correctly highlight space-delimited words in code spans
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
	log.Println("warning: StrikeThrough is a stub")
}

func (*vimDoc) TripleEmphasis(out *bytes.Buffer, text []byte) {
	out.Write(text)
}

func (*vimDoc) StrikeThrough(out *bytes.Buffer, text []byte) {
	// unimplemented
	log.Println("warning: StrikeThrough is a stub")
}

func (*vimDoc) FootnoteRef(out *bytes.Buffer, ref []byte, id int) {
	// unimplemented
	log.Println("warning: FootnoteRef is a stub")
}

// Low-level callbacks
func (v *vimDoc) Entity(out *bytes.Buffer, entity []byte) {
	out.Write(entity)
}

func (v *vimDoc) NormalText(out *bytes.Buffer, text []byte) {
	out.Write(text)
}

// Header and footer
func (v *vimDoc) DocumentHeader(out *bytes.Buffer) {
	if len(v.desc) > 0 {
		v.writeStraddle(out, []byte(v.filename), []byte(v.desc), 0)
	} else {
		out.WriteString(v.filename)
		out.WriteString("\n")
	}

	out.WriteString("\n")
}

func (v *vimDoc) DocumentFooter(out *bytes.Buffer) {
	var temp bytes.Buffer

	if v.tocPos > 0 && v.flags&flagNoToc == 0 {
		temp.Write(out.Bytes()[:v.tocPos])
		v.writeToc(&temp, v.rootHead, 0)
		temp.WriteString("\n")
		temp.Write(out.Bytes()[v.tocPos:])
	} else {
		temp.ReadFrom(out)
	}

	out.Reset()
	out.Write(v.fixupCode(temp.Bytes()))
}

func (v *vimDoc) GetFlags() int {
	return v.flags
}
