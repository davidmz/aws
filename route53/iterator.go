package route53

type iterator struct {
	fetcher   itFetcher
	values    []interface{}
	lastRead  bool
	lastError error
}

type itFetcher interface {
	itFetch() (vals []interface{}, last bool, err error)
}

func newIterator(f itFetcher) *iterator { return &iterator{fetcher: f} }

func (it *iterator) Next() bool {
	if len(it.values) > 1 {
		it.values = it.values[1:]
	} else if !it.lastRead {
		vals, last, err := it.fetcher.itFetch()
		if err != nil {
			it.lastError = err
			return false
		}
		it.values = vals
		it.lastRead = last
	} else {
		it.values = nil
	}
	return !(len(it.values) == 0 && it.lastRead)
}

func (it *iterator) Value() interface{} {
	if len(it.values) == 0 {
		return nil
	}
	return it.values[0]
}

func (it *iterator) Error() error { return it.lastError }
