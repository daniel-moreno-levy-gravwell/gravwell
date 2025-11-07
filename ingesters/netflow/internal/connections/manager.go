/*************************************************************************
 * Copyright 2025 Gravwell, Inc. All rights reserved.
 * Contact: <legal@gravwell.io>
 *
 * This software may be modified and distributed under the terms of the
 * BSD 2-clause license. See the LICENSE file for details.
 **************************************************************************/

// Package connections manages active connections to the netflow ingester
package connections

import (
	"sync"
)

type closer interface {
	Close() error
}

type Manager struct {
	sync.Mutex
	nextConnID int
	connClosers map[int]closer
}

func NewManager() *Manager {
	return &Manager{
		connClosers: make(map[int]closer),
	}
}

func (m *Manager) Add(c closer) int {
	m.Lock()
	m.nextConnID++
	id := m.nextConnID
	m.connClosers[id] = c
	m.Unlock()
	return id
}

func (m *Manager) Del(id int) {
	m.Lock()
	delete(m.connClosers, id)
	m.Unlock()
}

func (m *Manager) Count() int {
	m.Lock()
	defer m.Unlock()
	return len(m.connClosers)
}

func (m *Manager) CloseAll() {
	m.Lock()
	for _, v := range m.connClosers {
		v.Close()
	}
	defer m.Unlock()
}
