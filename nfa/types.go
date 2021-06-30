package nfa

type NFARule struct {
	ComparisonOperator string   `json:"comparisonOperator"`
	Key                string   `json:"key"`
	Value              []string `json:"value"`
}

type NFAFilterItem struct {
	Condition string    `json:"condition"`
	Rules     []NFARule `json:"rules"`
}

type NFAParameter struct {
	AggregateFunction string        `json:"aggregateFunction"`
	Limit             int           `json:"limit"`
	OrderBy           string        `json:"orderby"`
	PageSize          int           `json:"pageSize"`
	AggregateColumn   string        `json:"aggregateColumn"`
	Devices           []int         `json:"devices"`
	Locations         []int         `json:"locations"`
	Filters           NFAFilterItem `json:"filters"`
	GroupBy           []string      `json:"groupby"`
	Order             string        `json:"order"`
	GroupByDstPrefix  int           `json:"groupByDstPrefix"`
	GroupBySrcPrefix  int           `json:"groupBySrcPrefix"`
	RateUnit          string        `json:"rateUnit"`
}

type NFAFilter struct {
	Id          int          `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Parameters  NFAParameter `json:"parameters"`
	Report      string       `json:"report"`
	Shared      bool         `json:"shared"`
	Owner       string       `json:"owner"`
}
