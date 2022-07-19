package internal

type Obstacle struct {
	// Aperture is the size of the opening in the obstacle.
	// Must be odd. Even values will be rounded up to the next odd value.
	Aperture int
	*Location
}

func NewObstacle(aperture int, l *Location) *Obstacle {
	return &Obstacle{
		Aperture: aperture,
		Location: l,
	}
}

func (o *Obstacle) Collides(l Location) bool {
	sameColumn := o.X == l.X
	radius := o.Aperture / 2
	inAperture := abs(o.Y-l.Y) <= radius
	return sameColumn && !inAperture
}
