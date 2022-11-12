package mongoutil

import (
	"go.mongodb.org/mongo-driver/bson"
)

type Stage interface {
	Base
}

type Match struct {
	Querys []Query
}

func (m *Match) GetBson() bson.M {
	r := bson.M{}
	for _, q := range m.Querys {
		for k, v := range q.GetBson() {
			r[k] = v
		}
	}
	return bson.M{MatchOperator: r}
}

type Lookup struct {
	From         string
	LocalField   string
	ForeignField string
	As           string
}

func (m *Lookup) GetBson() bson.M {
	r := bson.M{}
	if m.From != "" {
		r[LookFrom] = m.From
	}
	if m.LocalField != "" {
		r[LookLocalField] = m.LocalField
	}
	if m.ForeignField != "" {
		r[LookForeignField] = m.ForeignField
	}
	if m.As != "" {
		r[LookAs] = m.As
	}
	return bson.M{LookOperator: r}
}

type Unwind struct {
	Path                       string
	IncludeArrayIndex          string
	PreserveNullAndEmptyArrays bool
}

func (u *Unwind) GetBson() bson.M {
	result := bson.M{}
	if u.Path == "" {
		return bson.M{}
	}
	result[UnwindPath] = u.Path
	if u.IncludeArrayIndex != "" {
		result[UnwindIncludeArrayIndex] = u.IncludeArrayIndex
	}
	result[UnwindPreserveNullAndEmptyArrays] = u.PreserveNullAndEmptyArrays
	return bson.M{UnwindOperator: result}
}

type Sort struct {
	Fields map[string]Order
}

func (s *Sort) GetBson() bson.M {
	r := bson.M{}
	for k, o := range s.Fields {
		r[k] = o.GetValue()
	}
	return bson.M{SortOperator: r}
}

type Order struct {
	IsMeta    bool
	Value     int8
	MetaValue string
}

func (m *Order) GetValue() interface{} {
	if m.IsMeta {
		return bson.M{MetaOperatot: m.Value}
	}
	return m.Value
}

type Facet struct {
	Stages map[string][]Stage
}

func (f *Facet) GetBson() bson.M {
	ssr := bson.M{}
	for k, stages := range f.Stages {
		var ss []bson.M
		for _, s := range stages {
			ss = append(ss, s.GetBson())
		}
		ssr[k] = ss
	}
	return bson.M{FacetOperator: ssr}
}

type Project struct {
	Specifications []Specification
}

func (p *Project) GetBson() bson.M {
	r := bson.M{}
	for _, s := range p.Specifications {
		r[s.Key] = s.Value
	}
	return bson.M{ProjectOperator: r}
}

type Specification struct {
	Key   string
	Value interface{}
}
