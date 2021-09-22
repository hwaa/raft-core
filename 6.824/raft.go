package raft

//
// this is an outline of the API that raft must expose to
// the service (or tester). see comments below for
// each of these functions for more details.
//
// rf = Make(...)
//   create a new Raft server.
// rf.Start(command interface{}) (index, term, isleader)
//   start agreement on a new log entry
// rf.GetState() (term, isLeader)
//   ask a Raft for its current term, and whether it thinks it is leader
// ApplyMsg
//   each time a new entry is committed to the log, each Raft peer
//   should send an ApplyMsg to the service (or tester)
//   in the same server.
//

import (
	"6.824/labgob"
	"bytes"
	"log"
	"math/rand"
	"sort"
	"time"

	"6.824/labrpc"
	"sync"
	"sync/atomic"
)

// ApplyMsg
// as each Raft peer becomes aware that successive log entries are
// committed, the peer should send an ApplyMsg to the service (or
// tester) on the same server, via the applyCh passed to Make(). set
// CommandValid to true to indicate that the ApplyMsg contains a newly
// committed log entry.
//
// in part 2D you'll want to send other kinds of messages (e.g.,
// snapshots) on the applyCh, but set CommandValid to false for these
// other uses.
//
type ApplyMsg struct {
	CommandValid bool
	Command      interface{}
	CommandIndex int

	// For 2D:
	SnapshotValid bool
	Snapshot      []byte
	SnapshotTerm  int
	SnapshotIndex int
}

// Raft
// A Go object implementing a single Raft peer.
//
type state int

const (
	Follower	state = iota
	Candidate
	Leader
)
const None int = -1

type Log struct {
	Term	int
	Command interface{}
}

type Raft struct {
	mu        sync.RWMutex          // Lock to protect shared access to this peer's state
	peers     []*labrpc.ClientEnd // RPC end points of all peers
	persister *Persister          // Object to hold this peer's persisted state
	me        int                 // this peer's index into peers[]
	dead      int32               // set by Kill()

	state	  state				//state of cur raft node
	currentTerm	int				//latest term has seen(initialized to 0, increases monotonically)"
	voteFor		int				//CandidateId that received vote in current Term(or null if none)"
	log			[]Log
	//for log compaction
	lastIncludedIndex	int		//the snap shot replaces all entries up through and including this index"
	lastIncludedTerm	int		//term of lastIncludedIndex"
	//Volatile state on all servers:
	commitIndex int		//index of highest log entry known to be committed(initialized to 0, increases monotonically)"
	lastApplied int		//index of highest log entry known to be applied to state machine(initialized to 0, increases monotonically)"
	//Volatile state on Leaders:(reinitialized after election)
	nextIndex	[]int	//for each server, index of the next log entry to be sent to that server"
	matchIndex	[]int	//for each server, index of the highest log entry known to be replicated to that server(initialized to 0, im)"
	//channel
	applyCh		chan ApplyMsg	// from Make()
	killCh		chan bool		//for Kill()
	//handle rpc
	voteCh		chan bool	//send agree or disagree
	appendLogCh	chan bool
}

// GetState return currentTerm and whether this server
// believes it is the leader.
func (rf *Raft) GetState() (int, bool) {

	var term int
	var isleader bool
	// Your code here (2A).
	rf.mu.RLock()
	defer rf.mu.RUnlock()
	term = rf.currentTerm
	isleader = rf.state == Leader
	return term, isleader
}

//
// save Raft's persistent state to stable storage,
// where it can later be retrieved after a crash and restart.
// see paper's Figure 2 for a description of what should be persistent.
//
func (rf *Raft) persist() {
	rf.persister.SaveRaftState(rf.encodeRaftState())
}

func(rf *Raft)encodeRaftState() []byte {
	w := new(bytes.Buffer)
	e := labgob.NewEncoder(w)
	_ = e.Encode(rf.currentTerm)
	_ = e.Encode(rf.voteFor)
	_ = e.Encode(rf.log)
	_ = e.Encode(rf.lastIncludedIndex)
	_ = e.Encode(rf.lastIncludedTerm)
	return w.Bytes()
}

