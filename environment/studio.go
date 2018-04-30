package environment

import types "github.com/codegamc/hollywood/types"

// its because actors congregate in the studio... get it?

/*

Idea of how this should work:

* The actor, on creation, gets a few functions added to its env by the studio
Those functions are:
	1. (message/write targetAddressAsInt myMessage[hwtype]) => returns true if address is
		valid and in studio
	2. (message/read) => returns a message (hwType) if there is one waiting, this is a
		blocking op. so it waits for that.

	To accomplish this, when an actor + new envi is created, the studio will add functions
		to the envi, and seed the cooresponding channels using function closures


*/

// MakeStudio creates a new studio
func MakeStudio() *Studio {

	nextID := int64(0)
	m := make(map[int64](chan types.HWType))

	// A gorouting should be created here that listens for a message on any of the mailboxes?

	return &Studio{nextID: nextID, actorMap: m}
}

// Studio is an object that is used to manage actors
type Studio struct {
	nextID   int64
	actorMap map[int64](chan types.HWType)
	//studioMail chan types.HWType
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
	mailbox := make(chan types.HWType)

	return Actor{ActorID: actorID, mail: mailbox, studio: s}
}

// SendMail is how the studio sends mail
func (s *Studio) SendMail(address types.HWInt, message types.HWType) {
	s.actorMap[address.Val] <- message

}
