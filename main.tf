
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
}

provider "jmon" {

}

terraform {
  required_providers {
    jmon = {
      source  = "github.com/matthewjohn/jmon"
      # Other parameters...
    }
  }
}
