package html

import (
	"html/template"
	"testing"
)

var html = []byte(`
<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8">
		<meta name="viewport" content="initial-scale=1, width=device-width">
		<title>Panoptikos</title>

		{{if .IsDevAppServer}}
			{{range .DevCssFiles}}
				<link href="{{.}}" rel="stylesheet" type="text/css">
			{{end}}
		{{else}}
			<link href="/{{.CompiledCssFile}}" rel="stylesheet" type="text/css">
		{{end}}
	</head>

	<body>
		<p id="some-class">Foo</p>
		<p id="some-other-class">Bar</p>

		{{if .IsDevAppServer}}
			{{range .DevJsFiles}}
				<script src="{{.}}"></script>
			{{end}}
		{{else}}
			<script src="/{{.CompiledJsFile}}"></script>
		{{end}}

		< div
			foo
			bar = "baz"
			baz1  baz2    baz3
		>
		</ div >

		< br >

		<!-- Comment 1 -->
		<script>var s = "Some JavaScript code"</script>

		<!-- Comment 2 -->
		<noscript>
			<div>Enable JavaScript.</div>
		</noscript>
	</body>
</html>
`)

func TestParagraphs(t *testing.T) {
	type Test struct {
		input    string
		expected template.HTML
	}

	tests := []*Test{
		// Single linebreak must result in <br>, multiple consecutive linebreaks in new paragraph
		{"Hello,\r\nworld!", template.HTML("<p>Hello,<br>world!</p>")},
		{"Hello,\r\n\r\nworld!", template.HTML("<p>Hello,</p><p>world!</p>")},
		{"Hello,\r\n\r\n\r\nworld!", template.HTML("<p>Hello,</p><p>world!</p>")},
		{"Hello,\r\n\r\n\r\n\r\nworld!", template.HTML("<p>Hello,</p><p>world!</p>")},
		// Leading and trailing linebreaks must be ignored
		{"\r\nHello,\r\nworld!\r\n", template.HTML("<p>Hello,<br>world!</p>")},
		{"\r\nHello,\r\n\r\nworld!\r\n", template.HTML("<p>Hello,</p><p>world!</p>")},
		{"\r\n\r\nHello,\r\nworld!\r\n\r\n", template.HTML("<p>Hello,<br>world!</p>")},
		{"\r\n\r\nHello,\r\n\r\nworld!\r\n\r\n", template.HTML("<p>Hello,</p><p>world!</p>")},
		// Already present HTML tags must be escaped
		{"<b>Hello,</b>\r\n\r\n\r\n<i>world!</i>", template.HTML("<p>&lt;b&gt;Hello,&lt;/b&gt;</p><p>&lt;i&gt;world!&lt;/i&gt;</p>")},
	}

	for _, test := range tests {
		if result := Paragraphs(test.input); result != test.expected {
			t.Errorf("f(%q) returned %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestRemoveComments(t *testing.T) {
	input := []byte("foo<!-- Comment with newline \n-->bar<!--\n Comment with two newlines \n-->baz")
	expectedOutput := []byte("foobarbaz")

	result := RemoveComments(input)

	if len(result) != len(expectedOutput) {
		t.Errorf("HTML comments weren’t removed correctly: '%s'", result)
		return
	}

	for i := range result {
		if result[i] != expectedOutput[i] {
			t.Errorf("HTML comments weren’t removed correctly: '%s'", result)
			return
		}
	}
}

func TestRemoveWhitespace(t *testing.T) {
	expectedResult := []byte(`<!DOCTYPE html><html><head><meta charset="utf-8"><meta name="viewport" content="initial-scale=1, width=device-width"><title>Panoptikos</title>{{if .IsDevAppServer}}{{range .DevCssFiles}}<link href="{{.}}" rel="stylesheet" type="text/css">{{end}}{{else}}<link href="/{{.CompiledCssFile}}" rel="stylesheet" type="text/css">{{end}}</head><body><p id="some-class">Foo</p><p id="some-other-class">Bar</p>{{if .IsDevAppServer}}{{range .DevJsFiles}}<script src="{{.}}"></script>{{end}}{{else}}<script src="/{{.CompiledJsFile}}"></script>{{end}}<div foo bar="baz" baz1 baz2 baz3></div><br><!-- Comment 1 --><script>var s = "Some JavaScript code"</script><!-- Comment 2 --><noscript><div>Enable JavaScript.</div></noscript></body></html>`)

	result := RemoveWhitespace(html)

	if len(result) != len(expectedResult) {
		t.Errorf("Whitespace wasn’t removed correctly: '%s'", result)
		return
	}

	for i := range result {
		if result[i] != expectedResult[i] {
			t.Errorf("Whitespace wasn’t removed correctly: '%s'", result)
			return
		}
	}
}

func BenchmarkRemoveWhitespace(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RemoveWhitespace(html)
	}
}
