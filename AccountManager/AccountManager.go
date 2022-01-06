package AccountManager

type NodeID = uint64

type AccountManager struct {
	WorkerNumberSet    map[uint32]string
	VoterNumberSet     map[uint32]string
	VoterSet           map[string]NodeID
	WorkerSet          map[string]NodeID
	WorkerCandidateSet map[string]NodeID

	WorkerCandidateList []string
}
