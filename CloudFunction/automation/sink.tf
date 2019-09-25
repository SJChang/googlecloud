resource "google_logging_project_sink" "sink" {
  name                   = "sink-threat-findings"
  destination            = "pubsub.googleapis.com/projects/${var.automationProject}/topics/${local.findings-topic}"
  filter                 = "resource.type = threat_detector"
  unique_writer_identity = true
  project                = "${var.threatfindingsProject}"
}

resource "google_project_iam_binding" "log-writer-pubsub" {
  role    = "roles/pubsub.publisher"
  project = "${var.automationProject}"

  members = [
    "${google_logging_project_sink.sink.writer_identity}",
  ]
}

resource "google_pubsub_topic" "topic" {
  name = "${local.findings-topic}"
}