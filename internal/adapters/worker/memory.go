package worker

import "github.com/alitto/pond/v2"

type Memory struct {
	pool pond.Pool
}

func NewMemory() *Memory {
	pool := pond.NewPool(100)
	return &Memory{pool: pool}
}

func (m *Memory) Submit(f func()) {
	m.pool.Submit(f)
}

func (m *Memory) Shutdown() {
	m.pool.StopAndWait()
}
