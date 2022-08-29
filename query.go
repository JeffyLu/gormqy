package gormqy

import (
	"fmt"
	"strings"
)

type Query struct {
	whereExprs []string
	whereVals  []interface{}
	selectCols []string
	orders     []string
	limit      uint64
}

func New() *Query {
	return &Query{}
}

func (q *Query) Condition() *Condition {
	return NewCondition(q)
}

func (q *Query) AndCondition() *Condition {
	return q.conditionAfter(LogicAnd, false)
}

func (q *Query) AndConditionAfterGroup() *Condition {
	return q.conditionAfter(LogicAnd, true)
}

func (q *Query) OrCondition() *Condition {
	return q.conditionAfter(LogicOr, false)
}

func (q *Query) OrConditionAfterGroup() *Condition {
	return q.conditionAfter(LogicOr, true)
}

func (q *Query) group() {
	q.whereExprs = []string{fmt.Sprintf("(%s)", strings.Join(q.whereExprs, " "))}
}

func (q *Query) conditionAfter(logic string, groupFirst bool) *Condition {
	if len(q.whereExprs) != 0 {
		if groupFirst {
			q.group()
		}
		q.whereExprs = append(q.whereExprs, logic)
	}
	return q.Condition()
}

func (q *Query) Where() (expr string, vals []interface{}) {
	if len(q.whereExprs) == 0 {
		return "", nil
	}
	return strings.Join(q.whereExprs, " "), q.whereVals
}

func (q *Query) SelectCols(cols ...string) {
	q.selectCols = append(q.selectCols, cols...)
}

func (q *Query) Select() string {
	if len(q.selectCols) == 0 {
		return "*"
	}
	return strings.Join(q.selectCols, ", ")
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
