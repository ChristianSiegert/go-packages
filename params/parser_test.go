package params_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/ChristianSiegert/go-packages/params"
)

// Dest1 is a destination for writing parsed URL values into.
type Dest1 struct {
	Bool1    bool
	Bool2    bool
	Bool3    bool
	Bool4    bool
	Bool5    bool
	Bool6    bool
	Float32  float32
	Float64  float64
	Int      int
	Int8     int8
	Int16    int16
	Int32    int32
	Int64    int64
	String1  string
	String2  string `param:"custom_name_string2"`
	Uint     uint
	Uint8    uint8
	Uint16   uint16
	Uint32   uint32
	Uint64   uint64
	Sbool1   []bool
	Sbool2   []bool
	Sbool3   []bool
	Sbool4   []bool
	Sbool5   []bool
	Sbool6   []bool
	Sfloat32 []float32
	Sfloat64 []float64
	Sint     []int
	Sint8    []int8
	Sint16   []int16
	Sint32   []int32
	Sint64   []int64
	Sstring1 []string
	Sstring2 []string `param:"custom_name_sstring2"`
	Suint    []uint
	Suint8   []uint8
	Suint16  []uint16
	Suint32  []uint32
	Suint64  []uint64
}

type Dest3 struct {
	Map map[string]string
}

