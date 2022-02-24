// Copyright 2021-2022 Zenauth Ltd.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"errors"
	"fmt"
	"strings"

	"entgo.io/ent/dialect/sql"
	responsev1 "github.com/cerbos/cerbos/api/genpb/cerbos/response/v1"
	"github.com/iancoleman/strcase"
)

var toSQLOp = map[string]sql.Op{
	"eq":   sql.OpEQ,
	"ne":   sql.OpNEQ,
	"lt":   sql.OpLT,
	"lte":  sql.OpLTE,
	"gt":   sql.OpGT,
	"gte":  sql.OpGTE,
	"add":  sql.OpAdd,
	"sub":  sql.OpSub,
	"mult": sql.OpMul,
	"div":  sql.OpDiv,
	"mod":  sql.OpMod,
	"in":   sql.OpIn,
}

var toEntField = map[string]string{
	"ownerId": "user_contacts",
}

var ErrExpressionExpected = errors.New("expected expression")

type BuildPredicateType func(e *responsev1.ResourcesQueryPlanResponse_Expression_Operand_Expression) (p *sql.Predicate, err error)

func (t BuildPredicateType) BuildPredicate(e *responsev1.ResourcesQueryPlanResponse_Expression_Operand_Expression) (p *sql.Predicate, err error) {
	return t(e)
}

func BuildPredicate(e *responsev1.ResourcesQueryPlanResponse_Expression_Operand_Expression) (p *sql.Predicate, err error) {
	if e == nil {
		return nil, nil
	}
	switch e.Expression.Operator {
	case "or", "and":
		ps := make([]*sql.Predicate, len(e.Expression.Operands))
		for i, o := range e.Expression.Operands {
			if oe, ok := o.GetNode().(*responsev1.ResourcesQueryPlanResponse_Expression_Operand_Expression); ok {
				ps[i], err = BuildPredicate(oe)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, ErrExpressionExpected
			}
		}
		if e.Expression.Operator == "or" {
			return sql.Or(ps...), nil
		}
		return sql.And(ps...), nil
	case "not":
		o := e.Expression.Operands[0]
		if oe, ok := o.GetNode().(*responsev1.ResourcesQueryPlanResponse_Expression_Operand_Expression); ok {
			p, err = BuildPredicate(oe)
			if err != nil {
				return nil, err
			}
			return sql.Not(p), nil
		}
		return nil, ErrExpressionExpected
	default:
		const numOperands = 2
		if len(e.Expression.Operands) != numOperands {
			return nil, fmt.Errorf("expected a binary operation: op = %q, # of operands = %d", e.Expression.Operator, len(e.Expression.Operands))
		}
		op, ok := toSQLOp[e.Expression.Operator]
		if !ok {
			return nil, fmt.Errorf("unsupported operation %q", e.Expression.Operator)
		}
		var args [2]func(builder *sql.Builder) *sql.Builder
		for i, v := range e.Expression.Operands {
			args[i], err = newBuilder(v)
			if err != nil {
				return nil, err
			}
		}
		return sql.P().Append(func(b *sql.Builder) {
			args[0](b)
			b.WriteOp(op)
			args[1](b)
		}), nil
	}
}

func newBuilder(operand *responsev1.ResourcesQueryPlanResponse_Expression_Operand) (func(*sql.Builder) *sql.Builder, error) {
	switch e := operand.Node.(type) {
	case *responsev1.ResourcesQueryPlanResponse_Expression_Operand_Expression:
		p, err := BuildPredicate(e)
		if err != nil {
			return nil, err
		}
		return func(b *sql.Builder) *sql.Builder {
			return b.Join(p)
		}, nil
	case *responsev1.ResourcesQueryPlanResponse_Expression_Operand_Variable:
		return func(b *sql.Builder) *sql.Builder {
			return b.Ident(getFieldName(e.Variable))
		}, nil
	case *responsev1.ResourcesQueryPlanResponse_Expression_Operand_Value:
		return func(b *sql.Builder) *sql.Builder {
			return b.Arg(e.Value.AsInterface())
		}, nil
	}
	return nil, errors.New("unknown Node type")
}

func getFieldName(name string) string {
	name = strings.TrimPrefix(name, "request.resource.attr.")
	name = strings.TrimPrefix(name, "R.attr.")

	if s, ok := toEntField[name]; ok {
		return s
	}

	return strcase.ToSnake(name)
}
