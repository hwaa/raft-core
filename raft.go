package raft

import  pb "./raftpb"

type StateType int

const (
	Follower StateType = iota
	Candidate 
	Leader
	PreCandidate
)

type raftLog struct {

}

type Logger struct {

}

type Config struct {
	ID				uint64
	ElectionTick 	int
	HeartbeatTick 	int
	Storage 		pb.Storage
	Applied  		uint64
}

type raft struct {
	// raft node id
	id 			uint64
	// currentTerm, Initialized to 0
	term 		uint64
	// voted id, initialized to -1
	voteFor 	uint64
	// the log
	raftLog		[]*raftLog
	//cur state
	state 		StateType
	//isLearner,used for learning, is true if the local node is a learner
	isLearner	bool
	// the leader id
	lead		uint64
	// number of ticks since it reached last election when it is leader
	//	or candidate
	electionElapsed		int

	//number of ticks since it reached last heartbeatTimeout
	// only leader keeps heartbeatElapsed
	heartbeatElapsed	int

	checkQuorum			int
	preVote				int

	heartbeatTimeout	int
	electionTimeout		int

	//	randomizedElectionTimeout is a random number between
	//	[electiontimeout, 2 * electiontimeout - 1]. It gets reset
	//	when raft changes	its state to follower or candidate.
	randomizedElectionTimeout	int
	disableProposalForwarding	bool


	tick  				func()
	step 				stepFunc

	logger 				Logger

	pendingReadIndexMessages	[]pb.Message

}

type stepFunc func(r *raft, m pb.Message) error

func (ra *raft)Step(m stepFunc) error {
	return nil
}

func (ra *raft)becomeLeader()  {
	if ra.state != Candidate {
		return
	}

}

func(ra *raft)becomeFollower() {

}

func (ra *raft)becomeCandidate()  {

}

