package queueCollection

type QueueCollections[TKey comparable, TVal any] struct {
	collections map[TKey]*queue[TVal]
}

func NewQueueCollections[TKey comparable, TVal any]() QueueCollections[TKey, TVal] {
	qc := QueueCollections[TKey, TVal]{}
	qc.collections = make(map[TKey]*queue[TVal])
	return qc
}

func (qc *QueueCollections[TKey, TVal]) Set(key TKey, val TVal) error {
	q, ok := qc.collections[key]
	if !ok {
		q = newQueue[TVal]()
		qc.collections[key] = q
	}
	err := q.set(val)
	return err
}

func (qc *QueueCollections[TKey, TVal]) Get() (vals []TVal) {
	for _, q := range qc.collections {
		val, err := q.get()
		if err == nil {
			vals = append(vals, val)
		}
	}
	return vals
}

func (qc *QueueCollections[TKey, TVal]) ValNum(key TKey) int {
	q, ok := qc.collections[key]
	if !ok {
		return 0
	}
	return q.len
}
