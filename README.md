# jmon-terraform-provider

Terraform provider for jmon

## Building and Installing

Install golang 1.15+

```
make
```

## Using in terraform

Define the provider:
```hcl
provider "jmon" {
  # This can be ommitted to default to http://localhost:5000
  url = "https://my-jmon-installation.com"
}

terraform {
  required_providers {
    jmon = {
      source  = "github.com/matthewjohn/jmon"
    }
  }
}
```

Create a minimal check:

```hcl
resource "jmon_check" "basic_check" {
  name = "My_Check"

  steps = <<EOF
# Check homepage
- goto: https://en.wikipedia.org/wiki/Main_Page
- check:
    title: Wikipedia, the free encyclopedia
EOF
}
```

An environment with a fully populated check
```hcl
resource "jmon_environment" "production" {
  name = "production"
}

resource "jmon_check" "full_check" {
  name = "Fully_Populated"

  environment         = jmon_environment.production.name
  timeout             = 60
  interval            = 300
  client              = "BROWSER_FIREFOX"
  screenshot_on_error = true

  steps = <<EOF
# Check homepage
- goto: https://en.wikipedia.org/wiki/Main_Page
- check:
    title: Wikipedia, the free encyclopedia
# Perform search
- find:
    id: searchform
    find:
    tag: input
    actions:
        - type: Pabalonium
        - press: enter
- check:
    url: "https://en.wikipedia.org/w/index.php?fulltext=Search&search=Pabalonium&title=Special%3ASearch&ns0=1"
- find:
    class: mw-search-nonefound
    check:
    text: There were no results matching the query.
- actions:
    - screenshot: Homepage
EOF
}

```
