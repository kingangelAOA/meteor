package mongoutil

import (
	"go.mongodb.org/mongo-driver/bson"
)

type Query interface {
	Base
}

type EqualityCondition struct {
	Key   string
	Value interface{}
}

func (ec *EqualityCondition) GetBson() bson.M {
	return bson.M{ec.Key: ec.Value}
}

type ComparisonQuery struct {
	Name       string
	Comparisons []Comparison
}

func (cq *ComparisonQuery) GetBson() bson.M {
	r := bson.M{}
	for _, c := range cq.Comparisons {
		r[c.Operator] = c.Value
	}
	return bson.M{cq.Name: r}
}

type Comparison struct {
	Operator string
	Value      interface{}
}

func (c *Comparison) GetBson() bson.M {
	return bson.M{c.Operator: c.Value}
}


