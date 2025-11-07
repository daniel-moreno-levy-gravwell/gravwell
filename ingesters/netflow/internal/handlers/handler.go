/*************************************************************************
 * Copyright 2025 Gravwell, Inc. All rights reserved.
 * Contact: <legal@gravwell.io>
 *
 * This software may be modified and distributed under the terms of the
 * BSD 2-clause license. See the LICENSE file for details.
 **************************************************************************/

// Package handlers defines the main logic for datagram ingestion
package handlers

import (
	"errors"
	"sync"
	"time"

	"github.com/gravwell/gravwell/v3/ingest"
	"github.com/gravwell/gravwell/v3/ingest/entry"
	"github.com/gravwell/gravwell/v3/ingest/log"
	"github.com/gravwell/gravwell/v3/ingesters/netflow/internal/connections"
)

var (
	ErrAlreadyListening = errors.New("already listening")
	ErrAlreadyClosed    = errors.New("already closed")
	ErrNotReady         = errors.New("not Ready")
)

type BindHandler interface {
	String() string
	Listen(string) error
	Close() error
	Start(int) error
}

type BindConfig struct {
	Tag                entry.EntryTag
	Ch                 chan *entry.Entry
	Wg                 *sync.WaitGroup
	IgnoreTS           bool
	LocalTZ            bool
	Igst               *ingest.IngestMuxer
	LastInfoDump       time.Time
	SessionDumpEnabled bool
	ConnManager        *connections.Manager
	Log                *log.Logger
}

func (bc BindConfig) Validate() error {
	if bc.Ch == nil {
		return errors.New("nil channel")
	}
	if bc.Wg == nil {
		return errors.New("nil wait group")
	}
	if bc.Igst == nil {
		return errors.New("nil ingest muxer")
	}
	return nil
}
