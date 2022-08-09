package gormqy

import "fmt"

type Condition struct {
	query *Query
	logic *Logic
	exprs []string
	vals  []interface{}
}

func NewCondition(q *Query) *Condition {
	logic := &Logic{}
	cond := &Condition{
		query: q,
	}
	logic.cond = cond
	cond.logic = logic
	return cond
}

func (c *Condition) Eq(column string, val interface{}) *Logic {
	return c.add(fmt.Sprintf("%s = ?", column), val)
}

func (c *Condition) NotEq(column string, val interface{}) *Logic {
	return c.add(fmt.Sprintf("%s != ?", column), val)
}

func (c *Condition) Ge(column string, val interface{}) *Logic {
	return c.add(fmt.Sprintf("%s >= ?", column), val)
}

func (c *Condition) Le(column string, val interface{}) *Logic {
	return c.add(fmt.Sprintf("%s <= ?", column), val)
}

func (c *Condition) Gt(column string, val interface{}) *Logic {
	return c.add(fmt.Sprintf("%s > ?", column), val)
}

func (c *Condition) Lt(column string, val interface{}) *Logic {
	return c.add(fmt.Sprintf("%s < ?", column), val)
}

func (c *Condition) Contain(column string, val string) *Logic {
	return c.add(fmt.Sprintf("%s LIKE ?", column), fmt.Sprintf("%%%s%%", val))
}

func (c *Condition) Prefix(column string, val string) *Logic {
	return c.add(fmt.Sprintf("%s LIKE ?", column), fmt.Sprintf("%s%%", val))
}

func (c *Condition) Suffix(column string, val string) *Logic {
	return c.add(fmt.Sprintf("%s LIKE ?", column), fmt.Sprintf("%%%s", val))
}

func (c *Condition) In(column string, val interface{}) *Logic {
	return c.add(fmt.Sprintf("%s IN (?)", column), val)
}

func (c *Condition) NotIn(column string, val interface{}) *Logic {
	return c.add(fmt.Sprintf("%s NOT IN (?)", column), val)
}

func (c *Condition) add(expr string, val interface{}) *Logic {
	c.exprs = append(c.exprs, expr)
	c.vals = append(c.vals, val)
	return c.logic
}
