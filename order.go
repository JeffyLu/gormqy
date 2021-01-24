package gormqy

import (
	"fmt"
	"strings"
)

const (
	OrderASC  = "ASC"
	OrderDesc = "DESC"
)

func (q *Query) DESC(column string) *Query {
	q.orders = append(q.orders, fmt.Sprintf("%s %s", column, OrderDesc))
	return q
}

func (q *Query) ASC(column string) *Query {
	q.orders = append(q.orders, fmt.Sprintf("%s %s", column, OrderASC))
	return q
}

func (q *Query) Order() (expr string) {
	return strings.Join(q.orders, ", ")
}
