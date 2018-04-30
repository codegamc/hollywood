package environment

// its because actors congregate in the studio... get it?

// MakeStudio creates a new studio
func MakeStudio() *Studio {

	nextID := int64(0)
	return &Studio{nextID: nextID}
}

// Studio is an object that is used to manage actors
type Studio struct {
	nextID int64
}

// getNextID returns the next valid potential ID for an actor
func (s *Studio) getNextID() int64 {
	id := s.nextID
	s.nextID = id + 1
	return id
}

// NewActor creates a new actor
func (s *Studio) NewActor() Actor {
	actorID := s.getNextID()
	return Actor{ActorID: actorID, studio: s}
}
