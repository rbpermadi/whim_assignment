package request

type CurrencyParameter struct {
	Limit  int
	Offset int
	Query  string
}

type ConversionParameter struct {
	Limit          int
	Offset         int
	CurrencyIDFrom int64
	CurrencyIDTo   int64
}
