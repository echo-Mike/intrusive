package fn

func I[E any](e E) E {
	return e
}

func Next[S ~[]E, E any](s S) func() E {
	i := 0
	return func() E {
		defer func() { i = i + 1 }()
		return s[i]
	}
}

func Append[E any](e E, l []E) []E {
	return append(l, e)
}

func Applicator[E any, U any](g func() E, f func(E) U) func() U {
	return func() U { return f(g()) }
}

func Apply[E any, U any](g func() E, f func(E) U, n int) []U {
	return Reduce(Applicator(g, f), Append[U], make([]U, 0, n), n)
}

func Map[S ~[]E, E any, U any](s S, f func(E) U) []U {
	return Apply(Next(s), f, len(s))
}

func Reducer[E any, U any](g func() E, f func(E, U) U, b U) func() U {
	return func() U {
		b = f(g(), b)
		return b
	}
}

func Reduce[E any, U any](g func() E, f func(E, U) U, b U, n int) (r U) {
	a := Reducer(g, f, b)
	for range n {
		r = a()
	}
	return
}

func Take[E any](g func() E, n int) []E {
	return Reduce(g, Append[E], make([]E, 0, n), n)
}

func Filter[S ~[]E, E any](s S, f func(E) bool) []E {
	return Reduce(Next(s), func(e E, c []E) []E {
		if f(e) {
			return append(c, e)
		}
		return c
	}, make([]E, 0, len(s)), len(s))
}

func Zipper[A any, B any](a func() A, b func() B) func() struct {
	first  A
	second B
} {
	return func() struct {
		first  A
		second B
	} {
		return struct {
			first  A
			second B
		}{a(), b()}
	}
}

func Zip[A any, B any](a func() A, b func() B, n int) []struct {
	first  A
	second B
} {
	return Reduce(Zipper(a, b), Append[struct {
		first  A
		second B
	}], make([]struct {
		first  A
		second B
	}, 0, n), n)
}

func AnyOf[S ~[]E, E any](s S, f func(E) bool) bool {
	return Reduce(Next(s), func(v E, r bool) bool { return r || f(v) }, false, len(s))
}

func AllOf[S ~[]E, E any](s S, f func(E) bool) bool {
	return Reduce(Next(s), func(v E, r bool) bool { return f(v) && r }, true, len(s))
}

func NoneOf[S ~[]E, E any](s S, f func(E) bool) bool {
	return AllOf(s, func(e E) bool { return !f(e) })
}
