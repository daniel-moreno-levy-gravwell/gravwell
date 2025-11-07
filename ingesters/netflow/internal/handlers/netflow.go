/*************************************************************************
 * Copyright 2025 Gravwell, Inc. All rights reserved.
 * Contact: <legal@gravwell.io>
 *
 * This software may be modified and distributed under the terms of the
 * BSD 2-clause license. See the LICENSE file for details.
 **************************************************************************/

package handlers

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"

	"github.com/gravwell/gravwell/v3/ingest/entry"
	"github.com/gravwell/gravwell/v3/netflow"
)

type NetflowV5Handler struct {
	BindConfig
	c     *net.UDPConn
	ready bool
}

func NewNetflowV5Handler(c BindConfig) (*NetflowV5Handler, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	return &NetflowV5Handler{
		BindConfig: c,
	}, nil
}

func (n *NetflowV5Handler) String() string {
	return `NetflowV5`
}

func (n *NetflowV5Handler) Listen(s string) (err error) {
	n.ConnManager.Lock()
	defer n.ConnManager.Unlock()
	if n.c != nil {
		err = ErrAlreadyListening
		return
	}
	var a *net.UDPAddr
	if a, err = net.ResolveUDPAddr("udp", s); err != nil {
		return
	}
	if n.c, err = net.ListenUDP("udp", a); err == nil {
		n.ready = true
	}
	return
}

func (n *NetflowV5Handler) Close() error {
	if n == nil {
		return ErrAlreadyClosed
	}
	n.ConnManager.Lock()
	defer n.ConnManager.Unlock()
	n.ready = false
	return n.c.Close()
}

func (n *NetflowV5Handler) Start(id int) error {
	n.ConnManager.Lock()
	defer n.ConnManager.Unlock()
	if !n.ready || n.c == nil {
		fmt.Println(n.ready, n.c)
		return ErrNotReady
	}
	if id < 0 {
		return errors.New("invalid id")
	}
	go n.routine(id)
	return nil
}

func (n *NetflowV5Handler) routine(id int) {
	defer n.Wg.Done()
	defer n.ConnManager.Del(id)
	var nf netflow.NFv5
	var l int
	var addr *net.UDPAddr
	var err error
	var ts entry.Timestamp
	tbuff := make([]byte, netflow.HeaderSize+(30*netflow.RecordSize))
	for {
		if _, addr, err = n.c.ReadFromUDP(tbuff); err != nil {
			return
		}
		if l, err = nf.ValidateSize(tbuff); err != nil {
			continue //there isn't much we can do about bad packets...
		}
		lbuff := make([]byte, l)
		copy(lbuff, tbuff[0:l])
		if n.IgnoreTS {
			ts = entry.Now()
		} else {
			ts = entry.UnixTime(int64(binary.BigEndian.Uint32(lbuff[8:12])), int64(binary.BigEndian.Uint32(lbuff[12:16])))
		}
		e := &entry.Entry{
			Tag:  n.Tag,
			SRC:  addr.IP,
			TS:   ts,
			Data: lbuff,
		}
		n.Ch <- e
	}
}
