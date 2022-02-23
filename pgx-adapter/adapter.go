package main

import (
	"errors"
	"fmt"
	responsev1 "github.com/cerbos/cerbos/api/genpb/cerbos/response/v1"
	"github.com/iancoleman/strcase"
	"strings"
)

var toSqlOp = map[string]string{
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

var toSqlField = map[string]string{}

var ExpressionExpectedError = errors.New("expected expression")

type BuildPredicateType func(e *responsev1.ResourcesQueryPlanResponse_Expression_Operand_Expression) (where string, args []interface{}, err error)

func (t BuildPredicateType) BuildPredicate(e *responsev1.ResourcesQueryPlanResponse_Expression_Operand_Expression) (where string, args []interface{}, err error) {
	return t(e)
}

func buildPredicateImpl(e *responsev1.ResourcesQueryPlanResponse_Expression_Operand_Expression, b *strings.Builder, args *[]interface{}) (err error) {
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
			if oe, ok := o.GetNode().(*responsev1.ResourcesQueryPlanResponse_Expression_Operand_Expression); !ok {
				return ExpressionExpectedError
			} else {
				err = buildPredicateImpl(oe, b, args)
				if err != nil {
					return err
				}
			}
		}
		b.WriteRune(')')
		return nil
	case "not":
		o := e.Expression.Operands[0]
		b.WriteRune('(')
		b.WriteString("NOT ")
		if oe, ok := o.GetNode().(*responsev1.ResourcesQueryPlanResponse_Expression_Operand_Expression); !ok {
			return ExpressionExpectedError
		} else {
			err = buildPredicateImpl(oe, b, args)
			if err != nil {
				return err
			}
			b.WriteRune(')')
			return nil
		}
	default:
		if len(e.Expression.Operands) != 2 {
			return fmt.Errorf("expected a binary operation: op = %q, # of operands = %d", e.Expression.Operator, len(e.Expression.Operands))
		}
		op, ok := toSqlOp[e.Expression.Operator]
		if !ok {
			return fmt.Errorf("unsupported operation %q", e.Expression.Operator)
		}
		b.WriteRune('(')
		for i, operand := range e.Expression.Operands {
			switch eo := operand.Node.(type) {
			case *responsev1.ResourcesQueryPlanResponse_Expression_Operand_Expression:
				err = buildPredicateImpl(eo, b, args)
				if err != nil {
					return err
				}
			case *responsev1.ResourcesQueryPlanResponse_Expression_Operand_Variable:
				b.WriteRune('"')
				b.WriteString(getFieldName(eo.Variable))
				b.WriteRune('"')
			case *responsev1.ResourcesQueryPlanResponse_Expression_Operand_Value:
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

func BuildPredicate(e *responsev1.ResourcesQueryPlanResponse_Expression_Operand_Expression) (where string, args []interface{}, err error) {
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

	if s, ok := toSqlField[name]; ok {
		return s
	}

	return strcase.ToSnake(name)
}
