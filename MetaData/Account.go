package MetaData

import (
	"fmt"
	"strconv"
)

//go:generate msgp
type Account struct {
	WorkerNumberSet     map[string]uint32 `msg:"WorkerNumberSet"`
	VoterNumberSet      map[string]uint32 `msg:"VoterNumberSet"`
	VoterSet            map[string]string `msg:"VoterSet"`
	WorkerSet           map[string]string `msg:"WorkerSet"`
	WorkerCandidateSet  map[string]string `msg:"WorkerCandidateSet"`
	WorkerCandidateList []string          `msg:"WorkerCandidateList"`
}

func (min Account) ToByteArray() []byte {
	data, _ := min.MarshalMsg(nil)
	return data
}

func (min *Account) FromByteArray(data []byte) {
	_, err := min.UnmarshalMsg(data)
	if err != nil {
		fmt.Println("err=", err)
	}
}

func (a *Account) SetVoterSet(set map[string]uint64) {
	voterSet := make(map[string]string)
	for k, v := range set {
		voterSet[k] = strconv.FormatUint(v, 10)
	}
	a.VoterSet = voterSet
}

func (a *Account) GetVoterSet() map[string]uint64 {
	voterSet := make(map[string]uint64)
	for k, v := range a.VoterSet {
		result, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			fmt.Println(err)
		}
		voterSet[k] = result
	}
	return voterSet
}

func (a *Account) SetWorkerCandidateSet(set map[string]uint64) {
	workerCandidateSet := make(map[string]string)
	for k, v := range set {
		workerCandidateSet[k] = strconv.FormatUint(v, 10)
	}
	a.WorkerCandidateSet = workerCandidateSet
}

func (a *Account) GetWorkerCandidateSet() map[string]uint64 {
	workerCandidateSet := make(map[string]uint64)
	for k, v := range a.WorkerCandidateSet {
		result, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			fmt.Println(err)
		}
		workerCandidateSet[k] = result
	}
	return workerCandidateSet
}

func (a *Account) SetWorkerSet(set map[string]uint64) {
	workerSet := make(map[string]string)
	for k, v := range set {
		workerSet[k] = strconv.FormatUint(v, 10)
	}
	a.WorkerSet = workerSet
}

func (a *Account) GetWorkerSet() map[string]uint64 {
	workerSet := make(map[string]uint64)
	for k, v := range a.WorkerSet {
		result, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			fmt.Println(err)
		}
		workerSet[k] = result
	}
	return workerSet
}

func (a *Account) SetWorkerNumberSet(set map[uint32]string) {
	workerNumberSet := make(map[string]uint32)

	for k, v := range set {
		workerNumberSet[v] = k
	}

	a.WorkerNumberSet = workerNumberSet
}

func (a *Account) GetWorkerNumberSet() map[uint32]string {
	workerNumberSet := make(map[uint32]string)

	for k, v := range a.WorkerNumberSet {
		workerNumberSet[v] = k
	}
	return workerNumberSet
}

func (a *Account) SetVoterNumberSet(set map[uint32]string) {
	voterNumberSet := make(map[string]uint32)

	for k, v := range set {
		voterNumberSet[v] = k
	}

	a.VoterNumberSet = voterNumberSet
}

func (a *Account) GetVorterNumberSet() map[uint32]string {
	vorterNumberSet := make(map[uint32]string)

	for k, v := range a.VoterNumberSet {
		vorterNumberSet[v] = k
	}
	return vorterNumberSet
}
