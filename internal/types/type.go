package types

import "fmt"

type Type string

const (
	ORDER_BOOK Type = "ORDER_BOOK"
)

func (t *Type) String() string {
	return fmt.Sprintf("%v", *t)
}
