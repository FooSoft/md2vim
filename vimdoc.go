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
	"log"

	"github.com/russross/blackfriday"
)

type vimDoc struct {
}

func VimDocRenderer() blackfriday.Renderer {
	return &vimDoc{}
}

// Block-level callbacks
func (*vimDoc) BlockCode(out *bytes.Buffer, text []byte, lang string) {
	log.Println("stubbing BlockCode")
}

func (*vimDoc) BlockQuote(out *bytes.Buffer, text []byte) {
	log.Println("stubbing BlockQuote")
}

func (*vimDoc) BlockHtml(out *bytes.Buffer, text []byte) {
	log.Println("stubbing BlockHtml")
}

func (*vimDoc) Header(out *bytes.Buffer, text func() bool, level int, id string) {
	log.Println("stubbing Header")
}

func (*vimDoc) HRule(out *bytes.Buffer) {
	log.Println("stubbing HRule")
}

func (*vimDoc) List(out *bytes.Buffer, text func() bool, flags int) {
	log.Println("stubbing List")
}

func (*vimDoc) ListItem(out *bytes.Buffer, text []byte, flags int) {
	log.Println("stubbing ListItem")
}

func (*vimDoc) Paragraph(out *bytes.Buffer, text func() bool) {
	log.Println("stubbing Paragraph")
}

func (*vimDoc) Table(out *bytes.Buffer, header []byte, body []byte, columnData []int) {
	log.Println("stubbing Table")
}

func (*vimDoc) TableRow(out *bytes.Buffer, text []byte) {
	log.Println("stubbing TableRow")
}

func (*vimDoc) TableHeaderCell(out *bytes.Buffer, text []byte, flags int) {
	log.Println("stubbing TableHeaderCell")
}

func (*vimDoc) TableCell(out *bytes.Buffer, text []byte, flags int) {
	log.Println("stubbing TableCell")
}

func (*vimDoc) Footnotes(out *bytes.Buffer, text func() bool) {
	log.Println("stubbing Footnotes")
}

func (*vimDoc) FootnoteItem(out *bytes.Buffer, name, text []byte, flags int) {
	log.Println("stubbing FootnoteItem")
}

func (*vimDoc) TitleBlock(out *bytes.Buffer, text []byte) {
	log.Println("stubbing TitleBlock")
}

// Span-level callbacks
func (*vimDoc) AutoLink(out *bytes.Buffer, link []byte, kind int) {
	log.Println("stubbing AutoLink")
}

func (*vimDoc) CodeSpan(out *bytes.Buffer, text []byte) {
	log.Println("stubbing CodeSpan")
}

func (*vimDoc) DoubleEmphasis(out *bytes.Buffer, text []byte) {
	log.Println("stubbing DoubleEmphasis")
}

func (*vimDoc) Emphasis(out *bytes.Buffer, text []byte) {
	log.Println("stubbing Emphasis")
}

func (*vimDoc) Image(out *bytes.Buffer, link []byte, title []byte, alt []byte) {
	log.Println("stubbing Image")
}

func (*vimDoc) LineBreak(out *bytes.Buffer) {
	log.Println("stubbing LineBreak")
}

func (*vimDoc) Link(out *bytes.Buffer, link []byte, title []byte, content []byte) {
	log.Println("stubbing Link")
}

func (*vimDoc) RawHtmlTag(out *bytes.Buffer, tag []byte) {
	log.Println("stubbing RawHtmlTag")
}

func (*vimDoc) TripleEmphasis(out *bytes.Buffer, text []byte) {
	log.Println("stubbing TripleEmphasis")
}

func (*vimDoc) StrikeThrough(out *bytes.Buffer, text []byte) {
	log.Println("stubbing StrikeThrough")
}

func (*vimDoc) FootnoteRef(out *bytes.Buffer, ref []byte, id int) {
	log.Println("stubbing FootnoteRef")
}

// Low-level callbacks
func (*vimDoc) Entity(out *bytes.Buffer, entity []byte) {
	log.Println("stubbing Entity")
}

func (*vimDoc) NormalText(out *bytes.Buffer, text []byte) {
	log.Println("stubbing NormalText")
}

// Header and footer
func (*vimDoc) DocumentHeader(out *bytes.Buffer) {
	log.Println("stubbing DocumentHeader")
}

func (*vimDoc) DocumentFooter(out *bytes.Buffer) {
	log.Println("stubbing DocumentFooter")
}

func (*vimDoc) GetFlags() int {
	return 0
}
