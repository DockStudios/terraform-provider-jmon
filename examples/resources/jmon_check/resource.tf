resource "jmon_check" "simple-example" {
  name = "simple-example-check"

  environment = "default"

  steps = <<EOF
- goto: https://example.com
EOF

}

resource "jmon_check" "full-example" {
  name        = "simple-example-check"
  environment = "default"

  steps = <<EOF
# Example check
- goto: https://www.example.com

# ... Further steps
# See https://github.com/DockStudios/jmon/blob/main/docs/step_reference.md for more information
EOF

  client   = "BROWSER_FIREFOX"
  enabled  = false
  timeout  = 30
  interval = 600

  screenshot_on_error = false

  attributes = {
    notification_slack_channel = "test"
  }
}

