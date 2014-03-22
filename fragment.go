package waitress

type Fragment struct {
	*Router
}

func NewFragment(ctx interface{}) *Fragment {
	return &Fragment{
		Router: NewRouter(ctx),
	}
}
