package value

// Dials is a sortable type alias for a slice of Dial.
type Dials []Dial

// ParseDialURLs parses the provided urls and returns dials.
func ParseDialURLs(urls []string, options ...DialOption) (Dials, error) {
	dials := Dials{}
	for _, url := range urls {
		dial, err := ParseDialURL(url, options...)
		if err != nil {
			return nil, err
		}
		dials = append(dials, dial)
	}
	return dials, nil
}

// MustParseDialURLs works like ParseDialURLs, but panics on error.
func MustParseDialURLs(urls []string, options ...DialOption) Dials {
	dials, err := ParseDialURLs(urls, options...)
	if err != nil {
		panic(err)
	}
	return dials
}

func (a Dials) Len() int {
	return len(a)
}

func (a Dials) Less(i, j int) bool {
	return a[i].Address() < a[j].Address()
}

func (a Dials) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
