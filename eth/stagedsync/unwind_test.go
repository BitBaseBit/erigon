package stagedsync

import (
	"testing"

	"github.com/ledgerwatch/erigon/common"
	"github.com/ledgerwatch/erigon/eth/stagedsync/stages"
	"github.com/ledgerwatch/erigon/ethdb"
	"github.com/stretchr/testify/assert"
)

func TestUnwindStackLoadFromDb(t *testing.T) {
	_, tx := ethdb.NewTestTx(t)

	stack := NewPersistentUnwindStack()

	stages := []stages.SyncStage{stages.Bodies, stages.Headers, stages.Execution}

	points := []uint64{10, 20, 30}

	for i := range stages {
		err := stack.Add(UnwindState{stages[i], points[i], common.Hash{}}, tx)
		assert.NoError(t, err)
	}

	stack2 := NewPersistentUnwindStack()
	for i := range stages {
		err := stack2.AddFromDB(tx, stages[i])
		assert.NoError(t, err)
	}

	assert.Equal(t, stack.unwindStack, stack2.unwindStack)
	assert.Equal(t, len(stages), len(stack2.unwindStack))
}

func TestUnwindStackLoadFromDbAfterDone(t *testing.T) {
	_, tx := ethdb.NewTestTx(t)

	stack := NewPersistentUnwindStack()

	stages := []stages.SyncStage{stages.Bodies, stages.Headers, stages.Execution}

	points := []uint64{10, 20, 30}

	for i := range stages {
		err := stack.Add(UnwindState{stages[i], points[i], common.Hash{}}, tx)
		assert.NoError(t, err)
	}

	u := stack.Pop()
	assert.NotNil(t, u)
	err := u.Done(tx)
	assert.NoError(t, err)

	stack2 := NewPersistentUnwindStack()
	for i := range stages {
		err := stack2.AddFromDB(tx, stages[i])
		assert.NoError(t, err)
	}

	assert.Equal(t, stack.unwindStack, stack2.unwindStack)
	assert.Equal(t, len(stages)-1, len(stack2.unwindStack))
}

func TestUnwindStackLoadFromDbNoDone(t *testing.T) {
	_, tx := ethdb.NewTestTx(t)

	stack := NewPersistentUnwindStack()

	stages := []stages.SyncStage{stages.Bodies, stages.Headers, stages.Execution}

	points := []uint64{10, 20, 30}

	for i := range stages {
		err := stack.Add(UnwindState{stages[i], points[i], common.Hash{}}, tx)
		assert.NoError(t, err)
	}

	u := stack.Pop()
	assert.NotNil(t, u)

	stack2 := NewPersistentUnwindStack()
	for i := range stages {
		err := stack2.AddFromDB(tx, stages[i])
		assert.NoError(t, err)
	}

	assert.NotEqual(t, stack.unwindStack, stack2.unwindStack)
	assert.Equal(t, len(stages), len(stack2.unwindStack))
}

func TestUnwindStackPopAndEmpty(t *testing.T) {
	_, tx := ethdb.NewTestTx(t)

	stack := NewPersistentUnwindStack()

	stages := []stages.SyncStage{stages.Bodies, stages.Headers, stages.Execution}

	points := []uint64{10, 20, 30}

	for i := range stages {
		err := stack.Add(UnwindState{stages[i], points[i], common.Hash{}}, tx)
		assert.NoError(t, err)
	}

	assert.False(t, stack.Empty())
	u := stack.Pop()
	assert.NotNil(t, u)

	assert.False(t, stack.Empty())
	u = stack.Pop()
	assert.NotNil(t, u)

	assert.False(t, stack.Empty())
	u = stack.Pop()
	assert.NotNil(t, u)

	assert.True(t, stack.Empty())
	u = stack.Pop()
	assert.Nil(t, u)

	assert.True(t, stack.Empty())
}

func TestUnwindOverrideWithLower(t *testing.T) {
	db := ethdb.NewTestDB(t)

	stack := NewPersistentUnwindStack()

	stages := []stages.SyncStage{stages.Bodies, stages.Headers, stages.Execution}

	points := []uint64{10, 20, 30}

	for i := range stages {
		err := stack.Add(UnwindState{stages[i], points[i], common.Hash{}}, db)
		assert.NoError(t, err)
	}

	assert.Equal(t, 3, len(stack.unwindStack))

	err := stack.Add(UnwindState{stages[0], 5, common.Hash{}}, db)
	assert.NoError(t, err)

	// we append if the next unwind is to the lower block
	assert.Equal(t, 4, len(stack.unwindStack))
}

func TestUnwindOverrideWithHigher(t *testing.T) {
	_, tx := ethdb.NewTestTx(t)

	stack := NewPersistentUnwindStack()

	stages := []stages.SyncStage{stages.Bodies, stages.Headers, stages.Execution}

	points := []uint64{10, 20, 30}

	for i := range stages {
		err := stack.Add(UnwindState{stages[i], points[i], common.Hash{}}, tx)
		assert.NoError(t, err)
	}

	assert.Equal(t, 3, len(stack.unwindStack))

	err := stack.Add(UnwindState{stages[0], 105, common.Hash{}}, tx)
	assert.NoError(t, err)

	// we ignore if next unwind is to the higher block
	assert.Equal(t, 3, len(stack.unwindStack))
}

func TestUnwindOverrideWithTheSame(t *testing.T) {
	_, tx := ethdb.NewTestTx(t)

	stack := NewPersistentUnwindStack()

	stages := []stages.SyncStage{stages.Bodies, stages.Headers, stages.Execution}

	points := []uint64{10, 20, 30}

	for i := range stages {
		err := stack.Add(UnwindState{stages[i], points[i], common.Hash{}}, tx)
		assert.NoError(t, err)
	}

	assert.Equal(t, 3, len(stack.unwindStack))

	err := stack.Add(UnwindState{stages[0], 10, common.Hash{}}, tx)
	assert.NoError(t, err)

	// we ignore if next unwind is to the higher block
	assert.Equal(t, 3, len(stack.unwindStack))
}
