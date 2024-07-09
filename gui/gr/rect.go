package gr

type Rect struct {
	X0, Y0, X1, Y1 int
}

func (r Rect) Contains(o Rect) bool {
	return (r.X0 <= o.X0) && (r.X1 >= o.X1) && (r.Y0 <= o.Y0) && (r.Y1 >= o.Y1)
}
