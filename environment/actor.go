package environment

import types "github.com/codegamc/hollywood/types"

// Actor stores the stuff about the actor
type Actor struct {
	ActorID int64
	mail    chan types.HWType
	studio  *Studio
}

// GetMail returns the next item in ac actor's mailbox (or waits)
func (a Actor) GetMail() types.HWType {
	return <-a.mail
}

// SendMail uses the Actor's studio to send mail  to another actor
func (a Actor) SendMail(couplet []types.HWType) {
	//address := couplet[0]
	// message := couplet[1]
	a.studio.SendMail(couplet[0].(types.HWInt), couplet[1])
}
