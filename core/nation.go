package core

// Nation represents a nation in the game.
type Nation interface {
	ID() NationID
	Name() string
}

// NationID is a unique identifier for a nation.
type NationID string

// MyNation represents the player's nation.
type MyNation struct {
	id   NationID
	name string
}

// NewMyNation creates a new MyNation instance.
func NewMyNation(id NationID, name string) *MyNation {
	return &MyNation{
		id:   id,
		name: name,
	}
}

// ID returns the nation ID.
func (n *MyNation) ID() NationID {
	return n.id
}

// Name returns the nation name.
func (n *MyNation) Name() string {
	return n.name
}

// OtherNation represents an NPC nation.
type OtherNation struct {
	id   NationID
	name string
}

// NewOtherNation creates a new OtherNation instance.
func NewOtherNation(id NationID, name string) *OtherNation {
	return &OtherNation{
		id:   id,
		name: name,
	}
}

// ID returns the nation ID.
func (n *OtherNation) ID() NationID {
	return n.id
}

// Name returns the nation name.
func (n *OtherNation) Name() string {
	return n.name
}
