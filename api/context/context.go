package context

type ctxKey int

const (
	// ParamKey is the request context key used to store params from the request path
	ParamKey ctxKey = iota + 1

	// RequestSchemeKey is the request context key used to store c.Scheme() created by
	// the PopulateRequestContext middleware.
	RequestSchemeKey
)
