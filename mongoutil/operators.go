package mongoutil

import (
	"go.mongodb.org/mongo-driver/bson"
)

type Base interface {
	GetBson() bson.M
}

type Number interface {
	~int8 | ~int | ~int16 | ~int32 | ~int64 | ~float64 | ~float32
}

type Expression interface {
	Base
}

type ABS[T Number | bson.M] struct {
	input T
}

func (a *ABS[T]) GetBson() bson.M {
	return bson.M{ABSOperator: a.input}
}

type Accumulator struct {
	Init           interface{}
	InitArgs       bson.A
	Accumulate     interface{}
	AccumulateArgs bson.A
	Merge          interface{}
	Finalize       interface{}
}

func (a *Accumulator) GetBson() bson.M {
	result := bson.M{}
	if a.Init != nil {
		result[AccumulatorInit] = a.Init
	}
	if a.InitArgs != nil {
		result[AccumulatorInitArgs] = a.InitArgs
	}
	if a.Accumulate != nil {
		result[AccumulatorAccumulate] = a.Accumulate
	}
	if a.AccumulateArgs != nil {
		result[AccumulatorAccumulateArgs] = a.AccumulateArgs
	}
	if a.Merge != nil {
		result[AccumulatorAccumulateArgs] = a.AccumulateArgs
	}
	if a.Finalize != nil {
		result[AccumulatorFinalize] = a.Finalize
	}
	result[AccumulatorLang] = "js"
	return bson.M{AccumulatorOperator: result}
}

type QueryKV struct {
	Left  string
	Right interface{}
}

type And struct {
	Expressions []Expression
}

func (a *And) GetBson() bson.M {
	var bes []bson.M
	for _, e := range a.Expressions {
		bes = append(bes, e.GetBson())
	}
	return bson.M{AndOperator: bes}
}

type Gte struct {
	
}