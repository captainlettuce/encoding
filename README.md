# encoding
Go encoding/json style unmarshaling of environment variables to struct. 

# Types
Currently supported types:
```Go
string
bool
int64
uint64
```

# Syntax
```Go
/* Environment variables:
v1=Yolo
v2=1312
*/

import (
    "github.com/captainlettuce/encoding/env"
    "os"
)

type someStruct struct {
    Value1 string `env:"v1"` // If no type tag is given it's assumed to be a string
    Value2 uint64 `env:"v2,uint64"`
}

func example() {
    vars := os.Environ()
    ss := &someStruct{}
    err := env.Unmarshal(vars, ss)
    if err != nil {
        // Error handling
    }
  
    println("%+v", ss)
    // &{Value1: Yolo, Value2: 1312}
}

```

# ToDo
Currently not inferring name and type from struct definition. I.e. Tags are neccessary.

Add support for all primitives

Add support for slice/array and maps
