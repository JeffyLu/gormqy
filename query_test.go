package gormqy

import (
	"testing"
)

func TestOrder(t *testing.T) {
	type order struct {
		column string
		method string
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
		q := New()
		for _, o := range c.orders {
			if o.method == OrderDesc {
				q.DESC(o.column)
			} else {
				q.ASC(o.column)
			}
		}
		if res := q.Order(); res != c.res {
			t.Fatalf("expect: %s, got: %s", c.res, res)
		}
	}
}

func TestCondition(t *testing.T) {
	{
		exprExp := "name = ? AND age > ?"
		t.Log(exprExp)

		expr, vals := New().
			Condition().
			Eq("name", "Tom").And().Gt("age", 10).End().
			Where()
		if expr != exprExp {
			t.Fatalf("expect: %s, got: %s", exprExp, expr)
		}
		if len(vals) != 2 || vals[0] != "Tom" || vals[1] != 10 {
			t.Fatalf("expect vals: [Tom 10], got: %+v", vals)
		}
	}
	{
		exprExp := "(name = ? AND age > ?) OR (name = ? AND age <= ?)"
		t.Log(exprExp)

		expr, vals := New().
			Condition().
			Eq("name", "Tom").And().Gt("age", 10).EndWithGroup().
			OrCondition().
			Eq("name", "Sam").And().Le("age", 5).EndWithGroup().
			Where()
		if expr != exprExp {
			t.Fatalf("expect: %s, got: %s", exprExp, expr)
		}
		if len(vals) != 4 || vals[0] != "Tom" || vals[1] != 10 || vals[2] != "Sam" || vals[3] != 5 {
			t.Fatalf("expect vals: [Tom 10 Sam 5], got: %+v", vals)
		}
	}
	{
		exprExp := "(name LIKE ? OR name LIKE ?) AND age >= ? AND point < ?"
		t.Log(exprExp)

		expr, vals := New().
			Condition().
			Prefix("name", "T").Or().Suffix("name", "d").EndWithGroup().
			AndCondition().
			Ge("age", 10).And().Lt("point", 70).End().
			Where()
		if expr != exprExp {
			t.Fatalf("expect: %s, got: %s", exprExp, expr)
		}
		if len(vals) != 4 || vals[0] != "T%" || vals[1] != "%d" || vals[2] != 10 || vals[3] != 70 {
			t.Fatalf("expect vals: [T%% %%d 10 70], got: %+v", vals)
		}
	}
	{
		exprExp := "name NOT IN (?)"
		t.Log(exprExp)

		expr, vals := New().
			Condition().
			NotIn("name", []string{"Tom", "Sam"}).End().
			Where()
		if expr != exprExp {
			t.Fatalf("expect: %s, got: %s", exprExp, expr)
		}
		if len(vals) != 1 {
			t.Fatalf("expect vals: [[Tom Sam]], got: %+v", vals)
		}
		names, ok := vals[0].([]string)
		if !ok {
			t.Fatalf("expect val type: []string, got: unknown")
		}
		if len(names) != 2 || names[0] != "Tom" || names[1] != "Sam" {
			t.Fatalf("expect names: [Tom Sam], got: %+v", names)
		}
	}
	{

		exprExp := "((a = ? AND b = ?) OR (c = ? AND d = ?)) AND e = ?"
		t.Log(exprExp)

		expr, vals := New().
			Condition().
			Eq("a", 1).And().
			Eq("b", 1).EndWithGroup().
			OrCondition().
			Eq("c", 1).And().
			Eq("d", 1).EndWithGroup().
			AndConditionAfterGroup().
			Eq("e", 1).End().
			Where()
		if expr != exprExp {
			t.Fatalf("expect: %s, got: %s", exprExp, expr)
		}
		if len(vals) != 5 || vals[0] != 1 || vals[1] != 1 || vals[2] != 1 || vals[3] != 1 || vals[4] != 1 {
			t.Fatalf("expect vals: [1 1 1 1 1], got: %+v", vals)
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
		pq := New().ToPageQuery(c.p, c.s)
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
