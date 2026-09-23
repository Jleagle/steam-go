package steamvdf

import "testing"

// FuzzReadBytes checks that no input, text or binary, makes the parsers panic.
func FuzzReadBytes(f *testing.F) {
	seeds := []string{
		"",
		`"root" { "k" "v" }`,
		`"root" { "a" { "b" "c" } "d" { } }`,
		`"a"`,
		"// comment only",
		"{",
		"}",
		`"root" { "k" "v"`,
		`"root" { "k" "v" } }`,
		`"root" { "k" "C:\Games" "n" "v" }`,
		`"root" { "q" "a\"b" "b" "\\" }`,
		"\"root\" { \"k\" \"v\" }\x00",
		"\x01\x00\x00\x00\x00root\x00\x01s\x00v\x00\x02i\x00\x01\x00\x00\x00\x08\x08",
	}
	for _, s := range seeds {
		f.Add([]byte(s))
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		_, _ = ReadBytes(data)
	})
}
