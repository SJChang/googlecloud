variable "automationProject" {
  type        = "string"
  description = "Project ID of cloud function automation."
}

variable "threatfindingsProject" {
  type        = "string"
  description = "Project ID that generates the threat findings."
}

variable "userFolder" {
  type        = "string"
  description = "Folder ID that external users are added."
}