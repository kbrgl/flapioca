package internal

type Obstacle struct {
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
	sameColumn := o.x == l.x
	radius := o.Aperture / 2
	inAperture := abs(o.y-l.y) <= radius
	return sameColumn && !inAperture
}
