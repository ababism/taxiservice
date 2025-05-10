package stringset

type Set map[string]struct{}

func (s Set) Add(item string) {
	s[item] = struct{}{}
}
func New(item ...string) Set {
	s := Set(make(map[string]struct{}))
	s.AddItems(item)
	return s
}

func (s Set) Intersect(other Set) Set {
	result := Set{}
	if s == nil || other == nil {
		return result
	}
	for item := range s {
		if _, found := other[item]; found {
			result.Add(item)
		}
	}
	return result
}

func (s Set) AddItems(items []string) {
	for _, item := range items {
		s[item] = struct{}{}
	}
}

func (s Set) Remove(item string) {
	delete(s, item)
}

func (s Set) Contains(item string) bool {
	if s == nil {
		return false
	}
	_, found := s[item]
	return found
}

func (s Set) Size() int {
	return len(s)
}

func (s Set) Items() []string {
	items := make([]string, 0, len(s))
	for item := range s {
		items = append(items, item)
	}
	return items
}

//func (s Set) Join(set Set) Set {
//	if s == nil {
//		return set
//	}
//	for item := range set {
//		s[item] = struct{}{}
//	}
//	return s
//}
