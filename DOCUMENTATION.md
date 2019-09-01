# Documentation for JSONC

## ToJSON()

Converts JSONC to JSON equivalent by removing all comments.

```go
func ToJSON(b []byte) []byte
```

Example:

```go
func main() {
  j := []byte(`{"foo": /*comment*/ "bar"}`)
  jc := jsonc.ToJSON(j)
  fmt.Println(string(jc)) // {"foo":"bar"}
}
```

## ReadFromFile()

Reads jsonc file and returns JSONC and JSON encodings.

```go
func ReadFromFile(filename string) ([]byte, []byte, error)
```

Example:

```go
func main() {
  jc, j, err := jsonc.ReadFromFile("data.jsonc")
  if err != nil {
    log.Fatal(err)
  }
  // jc and j contains JSONC and JSON, respectively.
}
```

## Unmarshal()

Parses the JSONC-encoded data and stores the result in the value pointed to by passed interface. Equivalent of calling `json.Unmarshal(jsonc.ToJSON(data), v)`.

```go
func Unmarshal(data []byte, v interface{}) error
```

Example:

```go
func main(){
  type UnmarshalTest struct {
    Foo    string `json:"foo"`
    True   bool   `json:"true"`
    Num    int    `json:"number"`
    Object struct {
      Test string `json:"test"`
    } `json:"object"`
    Array []int `json:"array"`
  }

  un := UnmarshalTest{}
  jc, _, _ := jsonc.ReadFromFile("data.jsonc")
  if err := jsonc.Unmarshal(jc, &un); err != nil {
    log.Fatal("Unable to unmarshal.", err)
  }
  fmt.Println(string(jc))
  fmt.Printf("%+v", un)
}
```

Output:

```sh
$ go run .\main.go
{
  /* This is an example
     for block comment. */
  "foo": "bar foo", // Comments can
  "true": false, // Improve readability.
  "number": 42, // Number will always be 42.
  /* Comments are ignored while
     generating JSON from JSONC. */
  "object": {
    "test": "done"
  },
  "array": [1, 2, 3]
}

{Foo:bar foo True:false Num:42 Object:{Test:done} Array:[1 2 3]}
```

## Valid()

Reports whether data is a valid JSONC encoding or not.

```go
func Valid(data []byte) bool
```

Example:

```go
func main() {
  jc1 := []byte(`{"foo":/*comment*/"bar"}`)
  jc2 := []byte(`{"foo":/*comment/"bar"}`)
  fmt.Println(jsonc.Valid(jc1)) // true
  fmt.Println(jsonc.Valid(jc2)) // false
}
```
