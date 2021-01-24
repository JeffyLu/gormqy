package gormqy

func (q *Query) AddLimit(limit uint64) *Query {
	q.limit = limit
	return q
}

func (q *Query) Limit() uint64 {
	return q.limit
}
