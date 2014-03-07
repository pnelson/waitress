package router

type Router struct {
	Converters map[string]NewConverter
}

func New() *Router {
	return &Router{
		Converters: map[string]NewConverter{
			"default": NewStringConverter,
			"string":  NewStringConverter,
			"path":    NewPathConverter,
			"any":     NewAnyConverter,
			"int":     NewIntConverter,
		},
	}
}
