/*************************************************************************
 * Copyright 2025 Gravwell, Inc. All rights reserved.
 * Contact: <legal@gravwell.io>
 *
 * This software may be modified and distributed under the terms of the
 * BSD 2-clause license. See the LICENSE file for details.
 **************************************************************************/

package main

import (
	"errors"
	"fmt"
	"net"
	"sync"
)

type SFlowV5Handler struct {
	bindConfig
	mtx   *sync.Mutex
	c     *net.UDPConn
	ready bool
}

func NewSFlowV5Handler(c bindConfig) (*SFlowV5Handler, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	return &SFlowV5Handler{
		bindConfig: c,
		mtx:        &sync.Mutex{},
	}, nil
}

func (s *SFlowV5Handler) String() string {
	return `sFlowV5`
}

func (s *SFlowV5Handler) Listen(addr string) (err error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if s.c != nil {
		err = ErrAlreadyListening
		return
	}
	var a *net.UDPAddr
	if a, err = net.ResolveUDPAddr("udp", addr); err != nil {
		return
	}
	if s.c, err = net.ListenUDP("udp", a); err == nil {
		s.ready = true
	}
	return
}

func (s *SFlowV5Handler) Close() error {
	if s == nil {
		return ErrAlreadyClosed
	}
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.ready = false
	return s.c.Close()
}

func (s *SFlowV5Handler) Start(id int) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if !s.ready || s.c == nil {
		fmt.Println(s.ready, s.c)
		return ErrNotReady
	}
	if id < 0 {
		return errors.New("invalid id")
	}
	go s.routine(id)
	return nil
}

func (s *SFlowV5Handler) routine(id int) {
	// TODO
}
