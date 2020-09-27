package domain

import (
	"StaticRaceDetector/utils"
	"encoding/json"
	"go/token"
	"golang.org/x/tools/go/ssa"
)

type OpKind int

const (
	GuardAccessRead OpKind = iota
	GuardAccessWrite
)

type GuardedAccess struct {
	ID      int
	Pos     token.Pos
	Value   ssa.Value
	State   *GoroutineState
	Lockset *Lockset
	OpKind  OpKind
}

type GuardedAccessJSON struct {
	ID          int
	Pos         token.Pos
	Value       int
	OpKind      OpKind
	LocksetJson *LocksetJson
	State       *GoroutineStateJSON
}

func (ga *GuardedAccess) ToJSON() GuardedAccessJSON {
	dumpJson := GuardedAccessJSON{}
	dumpJson.ID = ga.ID
	dumpJson.Pos = ga.Pos
	dumpJson.Value = int(ga.Value.Pos())
	dumpJson.OpKind = ga.OpKind
	dumpJson.LocksetJson = ga.Lockset.ToJSON()
	dumpJson.State = ga.State.ToJSON()
	return dumpJson
}
func (ga *GuardedAccess) MarshalJSON() ([]byte, error) {
	dump, err := json.Marshal(ga.ToJSON())
	return dump, err
}

func (ga *GuardedAccess) Intersects(gaToCompare *GuardedAccess) bool {
	if ga.ID == gaToCompare.ID || ga.State.GoroutineID == gaToCompare.State.GoroutineID {
		return true
	}
	if ga.OpKind == GuardAccessRead && gaToCompare.OpKind == GuardAccessRead {
		return true
	}
	for _, lockA := range ga.Lockset.ExistingLocks {
		for _, lockB := range gaToCompare.Lockset.ExistingLocks {
			if lockA.Pos() == lockB.Pos() {
				return true
			}
		}
	}
	return false
}

var GuardedAccessCounter = utils.NewCounter()

func AddGuardedAccess(pos token.Pos, value ssa.Value, kind OpKind, lockset *Lockset, goroutineState *GoroutineState) *GuardedAccess {
	goroutineState.Increment()
	return &GuardedAccess{ID: GuardedAccessCounter.GetNext(), Pos: pos, Value: value, Lockset: lockset.Copy(), OpKind: kind, State: goroutineState.Copy()}
}