func(rf *Raft) persistWithSnapshot(snapshot []byte) {
	rf.persister.SaveStateAndSnapshot(rf.encodeRaftState(), snapshot)
}

//
// restore previously persisted state.
//
func (rf *Raft) readPersist(data []byte) {
	if data == nil || len(data) < 1 { // bootstrap without any state?
		return
	}
	r := bytes.NewBuffer(data)
	d := labgob.NewDecoder(r)
	var currentTerm int
	var voteFor 	int
	var clog 		[]Log
	var lastIncludedIndex 	int
	var lastIncludedTerm	int
	if d.Decode(&currentTerm) != nil || d.Decode(&voteFor) != nil || d.Decode(&clog) != nil || d.Decode(&lastIncludedIndex) != nil || d.Decode(&lastIncludedTerm) != nil {
		log.Fatalf("readPersist ERROR at server %v", rf.me)
	} else {
		rf.mu.Lock()
		rf.currentTerm, rf.voteFor = currentTerm, voteFor
		rf.lastIncludedTerm, rf.lastIncludedIndex = lastIncludedTerm, lastIncludedIndex
		rf.log = clog
		rf.commitIndex, rf.lastApplied = rf.lastIncludedIndex, rf.lastIncludedIndex
		rf.mu.Unlock()
	}
}


// CondInstallSnapshot
// A service wants to switch to snapshot.  Only do so if Raft hasn't
// have more recent info since it communicate the snapshot on applyCh.
//
func (rf *Raft) CondInstallSnapshot(lastIncludedTerm int, lastIncludedIndex int, snapshot []byte) bool {

	// Your code here (2D).

	return true
}

// Snapshot the service says it has created a snapshot that has
// all info up to and including index. this means the
// service no longer needs the log through (and including)
// that index. Raft should now trim its log as much as possible.
func (rf *Raft) Snapshot(index int, snapshot []byte) {
	// Your code here (2D).

}

type InstallSnapShotArgs struct {
	Term				int // leader's term
	LeaderId			int	// for followers to redirect clients
	LastIncludedIndex	int	//the snapshot replaces all entries up through and including this index
	LastIncludedTerm	int	// term of lastIncludedIndex
	Data 				[]byte	//raw bytes of the snapshot chunk, starting at offset
}

type InstallSnapShotReply struct {
	Term	int //cur Term , for leader to update itself
}

func(rf *Raft) SendInstallSnapShot(server int, args * InstallSnapShotArgs, reply *InstallSnapShotReply) bool {
	ok := rf.peers[server].Call("Raft.InstallSnapShot", args, reply)
	return ok
}

func(rf *Raft)InstallSnapShot(args *InstallSnapShotArgs, reply *InstallSnapShotReply) {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	reply.Term = rf.currentTerm
	if args.Term < rf.currentTerm {
		rf.becomeFollower(args.Term)
	}
	send(rf.appendLogCh)
	if args.LastIncludedIndex <= rf.lastIncludedIndex {
		return
	}
	applyMsg := ApplyMsg{CommandValid: false, Snapshot: args.Data}

	if args.LastIncludedIndex < rf.logLen() - 1 {
		rf.log = append(make([]Log, 0),rf.log[args.LastIncludedIndex - rf.lastIncludedIndex:]... )
	}else {
		rf.log = []Log{{args.LastIncludedTerm, nil}}
	}

	rf.lastIncludedIndex, rf.lastIncludedTerm = args.LastIncludedIndex, args.LastIncludedTerm
	rf.persistWithSnapshot(args.Data)
	rf.commitIndex = Max(rf.commitIndex, rf.lastIncludedIndex)
	rf.lastApplied = Max(rf.lastApplied, rf.lastIncludedIndex)
	if rf.lastApplied > rf.lastIncludedIndex{return}
	rf.applyCh <- applyMsg
}

