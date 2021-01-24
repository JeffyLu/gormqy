package gormqy

import (
	"fmt"
	"strings"
)

const (
	LogicAnd = "AND"
	LogicOr  = "OR"
)

type Logic struct {
	cond *Condition
}

func (l *Logic) And() *Condition {
	l.cond.exprs = append(l.cond.exprs, LogicAnd)
	return l.cond
}

func (l *Logic) Or() *Condition {
	l.cond.exprs = append(l.cond.exprs, LogicOr)
	return l.cond
}

func (l *Logic) End() *Query {
	return l.end(false)
}

func (l *Logic) EndWithGroup() *Query {
	return l.end(true)
}

func (l *Logic) end(group bool) *Query {
	if len(l.cond.exprs) == 0 {
		return l.cond.query
	}

	expr := fmt.Sprintf("%s", strings.Join(l.cond.exprs, " "))
	if group {
		expr = "(" + expr + ")"
	}

	l.cond.query.exprs = append(l.cond.query.exprs, expr)
	l.cond.query.vals = append(l.cond.query.vals, l.cond.vals...)
	return l.cond.query
}
