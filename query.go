package gormqy

import (
	"fmt"
	"strings"
)

type Column string
type Operator string
type Logic string
type OrderMethod string

const (
	OpEq      Operator = "="
	OpLe      Operator = "<="
	OpGe      Operator = ">="
	OpLt      Operator = "<"
	OpGt      Operator = ">"
	OpNotEq   Operator = "!="
	OpContain Operator = "Contain"
	OpPrefix  Operator = "Prefix"
	OpSuffix  Operator = "Suffix"
	OpIn      Operator = "IN"
)

const (
	LogicOr  Logic = "OR"
	LogicAnd Logic = "AND"
)

const (
	OrderASC  OrderMethod = "ASC"
	OrderDesc OrderMethod = "DESC"
)

func Col(col string) Column {
	return Column(col)
}

func ConcatCol(sep string, cols ...Column) Column {
	cs := make([]string, 0, len(cols))
	for _, c := range cols {
		cs = append(cs, string(c))
	}
	return Column(fmt.Sprintf("concat_ws('%s', %s)", sep, strings.Join(cs, ", ")))
}

type Query struct {
	condCols []string
	condVals []interface{}
	orders   []string
	limit    uint64
}

func NewQuery() *Query {
	return &Query{}
}

func (q *Query) AddOrder(col Column, method OrderMethod) *Query {
	if method != OrderASC {
		method = OrderDesc
	}
	q.orders = append(q.orders, fmt.Sprintf("%s %s", col, method))
	return q
}

func (q *Query) Order() (expr string) {
	return strings.Join(q.orders, ", ")
}

func (q *Query) AddCondition(col Column, op Operator, value interface{}, logic Logic) *Query {
	if logic != LogicOr {
		logic = LogicAnd
	}
	var c string
	v := value
	switch op {
	case OpContain:
		c = fmt.Sprintf("%s LIKE ?", col)
		v = fmt.Sprintf("%%%s%%", value)
	case OpPrefix:
		c = fmt.Sprintf("%s LIKE ?", col)
		v = fmt.Sprintf("%s%%", value)
	case OpSuffix:
		c = fmt.Sprintf("%s LIKE ?", col)
		v = fmt.Sprintf("%%%s", value)
	case OpIn:
		c = fmt.Sprintf("%s IN (?)", col)
	default:
		c = fmt.Sprintf("%s %s ?", col, op)
	}
	q.condCols = append(q.condCols, c, string(logic))
	q.condVals = append(q.condVals, v)
	return q
}

func (q *Query) GroupConditions(logic Logic) *Query {
	q.condCols = []string{
		fmt.Sprintf("(%s)", strings.Join(q.condCols[:len(q.condCols)-1], " ")),
		string(logic),
	}
	return q
}

func (q *Query) Where() (expr string, vals []interface{}) {
	if len(q.condCols) == 0 {
		return "", nil
	}
	return strings.Join(q.condCols[:len(q.condCols)-1], " "), q.condVals
}

func (q *Query) AddLimit(limit uint64) *Query {
	q.limit = limit
	return q
}

func (q *Query) Limit() uint64 {
	return q.limit
}

type PageQuery struct {
	Query
	page   uint64
	size   uint64
	offset uint64
}

func NewPageQuery(page uint64, size uint64) *PageQuery {
	if page == 0 {
		page = 1
	}
	if size == 0 {
		size = 10
	}
	return &PageQuery{
		page:   page,
		size:   size,
		offset: (page - 1) * size,
	}
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

func (pq *PageQuery) AddOrder(col Column, method OrderMethod) *PageQuery {
	pq.Query.AddOrder(col, method)
	return pq
}

func (pq *PageQuery) AddCondition(col Column, op Operator, value interface{}, logic Logic) *PageQuery {
	pq.Query.AddCondition(col, op, value, logic)
	return pq
}

func (pq *PageQuery) GroupConditions(logic Logic) *PageQuery {
	pq.Query.condCols = []string{
		fmt.Sprintf("(%s)", strings.Join(pq.Query.condCols[:len(pq.Query.condCols)-1], " ")),
		string(logic),
	}
	return pq
}

func (pq *PageQuery) AddLimit(limit uint64) *PageQuery {
	pq.Query.limit = limit
	return pq
}
