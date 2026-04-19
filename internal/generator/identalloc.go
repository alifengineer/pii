package generator

import "fmt"

type identifierAllocator map[string]struct{}

func newIdentifierAllocator() identifierAllocator {
	return make(identifierAllocator)
}

func (ia identifierAllocator) allocate(want string) string {
	if _, exists := ia[want]; !exists {
		ia[want] = struct{}{}
		return want
	}
	for i := 0; ; i++ {
		name := fmt.Sprintf("%s_%d", want, i)
		if _, exists := ia[name]; !exists {
			ia[name] = struct{}{}
			return name
		}
	}
}
