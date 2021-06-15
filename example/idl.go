// +build idl

package example

import (
	"idl/example/store"
)

type appleBan struct {
	store.BuyPen  `path:"/store/buy/pen" typ:"json"`
}
