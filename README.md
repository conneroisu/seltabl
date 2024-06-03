# seltabl
[
https://www.phorm.ai/query?projectId=3e6f9d42-0098-4178-ab54-4a0b9c89353b
A golang library for configurably parsing html tables into stucts.

## Installation

```bash
go get github.com/conneroisu/seltabl
```

## Usage

```go
package main

import (
	"fmt"
	"github.com/conneroisu/seltabl"
	"github.com/conneroisu/seltabl/testdata"
)

type fixtureStruct struct {
	A string `json:"a" seltabl:"a" hSel:"tr:nth-child(1) td:nth-child(1)" dSel:"tr td:nth-child(1)" cSel:"$text"`
	B string `json:"b" seltabl:"b" hSel:"tr:nth-child(1) td:nth-child(2)" dSel:"tr td:nth-child(2)" cSel:"$text"`
}

func main() {
	p, err := seltabl.NewFromString[fixtureStruct](testdata.FixtureABNumTable)
	if err != nil {
		panic(err)
	}
	for _, pp := range p {
		fmt.Printf("pp %+v\n", pp)
	}
}
```

Output:

```bash
pp {A:1 B:2}
pp {A:3 B:4}
pp {A:5 B:6}
pp {A:7 B:8}
```
