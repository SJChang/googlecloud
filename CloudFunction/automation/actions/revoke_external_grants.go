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
        "automation/user"

        "context"
        "fmt"

        "cloud.google.com/go/pubsub"
)

/*
RevokeExternalGrants is the entry point of the Cloud Function.

This Cloud Function will read the incoming finding, if it's an ETD anomalous IAM grant
identicating an external member was invited to policy check to see if the external member
is in a list of disallowed domains.

Additionally check to see if the affected project is in the specified folder. If the grant
was to a domain explicitly disallowed and within the folder then remove the member from the
entire IAM policy for the resource.

TODO:
  - Disallowed email list should be an argument.
  - We currently remove any member that matches the disallowed emails. May want to only remove
        if they are explicitly found from a detector. Currently we'll remove an existing member
        that may not be intended.
*/
func RevokeExternalGrants(ctx context.Context, m pubsub.Message, c clients.ClientInt, folderIDs []string, disallowed []string) error {
        f := finding.NewFinding()

        if err := f.ReadFinding(&m); err != nil {
                return fmt.Errorf("failed to read finding: %q", err)
        }

        if eu := f.ExternalUsers(); len(eu) == 0 {
                return fmt.Errorf("no external users")
        }

        ancestors, err := c.GetProjectAncestry(f.ProjectID())
        if err != nil {
                return fmt.Errorf("failed to get project ancestry: %q", err)
        }

        for _, resource := range ancestors {
                for _, folderID := range folderIDs {
                        if resource != "folders/"+folderID {
                                continue
                        }

                        r := user.NewUser(c)
                        _, err = r.RemoveDomainsProject(f.ProjectID(), disallowed)
                        if err != nil {
                                return fmt.Errorf("failed to remove disallowed domains: %q", err)
                        }
                        return nil
                }
        }
        return nil
}