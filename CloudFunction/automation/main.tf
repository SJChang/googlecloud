locals {
  zone           = "us-central1-a"
  region         = "us-central1"
  findings-topic = "threat-findings"
}

provider "google" {
  project = "${var.automationProject}"
  region  = "${local.region}"
}