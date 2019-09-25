/*
Package actions provides the implementation of automated actions.

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
package actions

import (
        "automation/clients"
        "automation/finding"
        "automation/host"

        "context"
        "fmt"
        "strings"
        "time"

        "cloud.google.com/go/pubsub"
)

const (
        snapshotPrefix = "forensic-snapshots-"
        // snapshotTemplate is the name of the snapshot with disk, rule name and time included.
        snapshotTemplate = snapshotPrefix + "%s-%s"
        // allowSnapshotOlderThanDuration defines how old a snapshot must be before we overwrite.
        allowSnapshotOlderThanDuration = time.Minute * 5
)

// supportedRules contains a map of rules this function supports.
var supportedRules = map[string]bool{"bad_ip": true}

/*
   CreateSnapshot creates a snapshot of an instance's disk.
   For a given supported finding pull each disk associated with the affected instance.
   - Check to make sure we haven't created a snapshot for this finding recently.
   - Create a new snapshot for each disk labeled with the finding and current time.
*/

// CreateSnapshot creates a snapshot of an instance's disk.
func CreateSnapshot(ctx context.Context, m pubsub.Message, c clients.ClientInt) error {

        f := finding.NewFinding()
        h := host.NewHost(c)

        if err := f.ReadFinding(&m); err != nil {
                return fmt.Errorf("failed to read finding: %q", err)
        }

        if !supportedRules[f.RuleName()] {
                return nil
        }

        disks, err := h.ListInstanceDisks(f.ProjectID(), f.Zone(), f.Instance())
        if err != nil {
                return fmt.Errorf("failed to list disks: %q", err)
        }

        snapshots, err := h.ListProjectSnapshot(f.ProjectID())
        if err != nil {
                return fmt.Errorf("failed to list snapshots: %q", err)
        }

        for _, disk := range disks {

                rulename := strings.ReplaceAll(f.RuleName(), "_", "-")

                sn := fmt.Sprintf(snapshotTemplate, rulename, disk)
                prefix := snapshotPrefix + rulename + "-" + disk
                createSnapshot := true

                for _, snapshot := range snapshots.Items {
                        if extractDisk(snapshot.SourceDisk) != disk || !strings.HasPrefix(snapshot.Name, prefix) {
                                continue
                        }

                        isSnapshotNew, err := isSnapshotCreatedWithin(snapshot.CreationTimestamp, allowSnapshotOlderThanDuration)
                        if err != nil {
                                return fmt.Errorf("failed to parse snapshot timestamp: %q", err)
                        }
                        if isSnapshotNew {
                                createSnapshot = !isSnapshotNew
                                break
                        }
                }

                if !createSnapshot {
                        continue
                }

                if err = h.CreateDiskSnapshot(f.ProjectID(), f.Zone(), disk, sn); err != nil {
                        return fmt.Errorf("failed to create disk snapshot: %q", err)
                }

                if err = addSnapshotLabels(f.ProjectID(), f.Resource(), disk, f, h); err != nil {
                        return fmt.Errorf("failed to set snapshot labels: %q", err)
                }
        }
        return nil
}

// isSnapshotCreatedWithin checks if the previous snapshots created N mins ago.
func isSnapshotCreatedWithin(snapshotTime string, window time.Duration) (bool, error) {
        t, err := time.Parse(time.RFC3339, snapshotTime)
        if err != nil {
                return false, err
        }
        return time.Since(t) < window, nil
}

// extractDisk extracts the disk from disk resource.
func extractDisk(snapshot string) string {
        return strings.SplitAfter(snapshot, "/disks/")[1]
}

// addSnapshotLabels sets the labels of snapshots.
func addSnapshotLabels(ProjectID, resource, disk string, f *finding.Finding, h *host.Host) error {
        ips := f.BadIPs()
        b := strings.Join(ips, ",")

        // labelMax is the maximum size of labels we support.
        const labelMax = 60

        if len(b) > labelMax {
                b = b[:labelMax] + "..."
        }

        m := map[string]string{
                "disk-name":         disk,
                "network-indicator": b,
        }
        if err := h.SetSnapshotLabels(ProjectID, resource, m); err != nil {
                return fmt.Errorf("failed to set disk labels: %q", err)
        }
        return nil
}