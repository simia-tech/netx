package value

// Endpoints is a sortable type alias for a slice of Dial.
type Endpoints []Endpoint

// ParseEndpointURLs parses the provided urls and returns endpoints.
func ParseEndpointURLs(urls []string, options ...Option) (Endpoints, error) {
	dials := Endpoints{}
	for _, url := range urls {
		dial, err := ParseEndpointURL(url, options...)
		if err != nil {
			return nil, err
		}
		dials = append(dials, dial)
	}
	return dials, nil
}

// MustParseEndpointURLs works like ParseEndpointURLs, but panics on error.
func MustParseEndpointURLs(urls []string, options ...Option) Endpoints {
	dials, err := ParseEndpointURLs(urls, options...)
	if err != nil {
		panic(err)
	}
	return dials
}

func (e Endpoints) Len() int {
	return len(e)
}

func (e Endpoints) Less(i, j int) bool {
	return e[i].Address() < e[j].Address()
}

func (e Endpoints) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}
