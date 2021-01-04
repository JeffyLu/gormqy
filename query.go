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
	OpEq    Operator = "="
	OpLe    Operator = "<="
	OpGe    Operator = ">="
	OpLt    Operator = "<"
	OpGt    Operator = "<"
	OpNotEq Operator = "!="
	OpLike  Operator = "LIKE"
	OpIn    Operator = "IN"
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

func ConcatCol(cols ...Column) Column {
	var cs []string
	for _, c := range cols {
		cs = append(cs, string(c))
	}
	return Column(fmt.Sprintf("concat_ws(' | ', %s)", strings.Join(cs, ", ")))
}

type Query struct {
	condCols []string
	condVals []interface{}
	orders   []string
	limit    uint64
}

func (q *Query) OrderExpr() string {
	return strings.Join(q.orders, ", ")
}

func (q *Query) WhereExpr() (string, []interface{}) {
	if len(q.condCols) == 0 {
		return "", nil
	}
	return strings.Join(q.condCols[:len(q.condCols)-1], " "), q.condVals
}

func (q *Query) LimitExpr() uint64 {
	return q.limit
}

func (q *Query) Condition(col Column, op Operator, value interface{}, logic Logic) *Query {
	if logic != LogicOr {
		logic = LogicAnd
	}
	var c string
	var v interface{}
	switch op {
	case OpLike:
		c = fmt.Sprintf("%s LIKE ?", col)
		v = fmt.Sprintf("%%%s%%", value)
	case OpIn:
		c = fmt.Sprintf("%s IN (?)", col)
	default:
		c = fmt.Sprintf("%s %s ?", col, op)
	}
	q.condCols = append(q.condCols, c, string(logic))
	q.condVals = append(q.condVals, v)
	return q
}

func (q *Query) Order(col Column, method OrderMethod) *Query {
	if method != OrderASC {
		method = OrderDesc
	}
	q.orders = append(q.orders, fmt.Sprintf("%s %s", col, method))
	return q
}

func (q *Query) Limit(limit uint64) *Query {
	q.limit = limit
	return q
}

type PageQuery struct {
	Query
	Page   uint
	Size   uint
	Offset uint
}

func (pq *PageQuery) Validate() *PageQuery {
	if pq.Size == 0 {
		pq.Size = 10
	}
	if pq.Page == 0 {
		pq.Page = 1
	}
	pq.Offset = (pq.Page - 1) * pq.Size
	return pq
}
