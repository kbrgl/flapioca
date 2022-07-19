package internal

type Obstacles []*Obstacle

func NewObstacles(obstacles ...*Obstacle) Obstacles {
	return Obstacles(obstacles)
}

func (ol *Obstacles) Remove() {
	(*ol)[0] = nil
	*ol = (*ol)[1:]
}

func (ol *Obstacles) Add(obst *Obstacle) {
	*ol = append(*ol, obst)
}

func (ol Obstacles) Index(i int) *Obstacle {
	if i < 0 || i >= len(ol) {
		return nil
	}
	return (ol)[i]
}

func (ol Obstacles) Rightmost() *Obstacle {
	if len(ol) == 0 {
		return nil
	}
	return (ol)[len(ol)-1]
}