func(rf *Raft)sendSnapshot(server int) {
	args := InstallSnapShotArgs{
		rf.currentTerm,
		rf.lastIncludedIndex,
		rf.lastIncludedTerm,
		rf.me,
		rf.persister.ReadSnapshot(),
	}
	rf.mu.Unlock()
	reply := &InstallSnapShotReply{}
	ret := rf.SendInstallSnapShot(server, &args, reply)
	rf.mu.Lock()
	defer rf.mu.Unlock()
	if !ret || rf.state != Leader || rf.currentTerm != args.Term {
		return
	}
	if reply.Term > rf.currentTerm {
		rf.becomeFollower(reply.Term)
		return
	}
	rf.updateNextMatchIndex(server, rf.lastIncludedIndex)
}

// RequestVoteArgs
// example RequestVote RPC arguments structure.
// field names must start with capital letters!
//
type RequestVoteArgs struct {
	Term			int
	CandidateId		int
	LastLogIndex	int	//"index of candidate's last log entry"
	LastLogTerm		int	//"term of candidate's last log entry"
}

// RequestVoteReply
// example RequestVote RPC reply structure.
// field names must start with capital letters!
//
type RequestVoteReply struct {
	Term		int		//"currentTerm, for candidate to update itself"
	VoteGranted bool	//"true means  raft instance Granted vote for candidate with cur Term"
}

// RequestVote
// example RequestVote RPC handler.
//
func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	if args.Term > rf.currentTerm { // all server rule 1 if RPC request or response contains Term T > currentTerm
		rf.becomeFollower(args.Term) //set currentTerm = T. convert to follower
	}
	reply.Term =rf.currentTerm
	reply.VoteGranted = false
	if args.Term < rf.currentTerm || rf.voteFor != None && rf.voteFor != args.CandidateId {

	}else if args.LastLogTerm < rf.getLastLogTerm() || args.LastLogTerm == rf.getLastLogTerm() && args.LastLogIndex < rf.getLastLogIndex() {

	}else{
		rf.voteFor = args.CandidateId
		reply.VoteGranted = true
		rf.state = Follower
		rf.persist()
		send(rf.voteCh)
	}
}

//
// example code to send a RequestVote RPC to a server.
// server is the index of the target server in rf.peers[].
// expects RPC arguments in args.
// fills in *reply with RPC reply, so caller should
// pass &reply.
// the types of the args and reply passed to Call() must be
// the same as the types of the arguments declared in the
// handler function (including whether they are pointers).
//
// The labrpc package simulates a lossy network, in which servers
// may be unreachable, and in which requests and replies may be lost.
// Call() sends a request and waits for a reply. If a reply arrives
// within a timeout interval, Call() returns true; otherwise
// Call() returns false. Thus Call() may not return for a while.
// A false return can be caused by a dead server, a live server that
// can't be reached, a lost request, or a lost reply.
//
// Call() is guaranteed to return (perhaps after a delay) *except* if the
// handler function on the server side does not return.  Thus there
// is no need to implement your own timeouts around Call().
//
// look at the comments in ../labrpc/labrpc.go for more details.
//
// if you're having trouble getting RPC to work, check that you've
// capitalized all field names in structs passed over RPC, and
// that the caller passes the address of the reply struct with &, not
// the struct itself.
//
func (rf *Raft) sendRequestVote(server int, args *RequestVoteArgs, reply *RequestVoteReply) bool {
	ok := rf.peers[server].Call("Raft.RequestVote", args, reply)
	return ok
}

type AppendEntriesArgs struct {
	Term 			int
	LeaderId		int
	PrevLogIndex	int
	PrevLogTerm		int
	Entries 		[]Log
	LeaderCommit	int
}

type AppendEntriesReply struct {
	Term			int
	Success			bool
	ConflictIndex	int
	ConflictTerm	int
}

