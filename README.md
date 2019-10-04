# encoding
Go encoding/json style unmarshaling of environment variables to struct. 

Currently returns error if struct field is not found in environment

# Types
Currently supported types:
```Go
string
bool
int64
int
uint64
[]byte
map[string]bool
[]string
```

# Syntax
```Go
/* Environment variables:
v1=Yolo
v2=1312
VALUE4=swag
VALUE5=string1,string2
*/

import (
    "github.com/captainlettuce/encoding/env"
    "os"
)

type someStruct struct {
    Value1 string `env:"v1"`
    Value2 uint64 `env:"v2"`
    Value3 bool   `env:"-"` // '-' ignores field
    Value4 string           // If no name is given, field name is assumed in all caps, i.e VALUE4
    Value5 []string
}

func example() {
    vars := os.Environ()
    ss := &someStruct{}
    err := env.Unmarshal(vars, ss)
    if err != nil {
        // Error handling
    }
  
    println("%+v", ss)
    // &{Value1: "Yolo", Value2: 1312, Value3: false, Value4: "swag", Value5: []string{"string1", "string2"}}
}

```

# ToDo
Add support for all primitives

Add support for slice/array and maps
