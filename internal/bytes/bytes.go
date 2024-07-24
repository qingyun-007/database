package bytes

import "fmt"

type Buffer struct {
	Content interface{}
}

func (b Buffer) String() string {
	return fmt.Sprintf("%v", b.Content)
}

var Stderr interface{}
var Stdout interface{}
