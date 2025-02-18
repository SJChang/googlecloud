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
resource "google_cloudfunctions_function" "function" {
  name                  = "SnapshotDisk"
  description           = "Revokes IAM Event Threat Detection anomalous IAM grants."
  runtime               = "go112"
  available_memory_mb   = 128
  source_archive_bucket = "${google_storage_bucket.cloud_function_bucket.name}"
  source_archive_object = "${google_storage_bucket_object.cloud_function_zip.name}"
  timeout               = 60
  project               = "${var.automationProject}"
  region                = "${local.region}"
  entry_point           = "SnapshotDisk"

  event_trigger = {
    event_type = "providers/cloud.pubsub/eventTypes/topic.publish"
    resource   = "${local.findings-topic}"
  }
}

resource "google_storage_bucket" "cloud_function_bucket" {
  name       = "${var.automationProject}-function-bucket-finding"
  depends_on = ["local_file.cloudfunction-key-file"]
}

resource "google_storage_bucket_object" "cloud_function_zip" {
  name   = "cloud_function.zip"
  bucket = "${google_storage_bucket.cloud_function_bucket.name}"
  source = "${path.root}/cloud_function.zip"
}

data "archive_file" "cloud_function_zip" {
  type        = "zip"
  source_dir  = "${path.root}/automation"
  output_path = "${path.root}/cloud_function.zip"
  depends_on  = ["local_file.cloudfunction-key-file"]
}