func(rf *Raft) AppendEntries(args *AppendEntriesArgs, reply *AppendEntriesReply) {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	defer send(rf.appendLogCh)	// if election timeout elapses without receiving AppendEntries RPC from current Leader
	if args.Term > rf.currentTerm { // if RPC request or response contain Term > rf.CurrentTerm , convert to follower
		rf.becomeFollower(args.Term)
	}
	reply.Term = rf.currentTerm
	reply.Success = false
	reply.ConflictTerm = None
	reply.ConflictIndex = 0
	// return false if log doesn't contain an entry at prevLogIndex whose term matches prevLogTerm
	prevLogIndexTerm := -1
	logSize := rf.logLen()
	if args.PrevLogIndex >= rf.lastIncludedIndex && args.PrevLogIndex < rf.logLen() {
		prevLogIndexTerm = rf.GetLog(args.PrevLogIndex).Term
		//print( rf.GetLog(args.PrevLogIndex).Term,rf.log[args.PrevLogIndex].Term)
	}
	if prevLogIndexTerm != args.PrevLogTerm {
		reply.ConflictIndex = logSize
		if prevLogIndexTerm == -1 {
		}else {
			reply.ConflictTerm = prevLogIndexTerm
			i := rf.lastIncludedIndex
			for ; i < logSize; i++ {
				if rf.GetLog(i).Term == reply.ConflictTerm {
					reply.ConflictIndex = i
					break
				}
			}
		}
		return
	}
	//reply false if term < currentTerm
	if args.Term < rf.currentTerm{return}
	index := args.PrevLogIndex
	for i := 0; i < len(args.Entries); i++ {
		index++
		if index < logSize {
			if rf.GetLog(index).Term == args.Entries[i].Term {
				continue
			}else {
				rf.log = rf.log[:index - rf.lastIncludedIndex]
			}
		}
		rf.log = append(rf.log, args.Entries[i:]...)
		rf.persist()
		break
	}
	if args.LeaderCommit > rf.commitIndex {
		rf.commitIndex = Min(args.LeaderCommit, rf.getLastLogIndex())
		rf.updateLastApplied()
	}
	reply.Success = true
	return
}

func (rf *Raft)SendAppendEntries(server int, args *AppendEntriesArgs, reply *AppendEntriesReply) bool {
	ok := rf.peers[server].Call("Raft.AppendEntries", args, reply)
	return ok
}

func(rf *Raft)startAppendLog() {
	for i := 0; i < len(rf.peers); i++ {
		if i == rf.me{
			continue
		}
		go func(idx int) {
			for {
				rf.mu.Lock()
				if rf.state != Leader {
					rf.mu.Unlock()
					return
				}
				if rf.nextIndex[idx] - rf.lastIncludedIndex < 1 {
					rf.sendSnapshot(idx)
					return
				}
				args := AppendEntriesArgs{
					rf.currentTerm,
					rf.me,
					rf.getPrevLogIndex(idx),
					rf.getPrevLogTerm(idx),
					append(make([]Log, 0), rf.log[rf.nextIndex[idx] - rf.lastIncludedIndex:]...),
					rf.commitIndex,
				}
				rf.mu.Unlock()
				reply := &AppendEntriesReply{}
				ret := rf.SendAppendEntries(idx, &args, reply)
				rf.mu.Lock()
				if !ret || rf.state != Leader || rf.currentTerm != args.Term {
					rf.mu.Unlock()
					return
				}
				if reply.Term > rf.currentTerm {
					rf.becomeFollower(reply.Term)
					rf.mu.Unlock()
					return
				}
				if reply.Success {
					rf.updateNextMatchIndex(idx, args.PrevLogIndex + len(args.Entries))
					rf.mu.Unlock()
					return
				} else {
					tarIndex := reply.ConflictIndex
					//println(reply.ConflictTerm)
					if reply.ConflictTerm != None {
						logSize := rf.logLen()
						for i := rf.lastIncludedIndex; i < logSize; i++ {
							if rf.GetLog(i).Term != reply.ConflictTerm {
								continue
							}
							for i < logSize && rf.GetLog(i).Term == reply.ConflictTerm {
								i++
							}
							tarIndex = i
						}
					}
					rf.nextIndex[idx] = tarIndex
					rf.mu.Unlock()
				}
			}
		}(i)
	}
}


