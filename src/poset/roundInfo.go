package poset

import (
	"github.com/golang/protobuf/proto"
)

type pendingRound struct {
	Index   int64
	Decided bool
}

// RoundCreated wrapper for protobuf created round event messages
type RoundCreated struct {
	Message RoundCreatedMessage
}

// NewRoundInfo creates a new round info struct
func NewRoundCreated() *RoundCreated {
	return &RoundCreated{
		Message: RoundCreatedMessage{
			Events: make(map[string]*RoundEvent),
		},
	}
}

// RoundReceived constructor
func NewRoundReceived() *RoundReceived {
	return &RoundReceived{
		Rounds: []string{},
	}
}

// AddEvent add event to round info (optionally set clotho)
func (r *RoundCreated) AddEvent(x string, clotho bool) {
	_, ok := r.Message.Events[x]
	if !ok {
		r.Message.Events[x] = &RoundEvent{
			Clotho: clotho,
		}
	}
}

// SetConsensusEvent set an event as a consensus event
func (r *RoundCreated) SetConsensusEvent(x string) {
	e, ok := r.Message.Events[x]
	if !ok {
		e = &RoundEvent{}
	}
	e.Consensus = true
	r.Message.Events[x] = e
}

// SetRoundReceived set the received round for the given event
func (r *RoundCreated) SetRoundReceived(x string, round int64) {
	e, ok := r.Message.Events[x]

	if !ok {
		return
	}

	e.RoundReceived = round

	r.Message.Events[x] = e
}

// SetAtropos sets whether the given event is Atropos, otherwise it is Clotho when not found
func (r *RoundCreated) SetAtropos(x string, f bool) {
	e, ok := r.Message.Events[x]
	if !ok {
		e = &RoundEvent{
			Clotho: true,
		}
	}
	if f {
		e.Atropos = Trilean_TRUE
	} else {
		e.Atropos = Trilean_FALSE
	}
	r.Message.Events[x] = e
}

// ClothoDecided return true if no clothos' fame is left undefined
func (r *RoundCreated) ClothoDecided() bool {
	for _, e := range r.Message.Events {
		if e.Clotho && e.Atropos == Trilean_UNDEFINED {
			return false
		}
	}
	return true
}

// Clotho return clothos
func (r *RoundCreated) Clotho() []string {
	var res []string
	for x, e := range r.Message.Events {
		if e.Clotho {
			res = append(res, x)
		}
	}
	return res
}

// RoundEvents returns all non-consensus events for the created round 
func (r *RoundCreated) RoundEvents() []string {
	var res []string
	for x, e := range r.Message.Events {
		if !e.Consensus {
			res = append(res, x)
		}
	}
	return res
}

// ConsensusEvents returns all consensus events for the created round
func (r *RoundCreated) ConsensusEvents() []string {
	var res []string
	for x, e := range r.Message.Events {
		if e.Consensus {
			res = append(res, x)
		}
	}
	return res
}

// Atropos return Atropos
func (r *RoundCreated) Atropos() []string {
	var res []string
	for x, e := range r.Message.Events {
		if e.Clotho && e.Atropos == Trilean_TRUE {
			res = append(res, x)
		}
	}
	return res
}

// IsDecided checks if the event is a decided clotho
func (r *RoundCreated) IsDecided(clotho string) bool {
	w, ok := r.Message.Events[clotho]
	return ok && w.Clotho && w.Atropos != Trilean_UNDEFINED
}

// ProtoMarshal marshals the created round to protobuf
func (r *RoundCreated) ProtoMarshal() ([]byte, error) {
	var bf proto.Buffer
	bf.SetDeterministic(true)
	if err := bf.Marshal(&r.Message); err != nil {
		return nil, err
	}
	return bf.Bytes(), nil
}

// ProtoMarshal serialises the received round using protobuf
func (r *RoundReceived) ProtoMarshal() ([]byte, error) {
	var bf proto.Buffer
	bf.SetDeterministic(true)
	if err := bf.Marshal(r); err != nil {
		return nil, err
	}
	return bf.Bytes(), nil
}

func (r *RoundCreated) ProtoUnmarshal(data []byte) error {
	return proto.Unmarshal(data, &r.Message)
}

// ProtoUnmarshal de-serialises RoundReceived using protobuf
func (r *RoundReceived) ProtoUnmarshal(data []byte) error {
	return proto.Unmarshal(data, r)
}

// IsQueued returns whether the RoundCreated is queued for processing in PendingRounds
func (r *RoundCreated) IsQueued() bool {
	return r.Message.Queued
}

// Equals compares round events for equality
func (re *RoundEvent) Equals(that *RoundEvent) bool {
	return re.Consensus == that.Consensus &&
		re.Clotho == that.Clotho &&
		re.Atropos == that.Atropos &&
		re.RoundReceived == that.RoundReceived
}

// EqualsMapStringRoundEvent compares a map string of round eventss for equality
func EqualsMapStringRoundEvent(this map[string]*RoundEvent, that map[string]*RoundEvent) bool {
	if len(this) != len(that) {
		return false
	}
	for k, v := range this {
		v2, ok := that[k]
		if !ok || !v2.Equals(v) {
			return false
		}
	}
	return true
}

// Equals compares two round created structs for equality
func (r *RoundCreated) Equals(that *RoundCreated) bool {
	return r.Message.Queued == that.Message.Queued &&
		EqualsMapStringRoundEvent(r.Message.Events, that.Message.Events)
}
