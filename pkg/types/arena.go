package types

type Arena struct {
	ID      int32
	Name    string
	Mobs    map[int32]*Mob
	Players map[int32]*Player
}

func NewArena(id int32, name string) *Arena {
	return &Arena{
		ID:      id,
		Name:    name,
		Mobs:    make(map[int32]*Mob),
		Players: make(map[int32]*Player),
	}
}

func (a *Arena) Update() {
	// Update state in the arena
}
