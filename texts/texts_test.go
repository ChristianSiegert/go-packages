package texts

import "testing"

func TestTruncate(t *testing.T) {
	type Test struct {
		text      string
		maxLength int
		suffix    string
		exact     bool
		expected  string
	}

	tests := []*Test{
		// maxLength changes, suffix is “…”, exact is false
		{
			text:      "Lorem ipsum",
			maxLength: -1,
			suffix:    "…",
			exact:     false,
			expected:  "",
		},
		{
			text:      "Lorem ipsum",
			maxLength: 0,
			suffix:    "…",
			exact:     false,
			expected:  "",
		},
		{
			text:      "Lorem ipsum",
			maxLength: 1,
			suffix:    "…",
			exact:     false,
			expected:  "…",
		},
		{
			text:      "Lorem ipsum",
			maxLength: 2,
			suffix:    "…",
			exact:     false,
			expected:  "…",
		},
		{
			text:      "Lorem ipsum",
			maxLength: 6,
			suffix:    "…",
			exact:     false,
			expected:  "…",
		},
		{
			text:      "Lorem ipsum",
			maxLength: 7,
			suffix:    "…",
			exact:     false,
			expected:  "Lorem …",
		},
		{
			text:      "Lorem ipsum",
			maxLength: 8,
			suffix:    "…",
			exact:     false,
			expected:  "Lorem …",
		},
		{
			text:      "Lorem ipsum",
			maxLength: 9,
			suffix:    "…",
			exact:     false,
			expected:  "Lorem …",
		},
		{
			text:      "Lorem ipsum",
			maxLength: 10,
			suffix:    "…",
			exact:     false,
			expected:  "Lorem …",
		},
		{
			text:      "Lorem ipsum",
			maxLength: 11,
			suffix:    "…",
			exact:     false,
			expected:  "Lorem ipsum",
		},
		{
			text:      "Lorem ipsum",
			maxLength: 12,
			suffix:    "…",
			exact:     false,
			expected:  "Lorem ipsum",
		},

		// maxLength changes, suffix is “…”, exact is true
		{
			text:      "Lorem ipsum",
			maxLength: -1,
			suffix:    "…",
			exact:     true,
			expected:  "",
		},
		{
			text:      "Lorem ipsum",
			maxLength: 0,
			suffix:    "…",
			exact:     true,
			expected:  "",
		},
		{
			text:      "Lorem ipsum",
			maxLength: 1,
			suffix:    "…",
			exact:     true,
			expected:  "…",
		},
		{
			text:      "Lorem ipsum",
			maxLength: 2,
			suffix:    "…",
			exact:     true,
			expected:  "L…",
		},
		{
			text:      "Lorem ipsum",
			maxLength: 3,
			suffix:    "…",
			exact:     true,
			expected:  "Lo…",
		},
		{
			text:      "Lorem ipsum",
			maxLength: 4,
			suffix:    "…",
			exact:     true,
			expected:  "Lor…",
		},
		{
			text:      "Lorem ipsum",
			maxLength: 5,
			suffix:    "…",
			exact:     true,
			expected:  "Lore…",
		},
		{
			text:      "Lorem ipsum",
			maxLength: 6,
			suffix:    "…",
			exact:     true,
			expected:  "Lorem…",
		},
		{
			text:      "Lorem ipsum",
			maxLength: 7,
			suffix:    "…",
			exact:     true,
			expected:  "Lorem …",
		},
		{
			text:      "Lorem ipsum",
			maxLength: 8,
			suffix:    "…",
			exact:     true,
			expected:  "Lorem i…",
		},
		{
			text:      "Lorem ipsum",
			maxLength: 9,
			suffix:    "…",
			exact:     true,
			expected:  "Lorem ip…",
		},
		{
			text:      "Lorem ipsum",
			maxLength: 10,
			suffix:    "…",
			exact:     true,
			expected:  "Lorem ips…",
		},
		{
			text:      "Lorem ipsum",
			maxLength: 11,
			suffix:    "…",
			exact:     true,
			expected:  "Lorem ipsum",
		},
		{
			text:      "Lorem ipsum",
			maxLength: 12,
			suffix:    "…",
			exact:     true,
			expected:  "Lorem ipsum",
		},

		// maxLength changes, suffix is “...”, exact is false
		{
			text:      "Lorem ipsum",
			maxLength: -1,
			suffix:    "...",
			exact:     false,
			expected:  "",
		},
		{
			text:      "Lorem ipsum",
			maxLength: 0,
			suffix:    "...",
			exact:     false,
			expected:  "",
		},
		{
			text:      "Lorem ipsum",
			maxLength: 1,
			suffix:    "...",
			exact:     false,
			expected:  "",
		},
		{
			text:      "Lorem ipsum",
			maxLength: 2,
			suffix:    "...",
			exact:     false,
			expected:  "",
		},
		{
			text:      "Lorem ipsum",
			maxLength: 3,
			suffix:    "...",
			exact:     false,
			expected:  "...",
		},
		{
			text:      "Lorem ipsum",
			maxLength: 6,
			suffix:    "...",
			exact:     false,
			expected:  "...",
		},
		{
			text:      "Lorem ipsum",
			maxLength: 7,
			suffix:    "...",
			exact:     false,
			expected:  "...",
		},
		{
			text:      "Lorem ipsum",
			maxLength: 8,
			suffix:    "...",
			exact:     false,
			expected:  "...",
		},
		{
			text:      "Lorem ipsum",
			maxLength: 9,
			suffix:    "...",
			exact:     false,
			expected:  "Lorem ...",
		},
		{
			text:      "Lorem ipsum",
			maxLength: 10,
			suffix:    "...",
			exact:     false,
			expected:  "Lorem ...",
		},
		{
			text:      "Lorem ipsum",
			maxLength: 11,
			suffix:    "...",
			exact:     false,
			expected:  "Lorem ipsum",
		},
		{
			text:      "Lorem ipsum",
			maxLength: 12,
			suffix:    "...",
			exact:     false,
			expected:  "Lorem ipsum",
		},

		// maxLength changes, suffix is “...”, exact is true
		{
			text:      "Lorem ipsum",
			maxLength: -1,
			suffix:    "...",
			exact:     true,
			expected:  "",
		},
		{
			text:      "Lorem ipsum",
			maxLength: 0,
			suffix:    "...",
			exact:     true,
			expected:  "",
		},
		{
			text:      "Lorem ipsum",
			maxLength: 1,
			suffix:    "...",
			exact:     true,
			expected:  "",
		},
		{
			text:      "Lorem ipsum",
			maxLength: 2,
			suffix:    "...",
			exact:     true,
			expected:  "",
		},
		{
			text:      "Lorem ipsum",
			maxLength: 3,
			suffix:    "...",
			exact:     true,
			expected:  "...",
		},
		{
			text:      "Lorem ipsum",
			maxLength: 4,
			suffix:    "...",
			exact:     true,
			expected:  "L...",
		},
		{
			text:      "Lorem ipsum",
			maxLength: 5,
			suffix:    "...",
			exact:     true,
			expected:  "Lo...",
		},
		{
			text:      "Lorem ipsum",
			maxLength: 6,
			suffix:    "...",
			exact:     true,
			expected:  "Lor...",
		},
		{
			text:      "Lorem ipsum",
			maxLength: 7,
			suffix:    "...",
			exact:     true,
			expected:  "Lore...",
		},
		{
			text:      "Lorem ipsum",
			maxLength: 8,
			suffix:    "...",
			exact:     true,
			expected:  "Lorem...",
		},
		{
			text:      "Lorem ipsum",
			maxLength: 9,
			suffix:    "...",
			exact:     true,
			expected:  "Lorem ...",
		},
		{
			text:      "Lorem ipsum",
			maxLength: 10,
			suffix:    "...",
			exact:     true,
			expected:  "Lorem i...",
		},
		{
			text:      "Lorem ipsum",
			maxLength: 11,
			suffix:    "...",
			exact:     true,
			expected:  "Lorem ipsum",
		},
		{
			text:      "Lorem ipsum",
			maxLength: 12,
			suffix:    "...",
			exact:     true,
			expected:  "Lorem ipsum",
		},

		// If suffix is longer than text but text doesn’t need to be truncated,
		// the text must be returned instead of an empty string.
		{
			text:      "Lo",
			maxLength: 2,
			suffix:    "...",
			exact:     false,
			expected:  "Lo",
		},
		{
			text:      "Lo",
			maxLength: 2,
			suffix:    "...",
			exact:     true,
			expected:  "Lo",
		},

		// No suffix
		{
			text:      "Lorem ipsum",
			maxLength: 8,
			suffix:    "",
			exact:     false,
			expected:  "Lorem ",
		},
		{
			text:      "Lorem ipsum",
			maxLength: 8,
			suffix:    "",
			exact:     true,
			expected:  "Lorem ip",
		},

		// A more real-world-like test
		{
			text:      "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
			maxLength: 50,
			suffix:    "…",
			exact:     false,
			expected:  "Lorem ipsum dolor sit amet, consectetur …",
		},
		{
			text:      "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
			maxLength: 50,
			suffix:    "…",
			exact:     true,
			expected:  "Lorem ipsum dolor sit amet, consectetur adipiscin…",
		},
		{
			text:      "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
			maxLength: 50,
			suffix:    "...",
			exact:     false,
			expected:  "Lorem ipsum dolor sit amet, consectetur ...",
		},
		{
			text:      "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
			maxLength: 50,
			suffix:    "...",
			exact:     true,
			expected:  "Lorem ipsum dolor sit amet, consectetur adipisc...",
		},

		{
			text:      "Lorem ipsum dolor sit amet,\r\nconsectetur adipiscing elit.",
			maxLength: 40,
			suffix:    "…",
			exact:     false,
			expected:  "Lorem ipsum dolor sit amet,\r\n…",
		},
		{
			text:      "Lorem ipsum dolor sit amet,\r\nconsectetur adipiscing elit.",
			maxLength: 40,
			suffix:    "…",
			exact:     true,
			expected:  "Lorem ipsum dolor sit amet,\r\nconsectetu…",
		},
		{
			text:      "Lorem ipsum dolor sit amet,\r\nconsectetur adipiscing elit.",
			maxLength: 40,
			suffix:    "...",
			exact:     false,
			expected:  "Lorem ipsum dolor sit amet,\r\n...",
		},
		{
			text:      "Lorem ipsum dolor sit amet,\r\nconsectetur adipiscing elit.",
			maxLength: 40,
			suffix:    "...",
			exact:     true,
			expected:  "Lorem ipsum dolor sit amet,\r\nconsecte...",
		},

		{
			text:      "  Lorem ipsum   dolor\nsit\r\namet,\t  consectetur adipiscing elit.",
			maxLength: 45,
			suffix:    "…",
			exact:     false,
			expected:  "  Lorem ipsum   dolor\nsit\r\namet,\t  …",
		},
		{
			text:      "  Lorem ipsum   dolor\nsit\r\namet,\t  consectetur adipiscing elit.",
			maxLength: 45,
			suffix:    "…",
			exact:     true,
			expected:  "  Lorem ipsum   dolor\nsit\r\namet,\t  consectet…",
		},
		{
			text:      "  Lorem ipsum   dolor\nsit\r\namet,\t  consectetur adipiscing elit.",
			maxLength: 45,
			suffix:    "...",
			exact:     false,
			expected:  "  Lorem ipsum   dolor\nsit\r\namet,\t  ...",
		},
		{
			text:      "  Lorem ipsum   dolor\nsit\r\namet,\t  consectetur adipiscing elit.",
			maxLength: 45,
			suffix:    "...",
			exact:     true,
			expected:  "  Lorem ipsum   dolor\nsit\r\namet,\t  consect...",
		},
	}

	for _, test := range tests {
		if result := Truncate(test.text, test.maxLength, test.suffix, test.exact); result != test.expected {
			t.Errorf(
				"Truncate(%q, %d, %q, %t) returned %q, expected %q.",
				test.text,
				test.maxLength,
				test.suffix,
				test.exact,
				result,
				test.expected,
			)
		}
	}
}
