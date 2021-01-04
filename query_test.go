package gormqy

import (
	"testing"
)

func TestConcatCol(t *testing.T) {
	cases := []struct {
		sep  string
		cols []Column
		res  Column
	}{
		{sep: "|", cols: []Column{"name"}, res: "concat_ws('|', name)"},
		{sep: "|", cols: []Column{"name", "addr"}, res: "concat_ws('|', name, addr)"},
	}

	for _, c := range cases {
		t.Log(c)
		if res := ConcatCol(c.sep, c.cols...); res != c.res {
			t.Fatalf("expect: %s, got: %s", c.res, res)
		}
	}
}

func TestQueryOrder(t *testing.T) {
	type order struct {
		col    Column
		method OrderMethod
	}
	cases := []struct {
		orders []order
		res    string
	}{
		{
			orders: []order{
				{"id", OrderDesc},
			},
			res: "id DESC",
		},
		{
			orders: []order{
				{"id", OrderDesc},
				{"name", OrderASC},
			},
			res: "id DESC, name ASC",
		},
	}

	for _, c := range cases {
		t.Log(c)
		q := NewQuery()
		for _, o := range c.orders {
			q.AddOrder(o.col, o.method)
		}
		if res := q.Order(); res != c.res {
			t.Fatalf("expect: %s, got: %s", c.res, res)
		}
	}
}

func TestQueryCondition(t *testing.T) {
	type cond struct {
		col   Column
		op    Operator
		value interface{}
		logic Logic
	}
	cases := []struct {
		conds []cond
		expr  string
		vals  []interface{}
	}{
		{
			conds: []cond{
				{"id", OpEq, 1, LogicAnd},
			},
			expr: "id = ?",
			vals: []interface{}{1},
		},
		{
			conds: []cond{
				{"name", OpEq, "Tom", LogicAnd},
				{"age", OpLe, 20, LogicAnd},
			},
			expr: "name = ? AND age <= ?",
			vals: []interface{}{"Tom", 20},
		},
		{
			conds: []cond{
				{"name", OpPrefix, "T", LogicOr},
				{"name", OpSuffix, "m", LogicAnd},
			},
			expr: "name LIKE ? OR name LIKE ?",
			vals: []interface{}{"T%", "%m"},
		},
		{
			conds: []cond{
				{"name", OpContain, "x", LogicAnd},
				{"age", OpIn, [2]int{10, 20}, LogicAnd},
			},
			expr: "name LIKE ? AND age IN (?)",
			vals: []interface{}{"%x%", [2]int{10, 20}},
		},
	}

	for _, c := range cases {
		t.Log(c)
		q := NewQuery()
		for _, cond := range c.conds {
			q.AddCondition(cond.col, cond.op, cond.value, cond.logic)
		}
		expr, vals := q.Where()
		if expr != c.expr {
			t.Fatalf("expect expr: %s, got: %s", c.expr, expr)
		}
		for i := 0; i < len(c.vals); i++ {
			if vals[i] != c.vals[i] {
				t.Fatalf("expect vals: %v, got: %v", c.vals, vals)
			}
		}
	}
}

func TestPageQuery(t *testing.T) {
	cases := []struct {
		p      uint64
		s      uint64
		page   uint64
		size   uint64
		offset uint64
	}{
		{0, 0, 1, 10, 0},
		{2, 100, 2, 100, 100},
	}

	for _, c := range cases {
		t.Log(c)
		pq := NewPageQuery(c.p, c.s)
		if pq.Page() != c.page {
			t.Fatalf("expect page: %d, got: %d", c.page, pq.Page())
		}
		if pq.Size() != c.size {
			t.Fatalf("expect size: %d, got: %d", c.size, pq.Size())
		}
		if pq.Offset() != c.offset {
			t.Fatalf("expect offset: %d, got: %d", c.offset, pq.Offset())
		}
	}
}
