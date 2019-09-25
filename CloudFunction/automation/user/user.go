/*
Package user contains methods to change user resources.

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
package user

import (
        "automation/clients"

        "fmt"
        "regexp"
        "strings"

        stg "cloud.google.com/go/storage"
        crm "google.golang.org/api/cloudresourcemanager/v1"
)

type client interface {
        clients.CloudResourceManagerInt
        clients.StorageInt
}

// User struct
type User struct {
        c client
}

// NewUser returns a new instance of Useu.
func NewUser(c client) *User {
        return &User{c: c}
}

// RemoveDomainsProject removes all members from the given resource that end with the disallowed domains.
func (u *User) RemoveDomainsProject(projectID string, disallowedDomains []string) (*crm.Policy, error) {
        domains := strings.Replace(strings.Join(disallowedDomains, "|"), ".", `\.`, -1)
        regex := regexp.MustCompile(fmt.Sprintf(`@(%s)$`, domains))

        resp, err := u.c.GetPolicyProject(projectID)
        if err != nil {
                return nil, fmt.Errorf("failed to get project policy: %q", err)
        }

        p := u.removeMembersFromPolicy(regex, resp)

        setp, err := u.c.SetPolicyProject(projectID, p)
        if err != nil {
                return nil, fmt.Errorf("failed to set project policy: %q", err)
        }
        return setp, nil
}

// RemoveMembersProject removes the given members.
func (u *User) RemoveMembersProject(projectID string, disallowedUserEmails []string) (*crm.Policy, error) {
        joined := strings.Replace(strings.Join(disallowedUserEmails, "|"), ".", `\.`, -1)
        emailInRegex := regexp.MustCompile(fmt.Sprintf("^%s$", joined))

        resp, err := u.c.GetPolicyProject(projectID)

        if err != nil {
                return nil, fmt.Errorf("failed to get project policy: %q", err)
        }

        p := u.removeMembersFromPolicy(emailInRegex, resp)
        setp, err := u.c.SetPolicyProject(projectID, p)
        if err != nil {
                return nil, fmt.Errorf("failed to set project policy: %q", err)
        }
        return setp, nil
}

// removeMembersFromPolicy removes members that match the given regex.
func (u *User) removeMembersFromPolicy(regex *regexp.Regexp, policy *crm.Policy) *crm.Policy {
        for _, b := range policy.Bindings {
                members := []string{}
                for _, m := range b.Members {
                        if !regex.MatchString(m) {
                                members = append(members, m)
                        }
                }
                b.Members = members
        }
        return policy
}

// RemoveEntityFromBucket removes ACL Entity in the bucket.
func (u *User) RemoveEntityFromBucket(bucketName string, entity stg.ACLEntity) error {
        if err := u.c.RemoveBucketUsers(bucketName, entity); err != nil {
                return fmt.Errorf("failed to remove entity: %q", err)
        }
        return nil
}