// Start
// the service using Raft (e.g. a k/v server) wants to start
// agreement on the next command to be appended to Raft's log. if this
// server isn't the leader, returns false. otherwise start the
// agreement and return immediately. there is no guarantee that this
// command will ever be committed to the Raft log, since the leader
// may fail or lose an election. even if the Raft instance has been killed,
// this function should return gracefully.
//
// the first return value is the index that the command will appear at
// if it's ever committed. the second return value is the current
// term. the third return value is true if this server believes it is
// the leader.
//
func (rf *Raft) Start(command interface{}) (int, int, bool) {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	index := -1
	term := rf.currentTerm
	isLeader := rf.state == Leader
	if isLeader {
		index = rf.getLastLogIndex() + 1
		NewLog := Log{
			rf.currentTerm,
			command,
		}
		rf.log = append(rf.log, NewLog)
		rf.persist()
		rf.startAppendLog()
	}
	return index, term, isLeader
}

// Kill
// the tester doesn't halt goroutines created by Raft after each test,
// but it does call the Kill() method. your code can use killed() to
// check whether Kill() has been called. the use of atomic avoids the
// need for a lock.
//
// the issue is that long-running goroutines use memory and may chew
// up CPU time, perhaps causing later tests to fail and generating
// confusing debug output. any goroutine with a long-running loop
// should call killed() to check whether it should stop.
//
func (rf *Raft) Kill() {
	atomic.StoreInt32(&rf.dead, 1)
	// Your code here, if desired.
	send(rf.killCh)
}

func (rf *Raft) killed() bool {
	z := atomic.LoadInt32(&rf.dead)
	return z == 1
}

// The ticker go routine starts a new election if this peer hasn't received
// heartsbeats recently.
func (rf *Raft) ticker() {
	for rf.killed() == false {

		// Your code here to check if a leader election should
		// be started and to randomize sleeping time using
		// time.Sleep().

	}
}

func(rf *Raft)becomeFollower(term int) {
	rf.currentTerm = term
	rf.state = Follower
	rf.voteFor = None
	rf.persist()
}

func(rf *Raft)becomeLeader() {
	if rf.state != Candidate {
		return
	}
	rf.state = Leader
	rf.nextIndex = make([]int, len(rf.peers))
	rf.matchIndex = make([]int, len(rf.peers))
	for i := 0; i < len(rf.nextIndex); i++ {
		rf.nextIndex[i] = rf.getLastLogIndex() + 1
	}
}

func(rf *Raft)becomeCandidate() {
	rf.state = Candidate
	rf.currentTerm++
	rf.voteFor = rf.me
	rf.persist()
	go rf.startElection()
}

func(rf *Raft)startElection() {
	rf.mu.Lock()
	args := RequestVoteArgs{
		rf.currentTerm,
		rf.me,
		rf.getLastLogIndex(),
		rf.getLastLogTerm(),
	}
	rf.mu.Unlock()
	var vote int64 = 1
	for i := 0; i < len(rf.peers); i++ {
		if i == rf.me {
			continue
		}
		go func(idx int) {
				reply := &RequestVoteReply{}
				ret := rf.sendRequestVote(idx, & args, reply)
				if ret {
					rf.mu.Lock()
					defer rf.mu.Unlock()
					if reply.Term > rf.currentTerm {
						rf.becomeFollower(reply.Term)
						return
					}
					if rf.state != Candidate || rf.currentTerm != args.Term {
						return
					}
					if reply.VoteGranted{
						atomic.AddInt64(&vote,1)
					}
					if atomic.LoadInt64(&vote) > int64(len(rf.peers) / 2) {
						rf.becomeLeader()
						rf.startAppendLog()
						send(rf.voteCh) // after be leader, "select" goroutines will send heartbeats immediately
					}
				}
			return
		}(i)
	}
}

func send(ch chan bool) {
	select {
	case <-ch:
	default:
	}
	ch <- true
}

func (rf *Raft)getLastLogTerm() int {
	idx := rf.getLastLogIndex()
	if idx < rf.lastIncludedIndex {
		return -1
	}
	return rf.log[idx].Term
}

