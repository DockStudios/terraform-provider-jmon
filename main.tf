resource "jmon_check" "basic_check" {
  name = "My_Check"

  steps = <<EOF
# Check homepage
- goto: https://www.google.co.uk
- check:
    title: Google
EOF

  interval = 20
  #client = "BROWSER_FIREFOX"
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
