// +build idl

package store

import "github.com/idlgo/example1/pen"

type Request struct {
	PenCount int
}

type Response struct {
	PenL []pen.GetPenResponse
	PenM *map[string]pen.GetPenResponse
}

type (
	BuyPen func(request pen.GetPenRequest) Response
)
