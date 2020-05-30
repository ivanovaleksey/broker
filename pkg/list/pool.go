package list

import (
	"sync"
)

var ElementPoolCount int32
var ElementPool = sync.Pool{
	New: func() interface{} {
		// atomic.AddInt32(&ElementPoolCount, 1)
		return new(Element)
	},
}

func GetElementFromPool() *Element {
	el := ElementPool.Get().(*Element)
	return el
}

func PutElementToPool(el *Element) {
	el.next = nil
	el.Parts = nil
	el.Node = nil
	ElementPool.Put(el)
}
