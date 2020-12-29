package jsonc

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func b(s string) []byte { return []byte(s) }
func s(b []byte) string { return string(b) }

type testsStruct struct {
	validBlock    []byte
	validSingle   []byte
	invalidBlock  []byte
	invalidSingle []byte
}

var jsonTest, jsoncTest testsStruct

func init() {
	jsonTest = testsStruct{
		validBlock:   b(`{"foo":"bar foo","true":false,"number":42,"object":{"test":"done"},"array":[1,2,3],"url":"https://github.com","escape":"\"wo//rking"}`),
		invalidBlock: b(`{"foo":`),
	}
	jsoncTest = testsStruct{
		validBlock:    b(`{"foo": /** this is a bloc/k comm\"ent */ "bar foo", "true": /* true */ false, "number": 42, "object": { "test": "done" }, "array" : [1, 2, 3], "url" : "https://github.com", "escape":"\"wo//rking" }`),
		invalidBlock:  b(`{"foo": /* this is a block comment "bar foo", "true": false, "number": 42, "object": { "test": "done" }, "array" : [1, 2, 3], "url" : "https://github.com", "escape":"\"wo//rking }`),
		validSingle:   b("{\"foo\": // this is a single line comm\\\"ent\n\"bar foo\", \"true\": false, \"number\": 42, \"object\": { \"test\": \"done\" }, \"array\" : [1, 2, 3], \"url\" : \"https://github.com\", \"escape\":\"\\\"wo//rking\" }"),
		invalidSingle: b(`{"foo": // this is a single line comment "bar foo", "true": false, "number": 42, "object": { "test": "done" }, "array" : [1, 2, 3], "url" : "https://github.com", "escape":"\"wo//rking" }`),
	}
}

func TestToJSON(t *testing.T) {
	tests := []struct {
		name    string
		arg     []byte
		want    []byte
		wantErr bool
	}{
		{
			name: "Test for valid block comment.",
			arg:  jsoncTest.validBlock,
			want: jsonTest.validBlock,
		}, {
			name:    "Test for invalid block comment.",
			arg:     jsoncTest.invalidBlock,
			want:    jsonTest.invalidBlock,
			wantErr: true,
		}, {
			name: "Test for valid single line comment.",
			arg:  jsoncTest.validSingle,
			want: jsonTest.validBlock,
		}, {
			name:    "Test for invalid single line comment.",
			arg:     jsoncTest.invalidSingle,
			want:    jsonTest.invalidBlock,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToJSON(tt.arg)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToJSON() = %v, want %v", s(got), s(tt.want))
			}
			if !json.Valid(got) && !tt.wantErr {
				t.Errorf("ToJSON() = %v isn't valid.", s(got))
			}
		})
	}
}

func TestUnmarshal(t *testing.T) {
	type UnmarshalTest struct {
		Foo    string `json:"foo"`
		True   bool   `json:"true"`
		Num    int    `json:"number"`
		Object struct {
			Test string `json:"test"`
		} `json:"object"`
		Array  []int  `json:"array"`
		URL    string `json:"url"`
		Escape string `json:"escape"`
	}

	t.Run("Testing Unmarshal()", func(t *testing.T) {
		un := UnmarshalTest{}
		if err := Unmarshal(jsoncTest.validBlock, &un); err != nil {
			t.Errorf("Unmarshal() error = %v", err)
		}
		mr, err := json.Marshal(un)
		if err != nil {
			t.Errorf("Unmarshal() unable to marshal Unmarshal(). Error = %v, got = %v, want = %v", err, s(mr), s(jsonTest.validBlock))
		}
		if !reflect.DeepEqual(mr, jsonTest.validBlock) {
			t.Errorf("Unmarshal() didn't work correctly. Got = %v, want = %v", s(mr), s(jsonTest.validBlock))
		}
	})
}

func TestReadFromFile(t *testing.T) {
	tmp, err := ioutil.TempFile("", "ReadFromFileTest")
	if err != nil {
		t.Skip("Unable to create temp file.", err)
	}
	defer os.Remove(tmp.Name())

	if _, err := tmp.Write(jsoncTest.validBlock); err != nil {
		t.Skip("Unable to write to file.", err)
	}

	defer func() {
		if err := tmp.Close(); err != nil {
			t.Log("Unable to close the file.", err)
		}
	}()

	tests := []struct {
		name     string
		filename string
		want     []byte
		want1    []byte
		wantErr  bool
	}{
		{
			name:     "Valid use of ReadFromFile()",
			filename: tmp.Name(),
			want:     jsoncTest.validBlock,
			want1:    jsonTest.validBlock,
			wantErr:  false,
		},
		{
			name:     "Invalid use of ReadFromFile()",
			filename: "NonExistentFile",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ReadFromFile(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadFromFile() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ReadFromFile() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestValid(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want bool
	}{
		{
			name: "Valid test",
			data: b(`{"foo":/*comment*/"bar"}`),
			want: true,
		},
		{
			name: "Invalid test",
			data: b(`{"foo"://comment without ending`),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Valid(tt.data); got != tt.want {
				t.Errorf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkTranslate(b *testing.B) {
	for n := 0; n < b.N; n++ {
		translate(jsoncTest.validSingle)
		translate(jsoncTest.validBlock)
	}
}

func BenchmarkValid(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Valid(jsoncTest.validSingle)
		Valid(jsoncTest.validBlock)
	}
}