func (rf *Raft)getLastLogIndex() int {
	return rf.logLen() - 1
}

func (rf *Raft)logLen() int {
	return len(rf.log) + rf.lastIncludedIndex
}

func(rf *Raft)GetLog(i int) Log {
	return rf.log[i - rf.lastIncludedIndex]
}

func(rf *Raft)getPrevLogIndex(i int) int {
	return rf.nextIndex[i] - 1
}

func(rf *Raft)getPrevLogTerm(i int) int {
	prevLogIdx := rf.getPrevLogIndex(i)
	if prevLogIdx < rf.lastIncludedIndex {
		return -1
	}
	return rf.GetLog(prevLogIdx).Term
}

func(rf *Raft)updateNextMatchIndex(server int, matchIdx int) {
	rf.matchIndex[server] = matchIdx
	rf.nextIndex[server] = matchIdx + 1
	rf.updateCommitIndex()
}


func(rf *Raft) updateCommitIndex() {
	rf.matchIndex[rf.me] = rf.logLen() - 1
	copyMatchIdx := make([]int, len(rf.matchIndex))
	copy(copyMatchIdx, rf.matchIndex)
	sort.Slice(copyMatchIdx, func(i, j int) bool {
		return copyMatchIdx[i] < copyMatchIdx[j]
	})
	N := copyMatchIdx[len(copyMatchIdx)/2]
	//fmt.Println(N)
	if N > rf.commitIndex && rf.GetLog(N).Term == rf.currentTerm {
		rf.commitIndex = N
		//fmt.Printf("   #%v#", rf.commitIndex)
		rf.updateLastApplied()
	}
}

func(rf *Raft)updateLastApplied() {
	rf.lastApplied = Max(rf.lastApplied, rf.lastIncludedIndex)
	rf.commitIndex = Max(rf.commitIndex, rf.lastIncludedIndex)
	for rf.lastApplied < rf.commitIndex {
		rf.lastApplied++
		curLog := rf.GetLog(rf.lastApplied)
		applyMsg := ApplyMsg{CommandValid: true, Command: curLog.Command, CommandIndex: rf.lastApplied}
		rf.applyCh <- applyMsg
	}
}
// Make
// the service or tester wants to create a Raft server. the ports
// of all the Raft servers (including this one) are in peers[]. this
// server's port is peers[me]. all the servers' peers[] arrays
// have the same order. persister is a place for this server to
// save its persistent state, and also initially holds the most
// recent saved state, if any. applyCh is a channel on which the
// tester or service expects Raft to send ApplyMsg messages.
// Make() must return quickly, so it should start goroutines
// for any long-running work.
//
func Make(peers []*labrpc.ClientEnd, me int,
	persister *Persister, applyCh chan ApplyMsg) *Raft {
	rf := &Raft{}
	rf.peers = peers
	rf.persister = persister
	rf.me = me

	rf.state = Follower
	rf.currentTerm = 0
	rf.voteFor = None
	rf.log= make([]Log, 1)
	rf.commitIndex = 0
	rf.lastApplied = 0

	rf.voteCh = make(chan bool, 1)
	rf.applyCh = applyCh
	rf.killCh = make(chan bool, 1)
	rf.appendLogCh = make(chan bool, 1)

	rf.readPersist(persister.ReadRaftState())
	heartBeatTime := time.Duration(100) * time.Millisecond

	go func() {
		for {
			select {
			case <-rf.killCh:
				return
			default:
			}
			electionTime := time.Duration(rand.Intn(200)+300)*time.Millisecond
			state := rf.state
			switch state {
			case Follower, Candidate:
				select {
				case <-rf.voteCh:
				case <-rf.appendLogCh:
				case <-time.After(electionTime):
					rf.mu.Lock()
					rf.becomeCandidate()
					rf.mu.Unlock()
				}
			case Leader:
				time.Sleep(heartBeatTime)
				rf.startAppendLog()
			}
		}
	}()
	return rf
}