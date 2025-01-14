//-----------------------------------------------------------------------------
// Copyright (c) 2020 Detlef Stern
//
// This file is part of zettelstore.
//
// Zettelstore is free software: you can redistribute it and/or modify it under
// the terms of the GNU Affero General Public License as published by the Free
// Software Foundation, either version 3 of the License, or (at your option)
// any later version.
//
// Zettelstore is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License
// for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with Zettelstore. If not, see <http://www.gnu.org/licenses/>.
//-----------------------------------------------------------------------------

// Package directory manages the directory part of a dirstore.
package directory

import (
	"sync"
	"time"

	"zettelstore.de/z/domain"
	"zettelstore.de/z/place"
)

// Service specifies a directory scan service.
type Service struct {
	dirPath     string
	ticker      *time.Ticker
	cmds        chan dirCmd
	changeFuncs []place.ObserverFunc
	mxFuncs     sync.RWMutex
}

// NewService creates a new directory service.
func NewService(directoryPath string, rescanTime time.Duration) *Service {
	srv := &Service{
		dirPath: directoryPath,
		ticker:  time.NewTicker(rescanTime),
		cmds:    make(chan dirCmd),
	}
	return srv
}

// Start makes the directory service operational.
func (srv *Service) Start() {
	done := make(chan struct{})
	rawEvents := make(chan *fileEvent)
	events := make(chan *fileEvent)

	ready := make(chan int)
	go srv.directoryService(events, ready)
	go collectEvents(events, rawEvents)
	go watchDirectory(srv.dirPath, rawEvents, done)
	go ping(done, srv.ticker)
	<-ready
}

// Stop stops the directory service.
func (srv *Service) Stop() {
	srv.ticker.Stop()
	srv.ticker = nil
}

// Subscribe to invalidation events.
func (srv *Service) Subscribe(changeFunc place.ObserverFunc) {
	srv.mxFuncs.Lock()
	if changeFunc != nil {
		srv.changeFuncs = append(srv.changeFuncs, changeFunc)
	}
	srv.mxFuncs.Unlock()
}

func (srv *Service) notifyChange(all bool, zid domain.ZettelID) {
	srv.mxFuncs.RLock()
	changeFuncs := srv.changeFuncs
	srv.mxFuncs.RUnlock()
	for _, changeF := range changeFuncs {
		changeF(all, zid)
	}
}

// MetaSpec defines all possibilities where meta data can be stored.
type MetaSpec int

// Constants for MetaSpec
const (
	MetaSpecUnknown MetaSpec = iota
	MetaSpecNone             // no meta information
	MetaSpecFile             // meta information is in meta file
	MetaSpecHeader           // meta information is in header
)

// Entry stores everything for a directory entry.
type Entry struct {
	Zid         domain.ZettelID
	MetaSpec    MetaSpec // location of meta information
	MetaPath    string   // file path of meta information
	ContentPath string   // file path of zettel content
	ContentExt  string   // (normalized) file extension of zettel content
	Duplicates  bool     // multiple content files
}

// IsValid checks whether the entry is valid.
func (e *Entry) IsValid() bool {
	return e.Zid.IsValid()
}

// GetEntries returns an unsorted list of all current directory entries.
func (srv *Service) GetEntries() []Entry {
	resChan := make(chan resGetEntries)
	srv.cmds <- &cmdGetEntries{resChan}
	return <-resChan
}

// GetEntry returns the entry with the specified zettel id. If there is no such
// zettel id, an empty entry is returned.
func (srv *Service) GetEntry(zid domain.ZettelID) Entry {
	resChan := make(chan resGetEntry)
	srv.cmds <- &cmdGetEntry{zid, resChan}
	return <-resChan
}

// GetNew returns an entry with a new zettel id.
func (srv *Service) GetNew() Entry {
	resChan := make(chan resNewEntry)
	srv.cmds <- &cmdNewEntry{resChan}
	return <-resChan
}

// UpdateEntry notifies the directory of an updated entry.
func (srv *Service) UpdateEntry(entry *Entry) {
	resChan := make(chan struct{})
	srv.cmds <- &cmdUpdateEntry{entry, resChan}
	<-resChan
}

// RenameEntry notifies the directory of an renamed entry.
func (srv *Service) RenameEntry(curEntry, newEntry *Entry) error {
	resChan := make(chan resRenameEntry)
	srv.cmds <- &cmdRenameEntry{curEntry, newEntry, resChan}
	return <-resChan
}

// DeleteEntry removes a zettel id from the directory of entries.
func (srv *Service) DeleteEntry(zid domain.ZettelID) {
	resChan := make(chan struct{})
	srv.cmds <- &cmdDeleteEntry{zid, resChan}
	<-resChan
}
