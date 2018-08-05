# rotlog

Daily rotated log file implementation. 
Opens a file suffixed by today's date, e.g. `common.prefix.2018-08-01.log`, 
which can (or will) be added to default loggers output streams, and deletes older files with the same pattern.

## Installation

	$ go get github.com/todmitry/rotlog

## Example. Both uses are equivalent.

```go
package main

import (
	"log"
	"github.com/todmitry/rotlog"
)

func main() {
	rotlog.Set("/tmp", "rotlog.test", 9, os.Stdout)
	log.Println("Look for a file whose name starts with 'rotlog.test.' in the /tmp folder")

	file, _ := rotlog.Rotate("/tmp", "rotlog.more", 3)
	log.SetOutput(io.MultiWriter(file, os.Stderr))
	log.Println("There should be at least 1 but no more than 3 files whose name starts with 'rotlog.more.' in the /tmp folder")
}
```