// methods are HTTP methods the parser must support when parsing URL values.
var methods = []string{
	http.MethodConnect,
	http.MethodDelete,
	http.MethodGet,
	http.MethodHead,
	http.MethodOptions,
	http.MethodPatch,
	http.MethodPost,
	http.MethodPut,
	http.MethodTrace,
}

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		inputDest   interface{}
		inputParams url.Values
		expected    interface{}
		expectErr   bool
	}{
		{
			inputDest: &Dest1{},
			inputParams: url.Values{
				"Bool1":                []string{"1", "0"},
				"Bool2":                []string{"true", "false"},
				"Bool3":                []string{"yes", "no"},
				"Bool4":                []string{"0", "1"},
				"Bool5":                []string{"false", "true"},
				"Bool6":                []string{"no", "yes"},
				"Float32":              []string{"2", "3"},
				"Float64":              []string{"4", "5"},
				"Int":                  []string{"6", "7"},
				"Int8":                 []string{"8", "9"},
				"Int16":                []string{"10", "11"},
				"Int32":                []string{"12", "13"},
				"Int64":                []string{"14", "15"},
				"String1":              []string{"16", "17"},
				"Uint":                 []string{"18", "19"},
				"Uint8":                []string{"20", "21"},
				"Uint16":               []string{"22", "23"},
				"Uint32":               []string{"24", "25"},
				"Uint64":               []string{"26", "27"},
				"Sbool1":               []string{"1", "0"},
				"Sbool2":               []string{"true", "false"},
				"Sbool3":               []string{"yes", "no"},
				"Sbool4":               []string{"0", "1"},
				"Sbool5":               []string{"false", "true"},
				"Sbool6":               []string{"no", "yes"},
				"Sfloat32":             []string{"28", "29"},
				"Sfloat64":             []string{"30", "31"},
				"Sint":                 []string{"32", "33"},
				"Sint8":                []string{"34", "35"},
				"Sint16":               []string{"36", "37"},
				"Sint32":               []string{"38", "39"},
				"Sint64":               []string{"40", "41"},
				"Sstring1":             []string{"42", "43"},
				"Suint":                []string{"44", "45"},
				"Suint8":               []string{"46", "47"},
				"Suint16":              []string{"48", "49"},
				"Suint32":              []string{"50", "51"},
				"Suint64":              []string{"52", "53"},
				"custom_name_string2":  []string{"lorem"},
				"custom_name_sstring2": []string{"lorem", "ipsum"},
			},
			expected: &Dest1{
				Bool1:    true,
				Bool2:    true,
				Bool3:    true,
				Bool4:    false,
				Bool5:    false,
				Bool6:    false,
				Float32:  2,
				Float64:  4,
				Int:      6,
				Int8:     8,
				Int16:    10,
				Int32:    12,
				Int64:    14,
				String1:  "16",
				String2:  "lorem",
				Uint:     18,
				Uint8:    20,
				Uint16:   22,
				Uint32:   24,
				Uint64:   26,
				Sbool1:   []bool{true, false},
				Sbool2:   []bool{true, false},
				Sbool3:   []bool{true, false},
				Sbool4:   []bool{false, true},
				Sbool5:   []bool{false, true},
				Sbool6:   []bool{false, true},
				Sfloat32: []float32{28, 29},
				Sfloat64: []float64{30, 31},
				Sint:     []int{32, 33},
				Sint8:    []int8{34, 35},
				Sint16:   []int16{36, 37},
				Sint32:   []int32{38, 39},
				Sint64:   []int64{40, 41},
				Sstring1: []string{"42", "43"},
				Sstring2: []string{"lorem", "ipsum"},
				Suint:    []uint{44, 45},
				Suint8:   []uint8{46, 47},
				Suint16:  []uint16{48, 49},
				Suint32:  []uint32{50, 51},
				Suint64:  []uint64{52, 53},
			},
		},
		{
			inputDest: &Dest1{},
			inputParams: url.Values{
				"Bool1":                []string{"", "1", "0"},
				"Bool2":                []string{"", "true", "false"},
				"Bool3":                []string{"", "yes", "no"},
				"Bool4":                []string{"", "0", "1"},
				"Bool5":                []string{"", "false", "true"},
				"Bool6":                []string{"", "no", "yes"},
				"Float32":              []string{"", "2", "3"},
				"Float64":              []string{"", "4", "5"},
				"Int":                  []string{"", "6", "7"},
				"Int8":                 []string{"", "8", "9"},
				"Int16":                []string{"", "10", "11"},
				"Int32":                []string{"", "12", "13"},
				"Int64":                []string{"", "14", "15"},
				"String1":              []string{"", "16", "17"},
				"Uint":                 []string{"", "18", "19"},
				"Uint8":                []string{"", "20", "21"},
				"Uint16":               []string{"", "22", "23"},
				"Uint32":               []string{"", "24", "25"},
				"Uint64":               []string{"", "26", "27"},
				"Sbool1":               []string{"", "1", "0"},
				"Sbool2":               []string{"", "true", "false"},
				"Sbool3":               []string{"", "yes", "no"},
				"Sbool4":               []string{"", "0", "1"},
				"Sbool5":               []string{"", "false", "true"},
				"Sbool6":               []string{"", "no", "yes"},
				"Sfloat32":             []string{"", "28", "29"},
				"Sfloat64":             []string{"", "30", "31"},
				"Sint":                 []string{"", "32", "33"},
				"Sint8":                []string{"", "34", "35"},
				"Sint16":               []string{"", "36", "37"},
				"Sint32":               []string{"", "38", "39"},
				"Sint64":               []string{"", "40", "41"},
				"Sstring1":             []string{"", "42", "43"},
				"Suint":                []string{"", "44", "45"},
				"Suint8":               []string{"", "46", "47"},
				"Suint16":              []string{"", "48", "49"},
				"Suint32":              []string{"", "50", "51"},
				"Suint64":              []string{"", "52", "53"},
				"custom_name_string2":  []string{"", "lorem"},
				"custom_name_sstring2": []string{"", "lorem", "ipsum"},
			},
			expected: &Dest1{
				Bool1:    false,
				Bool2:    false,
				Bool3:    false,
				Bool4:    false,
				Bool5:    false,
				Bool6:    false,
				Float32:  0,
				Float64:  0,
				Int:      0,
				Int8:     0,
				Int16:    0,
				Int32:    0,
				Int64:    0,
				String1:  "",
				String2:  "",
				Uint:     0,
				Uint8:    0,
				Uint16:   0,
				Uint32:   0,
				Uint64:   0,
				Sbool1:   []bool{false, true, false},
				Sbool2:   []bool{false, true, false},
				Sbool3:   []bool{false, true, false},
				Sbool4:   []bool{false, false, true},
				Sbool5:   []bool{false, false, true},
				Sbool6:   []bool{false, false, true},
				Sfloat32: []float32{0, 28, 29},
				Sfloat64: []float64{0, 30, 31},
				Sint:     []int{0, 32, 33},
				Sint8:    []int8{0, 34, 35},
				Sint16:   []int16{0, 36, 37},
				Sint32:   []int32{0, 38, 39},
				Sint64:   []int64{0, 40, 41},
				Sstring1: []string{"", "42", "43"},
				Sstring2: []string{"", "lorem", "ipsum"},
				Suint:    []uint{0, 44, 45},
				Suint8:   []uint8{0, 46, 47},
				Suint16:  []uint16{0, 48, 49},
				Suint32:  []uint32{0, 50, 51},
				Suint64:  []uint64{0, 52, 53},
			},
		},
		// Test passing unsupported type
		{
			inputDest:   &Dest3{},
			inputParams: url.Values{"Map": []string{"foo"}},
			expected:    &Dest3{},
			expectErr:   true,
		},
		// Test not passing pointer to struct
		{
			inputDest: Dest1{},
			expectErr: true,
		},
		// Test empty request parameters
		{
			inputDest:   &Dest1{},
			inputParams: nil,
			expected:    &Dest1{},
		},
	}

	for _, test := range tests {
		for _, usePostForm := range []bool{false, true} {
			for _, method := range methods {
				// Create request that contains the parameters
				request := httptest.NewRequest(method, "/", nil)
				if usePostForm {
					request.PostForm = test.inputParams
				} else {
					request.Form = test.inputParams
				}

				// Create parser
				parser, err := params.NewParser(request, nil)
				if err != nil {
					t.Fatal(err)
				}

				// Test
				dest := test.inputDest

				if err := parser.Parse(dest); err != nil {
					if !test.expectErr {
						t.Fatalf("Parse failed: unexpected error: %s", err)
					}
					continue
				} else if test.expectErr {
					t.Fatalf("Parse failed: no error occured, expected error")
				} else if !reflect.DeepEqual(dest, test.expected) {
					t.Fatalf("Parse failed:\nexpected %#v\n\ngot %#v", test.expected, dest)
				}
			}
		}
	}
}
