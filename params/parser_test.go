package params_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strconv"
	"testing"

	"github.com/ChristianSiegert/go-packages/params"
)

type TestStruct1 struct {
	Field1  string
	Field2  float32
	Field3  float64
	Field4  int
	Field5  int8
	Field6  int16
	Field7  int32
	Field8  int64
	Field9  uint
	Field10 uint8
	Field11 uint16
	Field12 uint32
	Field13 uint64

	Field14 []string
	Field15 []float32
	Field16 []float64
	Field17 []int
	Field18 []int8
	Field19 []int16
	Field20 []int32
	Field21 []int64
	Field22 []uint
	Field23 []uint8
	Field24 []uint16
	Field25 []uint32
	Field26 []uint64

	Field27 bool
	Field28 bool
	Field29 bool
	Field30 []bool
	Field31 []bool
	Field32 []bool
}

type TestStruct2 struct {
	SomeFieldName string `param:"real_field_name"`
}

type TestStruct3 struct {
	Field map[string]string
}

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
		arg          interface{}
		parameters   url.Values
		expectedDest interface{}
		expectedErr  bool
	}{
		{
			arg:          &TestStruct1{},
			parameters:   test1URLValues(),
			expectedDest: test1Expected(),
		},
		{
			arg:          &TestStruct1{},
			parameters:   url.Values{"field27": []string{"no", "yes"}},
			expectedDest: &TestStruct1{Field27: false},
		},
		{
			arg:          &TestStruct1{},
			parameters:   url.Values{"field27": []string{"notbool"}},
			expectedDest: &TestStruct1{},
			expectedErr:  true,
		},
		{
			arg:          &TestStruct2{},
			parameters:   url.Values{"real_field_name": []string{"lorem"}},
			expectedDest: &TestStruct2{"lorem"},
		},
		{
			arg:         &TestStruct3{},
			parameters:  url.Values{"field": []string{"foo"}},
			expectedErr: true,
		},
		{
			arg:         TestStruct1{},
			expectedErr: true,
		},
		{
			arg:          &TestStruct1{},
			parameters:   nil,
			expectedDest: &TestStruct1{},
		},
	}

	for _, test := range tests {
		for _, usePostForm := range []bool{false, true} {
			for _, method := range methods {
				// Create request that contains the parameters
				request := httptest.NewRequest(method, "/", nil)
				if usePostForm {
					request.PostForm = test.parameters
				} else {
					request.Form = test.parameters
				}

				// Create parser
				parser, err := params.NewParser(request, nil)
				if err != nil {
					t.Fatal(err)
				}

				// Test
				dest := test.arg

				if err := parser.Parse(dest); err != nil {
					if !test.expectedErr {
						t.Fatalf("Parse failed: unexpected error %q", err)
					}
					continue
				} else if test.expectedErr {
					t.Fatalf("Parse failed: no error occured, expected error")
				}

				if !reflect.DeepEqual(dest, test.expectedDest) {
					t.Fatalf("Parse failed:\n%#v\n%#v", test.expectedDest, dest)
				}
			}
		}
	}
}

func test1URLValues() url.Values {
	urlValues := url.Values{}
	for i := 0; i < 27; i++ {
		name := "field" + strconv.Itoa(i+1)
		value1 := strconv.Itoa(2*i + 1)
		value2 := strconv.Itoa(2*i + 2)
		urlValues[name] = []string{value1, value2}
	}
	urlValues["field27"] = []string{"1", "0"}
	urlValues["field28"] = []string{"true", "false"}
	urlValues["field29"] = []string{"yes", "no"}
	urlValues["field30"] = urlValues["field27"]
	urlValues["field31"] = urlValues["field28"]
	urlValues["field32"] = urlValues["field29"]
	return urlValues
}

func test1Expected() interface{} {
	return &TestStruct1{
		Field1:  "1",
		Field2:  3.0,
		Field3:  5.0,
		Field4:  7,
		Field5:  9,
		Field6:  11,
		Field7:  13,
		Field8:  15,
		Field9:  17,
		Field10: 19,
		Field11: 21,
		Field12: 23,
		Field13: 25,

		Field14: []string{"27", "28"},
		Field15: []float32{29, 30},
		Field16: []float64{31, 32},
		Field17: []int{33, 34},
		Field18: []int8{35, 36},
		Field19: []int16{37, 38},
		Field20: []int32{39, 40},
		Field21: []int64{41, 42},
		Field22: []uint{43, 44},
		Field23: []uint8{45, 46},
		Field24: []uint16{47, 48},
		Field25: []uint32{49, 50},
		Field26: []uint64{51, 52},

		Field27: true,
		Field28: true,
		Field29: true,
		Field30: []bool{true, false},
		Field31: []bool{true, false},
		Field32: []bool{true, false},
	}
}
