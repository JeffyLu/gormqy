package gormqy

import (
	"strings"
)

type Query struct {
	exprs  []string
	vals   []interface{}
	orders []string
	limit  uint64
}

func New() *Query {
	return &Query{}
}

func (q *Query) Condition() *Condition {
	return NewCondition(q)
}

func (q *Query) AndCondition() *Condition {
	return q.conditionAfter(LogicAnd)
}

func (q *Query) OrCondition() *Condition {
	return q.conditionAfter(LogicOr)
}

func (q *Query) conditionAfter(logic string) *Condition {
	if len(q.exprs) != 0 {
		q.exprs = append(q.exprs, logic)
	}
	return q.Condition()
}

func (q *Query) Where() (expr string, vals []interface{}) {
	if len(q.exprs) == 0 {
		return "", nil
	}
	return strings.Join(q.exprs, " "), q.vals
}

func (q *Query) ToPageQuery(page uint64, size uint64) *PageQuery {
	if page == 0 {
		page = 1
	}
	if size == 0 {
		size = 10
	}
	return &PageQuery{
		Query:  q,
		page:   page,
		size:   size,
		offset: (page - 1) * size,
	}
}

type PageQuery struct {
	*Query
	page   uint64
	size   uint64
	offset uint64
}

func (pq *PageQuery) Page() uint64 {
	return pq.page
}

func (pq *PageQuery) Size() uint64 {
	return pq.size
}

func (pq *PageQuery) Offset() uint64 {
	return pq.offset
}
