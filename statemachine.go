package raft

type StateMachine interface {
	Save() // save state machine to persistent storage
	Recoverty() // recover state machine from persistent storage
}