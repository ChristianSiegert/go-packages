// Package html removes whitespace and comments from HTML code.
package html

import (
	"bytes"
	"html/template"
	"regexp"
	"strings"
)

var (
	regExpHtmlComment                   = regexp.MustCompile("<!--(.|[\r\n])*?-->")
	regExpLineBreak                     = regexp.MustCompile("[\r\n]+")
	regExpParagraphDelimiter            = regexp.MustCompile("(\r\n){2,}")
	regExpTag                           = regexp.MustCompile("<(.|[\r\n])*?>")
	regExpWhitespaceAtStart             = regexp.MustCompile("^[ \f\n\r\t\v]+")
	regExpWhitespaceAtEnd               = regexp.MustCompile("[ \f\n\r\t\v]+$")
	regExpWhitespaceBetweenTags         = regexp.MustCompile(">[ \f\n\r\t\v]+<")
	regExpWhitespaceBetweenActions      = regexp.MustCompile("}}[ \f\n\r\t\v]+{{")
	regExpWhitespaceBetweenTagAndAction = regexp.MustCompile(">[ \f\n\r\t\v]+{{")
	regExpWhitespaceBetweenActionAndTag = regexp.MustCompile("}}[ \f\n\r\t\v]+<")
	regExpWhitespaceInsideTagStart      = regexp.MustCompile("<(/)?[ \f\n\r\t\v]+")
	regExpWhitespaceInsideTagEnd        = regexp.MustCompile("[ \f\n\r\t\v]+(/?)>")
	regExpWhitespaceInsideTagEqualSign  = regexp.MustCompile("[ \f\n\r\t\v]*=[ \f\n\r\t\v]*")
	regExpWhitespaceInsideTag           = regexp.MustCompile("[ \f\n\r\t\v]{2,}")
)

// Paragraphs takes a plain text string, replaces single line breaks by <br>,
// and wraps <p></p> tags around text blocks that are separated by two or more
// line breaks.
func Paragraphs(s string) template.HTML {
	s = strings.TrimSpace(s)
	s = template.HTMLEscapeString(s)
	s = regExpParagraphDelimiter.ReplaceAllString(s, "</p><p>")
	s = strings.Replace(s, "\r\n", "<br>", -1)
	return template.HTML("<p>" + s + "</p>")
}

// RemoveComments removes HTML comments.
func RemoveComments(html []byte) []byte {
	return regExpHtmlComment.ReplaceAll(html, []byte(""))
}

// RemoveWhitespace removes whitespace between tags, actions, and at the
// beginning and end of the HTML code.
func RemoveWhitespace(html []byte) []byte {
	html = regExpWhitespaceBetweenTags.ReplaceAll(html, []byte("><"))
	html = regExpWhitespaceBetweenActions.ReplaceAll(html, []byte("}}{{"))
	html = regExpWhitespaceBetweenTagAndAction.ReplaceAll(html, []byte(">{{"))
	html = regExpWhitespaceBetweenActionAndTag.ReplaceAll(html, []byte("}}<"))
	html = regExpWhitespaceAtStart.ReplaceAll(html, []byte(""))
	html = regExpWhitespaceAtEnd.ReplaceAll(html, []byte(""))

	tags := regExpTag.FindAll(html, -1)

	for _, tag := range tags {
		cleanTag := regExpLineBreak.ReplaceAll(tag, []byte(" "))
		cleanTag = regExpWhitespaceInsideTagStart.ReplaceAll(cleanTag, []byte("<$1"))
		cleanTag = regExpWhitespaceInsideTagEnd.ReplaceAll(cleanTag, []byte("$1>"))
		cleanTag = regExpWhitespaceInsideTagEqualSign.ReplaceAllLiteral(cleanTag, []byte("="))
		cleanTag = regExpWhitespaceInsideTag.ReplaceAllLiteral(cleanTag, []byte(" "))
		html = bytes.Replace(html, tag, cleanTag, 1)
	}

	return html
}
