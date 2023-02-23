resource "jmon_check" "basic_check" {
  name = "My_Check"

  steps = <<EOF
# Check homepage
- goto: https://en.wikipedia.org/wiki/Main_Page
- check:
    title: Wikipedia, the free encyclopedia
EOF

  interval = 200
  client = "BROWSER_FIREFOX"
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
