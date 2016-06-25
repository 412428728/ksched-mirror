// Copyright 2016 The ksched Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Common type definitions
// C++ file: https://github.com/camsas/firmament/blob/master/src/base/types.h

package types

import (
	"sync"

	rs "github.com/coreos/ksched/pkg/types/resourcestatus"
	pb "github.com/coreos/ksched/proto"
)

type (
	TaskOutputID uint32
	TaskID       uint64
	EquivClass   uint64
	ResourceID   uint64
	JobID        uint64
)

// Thread safe maps: Acquire and release lock on read/write
// When initializing the map type, make sure to make map
// NOTE: These maps only take pointer values, so change wherever a direct struct is passed below
type ResourceMap struct {
	RWMu sync.RWMutex
	M    map[ResourceID]*rs.ResourceStatus
}

type JobMap struct {
	RWMu sync.RWMutex
	M    map[JobID]*pb.JobDescriptor
	// NOTE: changing to be a pointer to JobDescriptor since
	// you cannot have pointers to map values
}

type TaskMap struct {
	RWMu sync.RWMutex
	M    map[TaskID]*pb.TaskDescriptor
}

// Multimap: multiple values can map to a single key
type MultiMap map[uint64]map[uint64]struct{}

// Maps utility functions used by Firmament
// Implemented as generic templates in C++
// See: https://github.com/camsas/firmament/blob/master/src/misc/map-util.h
// TODO: Decide later what to do with this redundant mess

// Perform a lookup in a map.
// If the key is present in the map then the value associated with that
// key is returned, otherwise the value passed as a default is returned.
func (rm *ResourceMap) FindWithDefault(k ResourceID, dV *rs.ResourceStatus) *rs.ResourceStatus {
	rm.RWMu.RLock()
	defer rm.RWMu.RUnlock()
	v, ok := rm.M[k]
	if !ok {
		v = dV
	}
	return v
}

func (jm *JobMap) FindWithDefault(k JobID, dV *pb.JobDescriptor) *pb.JobDescriptor {
	jm.RWMu.RLock()
	defer jm.RWMu.RUnlock()
	v, ok := jm.M[k]
	if !ok {
		v = dV
	}
	return v
}

func (tm *TaskMap) FindWithDefault(k TaskID, dV *pb.TaskDescriptor) *pb.TaskDescriptor {
	tm.RWMu.RLock()
	defer tm.RWMu.RUnlock()
	v, ok := tm.M[k]
	if !ok {
		v = dV
	}
	return v
}

// NOTE: Did not implement FindOrNull since we cannot have pointers to map values in Go
// Perform a lookup in a map.
// Same as above but the returned pointer is not const and can be used to change
// the stored value.
// func (t * Type) FindOrNull(k KEY) () {
//		v, ok := &m[key]//Not allowed
// }

// Perform a lookup in a map whose values are pointers.
// If the key is present a const pointer to the associated value is returned,
// otherwise a NULL pointer is returned.
// This function does not distinguish between a missing key and a key mapped
// to a NULL value.
func (rm *ResourceMap) FindPtrOrNull(k ResourceID) *rs.ResourceStatus {
	rm.RWMu.RLock()
	defer rm.RWMu.RUnlock()
	v := rm.M[k] // Should be nil for missing keys by default
	return v
}

func (jm *JobMap) FindPtrOrNull(k JobID) *pb.JobDescriptor {
	jm.RWMu.RLock()
	defer jm.RWMu.RUnlock()
	v := jm.M[k]
	return v
}

func (tm *TaskMap) FindPtrOrNull(k TaskID) *pb.TaskDescriptor {
	tm.RWMu.RLock()
	defer tm.RWMu.RUnlock()
	v := tm.M[k]
	return v
}

// Change the value associated with a particular key in a map
// If the key is not present in the map the key and value are inserted,
// otherwise the value is updated to be a copy of the value provided.
// True indicates that an insert took place, false indicates an update.
func (rm *ResourceMap) InsertOrUpdate(k ResourceID, val *rs.ResourceStatus) bool {
	rm.RWMu.Lock()
	defer rm.RWMu.Unlock()
	_, ok := rm.M[k]
	rm.M[k] = val
	return !ok
}

func (jm *JobMap) InsertOrUpdate(k JobID, val *pb.JobDescriptor) bool {
	jm.RWMu.Lock()
	defer jm.RWMu.Unlock()
	_, ok := jm.M[k]
	jm.M[k] = val
	return !ok
}

func (tm *TaskMap) InsertOrUpdate(k TaskID, val *pb.TaskDescriptor) bool {
	tm.RWMu.Lock()
	defer tm.RWMu.Unlock()
	_, ok := tm.M[k]
	tm.M[k] = val
	return !ok
}

// Insert a new key and value into a map.
// If the key is not present in the map the key and value are
// inserted, otherwise nothing happens. True indicates that an insert
// took place, false indicates the key was already present.
func (rm *ResourceMap) InsertIfNotPresent(k ResourceID, val *rs.ResourceStatus) bool {
	rm.RWMu.Lock()
	defer rm.RWMu.Unlock()
	_, ok := rm.M[k]
	if !ok {
		rm.M[k] = val
	}
	return !ok
}

func (jm *JobMap) InsertIfNotPresent(k JobID, val *pb.JobDescriptor) bool {
	jm.RWMu.Lock()
	defer jm.RWMu.Unlock()
	_, ok := jm.M[k]
	if !ok {
		jm.M[k] = val
	}
	return !ok
}

func (tm *TaskMap) InsertIfNotPresent(k TaskID, val *pb.TaskDescriptor) bool {
	tm.RWMu.Lock()
	defer tm.RWMu.Unlock()
	_, ok := tm.M[k]
	if !ok {
		tm.M[k] = val
	}
	return !ok
}

// Perform a lookup in map.
// If the key is present and value is non-NULL then a copy of the value
// associated with the key is made into *val. Returns whether key was present.
func (rm *ResourceMap) FindCopy(k ResourceID, val *rs.ResourceStatus) bool {
	rm.RWMu.RLock()
	defer rm.RWMu.RUnlock()
	v, ok := rm.M[k]
	if ok && (v != nil) {
		*val = *v // since we know that the values are pointers for all maps
	}
	return ok
}

func (jm *JobMap) FindCopy(k JobID, val *pb.JobDescriptor) bool {
	jm.RWMu.RLock()
	defer jm.RWMu.RUnlock()
	v, ok := jm.M[k]
	if ok && (v != nil) {
		*val = *v // since we know that the values are pointers for all maps
	}
	return ok
}

func (tm *TaskMap) FindCopy(k TaskID, val *pb.TaskDescriptor) bool {
	tm.RWMu.RLock()
	defer tm.RWMu.RUnlock()
	v, ok := tm.M[k]
	if ok && (v != nil) {
		*val = *v // since we know that the values are pointers for all maps
	}
	return ok
}

// Test to see if a map contains a particular key.
// Returns true if the key is in the collection.
func (rm *ResourceMap) ContainsKey(k ResourceID) bool {
	rm.RWMu.RLock()
	defer rm.RWMu.RUnlock()
	_, ok := rm.M[k]
	return ok
}

func (jm *JobMap) ContainsKey(k JobID) bool {
	jm.RWMu.RLock()
	defer jm.RWMu.RUnlock()
	_, ok := jm.M[k]
	return ok
}

func (tm *TaskMap) ContainsKey(k TaskID) bool {
	tm.RWMu.RLock()
	defer tm.RWMu.RUnlock()
	_, ok := tm.M[k]
	return ok
}
