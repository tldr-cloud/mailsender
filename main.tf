provider "google" {
  project = "mailsender-288100"
  region  = "us-central1"
}

data "archive_file" "function_archive" {
  type        = "zip"
  source_dir  = "${path.root}/src"
  output_path = "${path.root}/out.gzip"
}

resource "google_pubsub_topic" "publish-newsletter" {
  name   = "publish-newsletter"
}

resource "google_storage_bucket_object" "code" {
  name                = format("%s#%s", "${path.root}/out.gzip", data.archive_file.function_archive.output_md5)
  bucket              = "mailsender-deployment"
  source              = data.archive_file.function_archive.output_path
  content_disposition = "attachment"
  content_encoding    = "gzip"
  content_type        = "application/zip"
}

resource "google_cloudfunctions_function" "mailsender" {
  name                  = "mailsender"
  available_memory_mb   = 128
  runtime               = "go113"
  entry_point           = "PubSubMessageHandler"
  event_trigger {
    event_type = "providers/cloud.pubsub/eventTypes/topic.publish"
    resource   = google_pubsub_topic.publish-newsletter.name
    failure_policy {
      retry = false
    }
  }
  service_account_email = "mailsender@mailsender-288100.iam.gserviceaccount.com"
  timeout               = 30

  source_archive_bucket = google_storage_bucket_object.code.bucket
  source_archive_object = google_storage_bucket_object.code.name
}

resource "google_cloudfunctions_function" "subscribe-handler" {
  name                  = "subscribe-handler"
  available_memory_mb   = 128
  runtime               = "go113"
  entry_point           = "ProcessNewSubscribeMsg"
  trigger_http          = true
  service_account_email = "mailsender@mailsender-288100.iam.gserviceaccount.com"
  timeout               = 30

  source_archive_bucket = google_storage_bucket_object.code.bucket
  source_archive_object = google_storage_bucket_object.code.name
}

resource "google_cloudfunctions_function" "unsubscribe-handler" {
  name                  = "unsubscribe-handler"
  available_memory_mb   = 128
  runtime               = "go113"
  entry_point           = "ProcessUnSubscribeMsg"
  trigger_http          = true
  service_account_email = "mailsender@mailsender-288100.iam.gserviceaccount.com"
  timeout               = 30

  source_archive_bucket = google_storage_bucket_object.code.bucket
  source_archive_object = google_storage_bucket_object.code.name
}

resource "google_cloudfunctions_function" "subscribe-verification-handler" {
  name                  = "subscribe-verification-handler"
  available_memory_mb   = 128
  runtime               = "go113"
  entry_point           = "ProcessNewSubscribeConfirmationMsg"
  trigger_http          = true
  service_account_email = "mailsender@mailsender-288100.iam.gserviceaccount.com"
  timeout               = 30

  source_archive_bucket = google_storage_bucket_object.code.bucket
  source_archive_object = google_storage_bucket_object.code.name
}

resource "google_cloudfunctions_function_iam_member" "invoker-subscribe-verification-handler" {
  project        = google_cloudfunctions_function.subscribe-verification-handler.project
  region         = google_cloudfunctions_function.subscribe-verification-handler.region
  cloud_function = google_cloudfunctions_function.subscribe-verification-handler.name

  role   = "roles/cloudfunctions.invoker"
  member = "allUsers"
}
