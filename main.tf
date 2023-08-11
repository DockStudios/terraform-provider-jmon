
resource "jmon_environment" "production" {
  name = "production"
}

resource "jmon_check" "prod_check" {
  name = "Blank_Check"
  environment = jmon_environment.production.id

  steps = <<EOF
- actions:
   - screenshot: blank
EOF
}

resource "jmon_check" "full_check" {
  name = "Check_Google_Homepage"

  environment = "default"
  timeout = 30

  steps = <<EOF
# Check homepage
- goto: https://www.google.co.uk
- check:
    title: Google
- actions:
  - screenshot: Homepage
EOF

  interval = 20

  attributes = {
    notification_slack_channel = "test"
  }
}


resource "jmon_check" "api_check" {
  name = "Check_DummyJson_Api"

  environment = "default"

  steps = <<EOF
- goto: https://dummyjson.com/products/1
- check:
    json:
      selector: '.id'
      equals: 1
    json:
      selector: '.images[0]'
      contains: 1.jpg
EOF
}

provider "jmon" {
  url     = "http://localhost:5000"
}

terraform {
  required_providers {
    jmon = {
      source  = "github.com/dockstudios/jmon"
      # Other parameters...
    }
  }
}
