package mongoutil

const (
	ABSOperator = "$abs"

	FacetOperator   = "$facet"
	ProjectOperator = "$project"
	SortOperator    = "$sort"
	
	GteOperator = "$gte"
	LteOperator = "$lte"
	AndOperator = "$and"
)

const (
	MetaOperatot  = "$meta"
	MetaTextScore = "textScore"
	MetaIndexKey  = "indexKey"
)

const (
	AccumulatorOperator       = "$accumulator"
	AccumulatorInit           = "init"
	AccumulatorInitArgs       = "initArgs"
	AccumulatorAccumulate     = "accumulate"
	AccumulatorAccumulateArgs = "accumulateArgs"
	AccumulatorMerge          = "merge"
	AccumulatorFinalize       = "finalize"
	AccumulatorLang           = "lang"
)

const (
	MatchOperator = "$match"
)

const (
	LookOperator     = "$lookup"
	LookFrom         = "from"
	LookLocalField   = "localField"
	LookForeignField = "foreignField"
	LookAs           = "as"
)

const (
	UnwindOperator                   = "$unwind"
	UnwindPath                       = "path"
	UnwindIncludeArrayIndex          = "includeArrayIndex"
	UnwindPreserveNullAndEmptyArrays = "preserveNullAndEmptyArrays"
)

const (
	InOperator = "$in"
)
