package functional

// Any of the element of the list matches the predicate
// 'a list -> 'a -> bool -> bool
func ListAnyOf[T any](l []T, p func(T) bool) bool {
	for _, v := range l {
		if p(v) {
			return true
		}
	}
	return false
}

// None of the element of the list matches the predicate
// 'a list -> 'a -> bool -> bool
func ListNoneOf[T any](l []T, p func(T) bool) bool {
	return !ListAnyOf(l, p)
}

// list map
// 'a list -> 'a -> 'b -> 'b list
func ListMap[T, U any](l []T, convert func(T) U) []U {
	var nl []U
	for _, v := range l {
		nl = append(nl, convert(v))
	}
	return nl
}

// list try map
// 'a list -> 'a -> 'b -> 'b list
func ListTryMap[T, U any](l []T, convert func(T) (U, error)) ([]U, error) {
	var nl []U
	for _, v := range l {
		cv, err := convert(v)
		if err != nil {
			return nil, err
		}
		nl = append(nl, cv)
	}
	return nl, nil
}

// list filter
// 'a list -> 'a -> bool -> 'a list
func ListFilter[T any](l []T, predicate func(T) bool) []T {
	var nl []T
	for _, v := range l {
		if predicate(v) {
			nl = append(nl, v)
		}
	}
	return nl
}

// list try apply
// 'a list -> 'a -> bool -> 'a list
func ListTryApply[T any](l []T, apply func(T) error) error {
	for _, v := range l {
		if err := apply(v); err != nil {
			return err
		}
	}
	return nil
}

// list find if
// 'a list -> 'a -> bool -> 'a ptr
func ListFindIf[T any](l []T, predicate func(T) bool) *T {
	for _, v := range l {
		if predicate(v) {
			return &v
		}
	}
	return nil
}

func ListIn[T comparable](l []T, v T) bool {
	for _, u := range l {
		if u == v {
			return true
		}
	}
	return false
}

func ListUniq[T comparable](l []T) []T {
	var nl []T
	for _, v := range l {
		found := false
		for _, w := range nl {
			if w == v {
				found = true
				break
			}
		}
		if !found {
			nl = append(nl, v)
		}
	}
	return nl
}

func ListEqual[T comparable](l1, l2 []T) bool {
	if len(l1) != len(l2) {
		return false
	}
	for i := 0; i < len(l1); i++ {
		if l1[i] != l2[i] {
			return false
		}
	}
	return true
}
