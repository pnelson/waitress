package waitress

import (
	"github.com/pnelson/waitress/middleware"
)

type fragment interface {
	Fragment() *Fragment
}

type Fragment struct {
	*middleware.Builder
	*Router
}

func NewFragment() *Fragment {
	return &Fragment{
		Builder: &middleware.Builder{},
		Router:  NewRouter(),
	}
}
