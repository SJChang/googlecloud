# Package automation contains the Cloud Function code to automate actions.

# Copyright 2019 Google LLC

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#       https://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
resource "google_service_account" "automation-service-account" {
  account_id   = "automation-service-account"
  display_name = "Service account used by automation Cloud Function"
  project      = "${var.automationProject}"
}

// Role "compute.instanceAdmin" required to get disk lists and create snapshots for GCE instances.
// This can be applied either in folder level or project level.
// This is used in actions/create-snapshot.go which creates new snapshots of the disk in the event
// of certain detectors triggering. These snapshots can help analysis of the event as the disk is
// captured at the time the activity occurred. This binding can be removed if the action is not
// being used.
resource "google_project_iam_binding" "gce-snapshot-bind" {
  project = "${var.automationProject}"
  role    = "roles/compute.instanceAdmin.v1"
  members = ["serviceAccount:${google_service_account.automation-service-account.email}"]
}

resource "google_folder_iam_binding" "cloudfunction-folder-bind" {
  folder  = "folders/${var.userFolder}"
  role    = "roles/resourcemanager.folderAdmin"
  members = ["serviceAccount:${google_service_account.automation-service-account.email}"]
}

resource "google_service_account_key" "cloudfunction-key" {
  service_account_id = "${google_service_account.automation-service-account.name}"
}

resource "local_file" "cloudfunction-key-file" {
  content  = "${base64decode(google_service_account_key.cloudfunction-key.private_key)}"
  filename = "./automation/auth.json"
}