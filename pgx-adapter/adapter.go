// Copyright 2021-2022 Zenauth Ltd.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"errors"
	"fmt"
	"strings"

	enginev1 "github.com/cerbos/cerbos/api/genpb/cerbos/engine/v1"
	"github.com/iancoleman/strcase"
)

var toSQLOp = map[string]string{
	"eq":   "=",
	"ne":   "<>",
	"lt":   "<",
	"lte":  "<=",
	"gt":   ">",
	"gte":  ">=",
	"add":  "+",
	"sub":  "-",
	"mult": "*",
	"div":  "/",
	"mod":  "%",
	"in":   "IN",
}

var toSQLField = map[string]string{}

var ErrExpressionExpected = errors.New("expected expression")

type filterOpExpression = enginev1.PlanResourcesFilter_Expression_Operand_Expression
type filterOpValue = enginev1.PlanResourcesFilter_Expression_Operand_Value
type filterOpVariable = enginev1.PlanResourcesFilter_Expression_Operand_Variable
type filterOp = enginev1.PlanResourcesFilter_Expression_Operand

type BuildPredicateType func(e *filterOpExpression) (where string, args []interface{}, err error)

func (t BuildPredicateType) BuildPredicate(e *filterOpExpression) (where string, args []interface{}, err error) {
	return t(e)
}

func buildPredicateImpl(e *filterOpExpression, b *strings.Builder, args *[]interface{}) (err error) {
	switch e.Expression.Operator {
	case "or", "and":
		b.WriteRune('(')
		op := strings.ToUpper(e.Expression.Operator)
		n := len(e.Expression.Operands)
		for i, o := range e.Expression.Operands {
			if i > 0 && n > 1 {
				b.WriteRune(' ')
				b.WriteString(op)
				b.WriteRune(' ')
			}
			if oe, ok := o.GetNode().(*filterOpExpression); ok {
				err = buildPredicateImpl(oe, b, args)
				if err != nil {
					return err
				}
			} else {
				return ErrExpressionExpected
			}
		}
		b.WriteRune(')')
		return nil
	case "not":
		o := e.Expression.Operands[0]
		b.WriteRune('(')
		b.WriteString("NOT ")
		if oe, ok := o.GetNode().(*filterOpExpression); ok {
			err = buildPredicateImpl(oe, b, args)
			if err != nil {
				return err
			}
			b.WriteRune(')')
			return nil
		}
		return ErrExpressionExpected
	default:
		if len(e.Expression.Operands) != 2 { //nolint:gomnd
			return fmt.Errorf("expected a binary operation: op = %q, # of operands = %d", e.Expression.Operator, len(e.Expression.Operands))
		}
		op, ok := toSQLOp[e.Expression.Operator]
		if !ok {
			return fmt.Errorf("unsupported operation %q", e.Expression.Operator)
		}
		b.WriteRune('(')
		for i, operand := range e.Expression.Operands {
			switch eo := operand.Node.(type) {
			case *filterOpExpression:
				err = buildPredicateImpl(eo, b, args)
				if err != nil {
					return err
				}
			case *filterOpVariable:
				b.WriteRune('"')
				b.WriteString(getFieldName(eo.Variable))
				b.WriteRune('"')
			case *filterOpValue:
				*args = append(*args, eo.Value.AsInterface())
				b.WriteString(fmt.Sprintf("$%d", len(*args)))
			}
			if i == 0 {
				b.WriteRune(' ')
				b.WriteString(op)
				b.WriteRune(' ')
			}
		}
		b.WriteRune(')')
	}

	return nil
}

func BuildPredicate(e *filterOpExpression) (where string, args []interface{}, err error) {
	if e == nil {
		return "", nil, nil
	}
	b := new(strings.Builder)
	err = buildPredicateImpl(e, b, &args)
	where = b.String()
	n := len(where)
	if n > 0 && where[0] == '(' && where[n-1] == ')' {
		where = where[1 : n-1]
	}
	return where, args, err
}

func getFieldName(name string) string {
	name = strings.TrimPrefix(name, "request.resource.attr.")
	name = strings.TrimPrefix(name, "R.attr.")

	if s, ok := toSQLField[name]; ok {
		return s
	}

	return strcase.ToSnake(name)
}
