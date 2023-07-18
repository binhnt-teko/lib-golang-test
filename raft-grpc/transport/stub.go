/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package transport

import (
	"fmt"
	"sync"

	"github.com/pkg/errors"
)

// Stub holds all information about the remote node,
// including the RemoteContext for it, and serializes
// some operations on it.
type Stub struct {
	lock sync.RWMutex
	// ID is unique among all members, and cannot be 0.
	ID uint64
	// Endpoint is the endpoint of the node, denoted in %s:%d format
	Endpoint string
	// Identity []byte

	*RemoteContext
}

// Active returns whether the Stub
// is active or not
func (stub *Stub) Active() bool {
	stub.lock.RLock()
	defer stub.lock.RUnlock()
	return stub.isActive()
}

// Active returns whether the Stub
// is active or not.
func (stub *Stub) isActive() bool {
	return stub.RemoteContext != nil
}

// Deactivate deactivates the Stub and
// ceases all communication operations
// invoked on it.
func (stub *Stub) Deactivate() {
	stub.lock.Lock()
	defer stub.lock.Unlock()
	if !stub.isActive() {
		return
	}
	stub.RemoteContext = nil
}

// Activate creates a remote context with the given function callback
// in an atomic manner - if two parallel invocations are invoked on this Stub,
// only a single invocation of createRemoteStub takes place.
func (stub *Stub) Activate(createRemoteContext func() (*RemoteContext, error)) error {
	stub.lock.Lock()
	defer stub.lock.Unlock()
	// Check if the stub has already been activated while we were waiting for the lock
	if stub.isActive() {
		return nil
	}
	fmt.Printf("Stub.Activate: call createRemoteContext \n")
	remoteStub, err := createRemoteContext()
	if err != nil {
		return errors.WithStack(err)
	}

	stub.RemoteContext = remoteStub
	return nil
}
