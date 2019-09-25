package host

import (
        "automation/clients"
        "testing"

      cs "google.golang.org/api/compute/v1"

      "gopkg.in/d4l3k/messagediff.v1"
)

const (
        projectID = "test-project"
        zone      = "test-zone"
        disk      = "test-disk"
)

func TestCreateDiskSnapshot(t *testing.T) {
        tests := []struct {
                name             string
                snapshotName string
                expectedError    error
        }{
                {
                        name:             "create disk snapshot",
                        snapshotName: "test-snapshot",
                        expectedError:    nil,
                },
        }
        for _, tt := range tests {
                t.Run(tt.name, func(t *testing.T) {
                        mock := &clients.MockClients{}
                        s := NewHost(mock)
                        if err := s.CreateDiskSnapshot(projectID, zone, disk, tt.snapshotName); err != tt.expectedError {
                                t.Errorf("%v failed exp:%v got: %v", tt.name, tt.expectedError, err)
                        }
                })
        }
}

func TestListProjectSnapshot(t *testing.T){
    tests := []struct {
                name             string
                input             *cs.SnapshotList
                expectedError    error
                expectedOutput             *cs.SnapshotList
        }{
                {
                        name:             "List project snapshot",
                        input:             createSnapshots(),
                        expectedError:    nil,
                        expectedOutput: createSnapshots(),
                },
        }
        for _, tt := range tests {
                t.Run(tt.name, func(t *testing.T) {
                        mock := &clients.MockClients{}
                        s := NewHost(mock)
                        mock.AddSavedCreateSnapshotFake(tt.input)
                        ss, err := s.ListProjectSnapshot(projectID)
                        if diff, equal := messagediff.PrettyDiff(ss, tt.expectedOutput); !equal {
                                t.Errorf("%v failed, difference: %+v", tt.name, diff)
                        }
                        if err != tt.expectedError {
                                t.Errorf("%v failed exp:%v got: %v", tt.name, tt.expectedError, err)
                        }
                })
        }
}

func TestListInstanceDisks(t *testing.T){
}

func ListProjectInstances(t *testing.T){
}


func createSnapshots() *cs.SnapshotList {
        return *cs.SnapshotList{
                        CreationTimestamp:    "2019-08-01T22:17:07.159-07:00",
                        Name:"snapshot-2",
                        SourceDisk:"https://www.googleapis.com/compute/v1/projects/regal-height-244217/zones/us-central1-a/disks/disk-1",
        }
}