/*
Package host contains methods to change host resources.

Copyright 2019 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

        https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package host

import (
        "automation/clients"
        "fmt"
        "strings"
        "time"

        cs "google.golang.org/api/compute/v1"
)

type client interface {
        clients.ComputeServiceInt
}

// Host struct
type Host struct {
        c client
}

// NewHost returns a new snapshot of a specified persistent disk.
func NewHost(c client) *Host {
        return &Host{c: c}
}

// CreateDiskSnapshot creates a snapshot.
func (h *Host) CreateDiskSnapshot(projectID, zone, disk, name string) error {
        rb := &cs.Snapshot{
                Description:       "Snapshot of " + disk,
                Name:              name,
                CreationTimestamp: time.Now().Format(time.RFC3339),
        }
        _, err := h.c.CreateSnapshot(projectID, zone, disk, rb)
        if err != nil {
                return fmt.Errorf("failed to create snapshot: %q", err)
        }
        return nil
}

// ListProjectSnapshot returns a list of snapshots.
func (h *Host) ListProjectSnapshot(projectID string) (*cs.SnapshotList, error) {
        snapshots, err := h.c.ListProjectSnapshots(projectID)
        if err != nil {
                return nil, fmt.Errorf("failed to list snapshots: %q", err)
        }
        return snapshots, nil
}

// ListInstanceDisks returns a list of disk names for a given instance.
func (h *Host) ListInstanceDisks(projectID, zone, instance string) ([]string, error) {
        ds, err := h.c.ListDisks(projectID, zone, instance)
        if err != nil {
                return nil, fmt.Errorf("failed to list disks: %q", err)
        }
        dl := []string{}
        for _, d := range ds.Items {
                if !h.diskBelongsToInstance(d, instance) {
                        return []string{}, nil
                }
                dl = append(dl, d.Name)
        }
        return dl, nil
}

// diskBelongsToInstance returns if the disk is attributed to the given instance.
func (h *Host) diskBelongsToInstance(disks *cs.Disk, instance string) bool {
        for _, u := range disks.Users {
                if !strings.HasSuffix(u, "/instances/"+instance) {
                        continue
                }
                return true
        }
        return false
}

// SetSnapshotLabels sets the labels on a snapshot.
func (h *Host) SetSnapshotLabels(projectID, resource string, m map[string]string) error {
        rb := &cs.GlobalSetLabelsRequest{Labels: m}
        _, err := h.c.SetLabels(projectID, resource, rb)
        if err != nil {
                return fmt.Errorf("failed to set disk labels: %q", err)
        }
        return nil

}