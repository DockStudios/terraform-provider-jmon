
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

resource "jmon_check" "example_plugin_call" {
  name = "Example_Plugin_Call"

  environment = "default"

  steps = <<EOF
# Call "example-plugin" plugin, passing in example_argument argument
 - call_plugin:
     example-plugin:
       example_argument: 'example_value'

# Call example plugin a second time
 - call_plugin:
     example-plugin:
       example_argument: 'example_value'

# Go to URL that contains run variable, set by plugin
 - goto: https://example.com/{variable_set_by_example_plugin}
EOF

  attributes = {
    example_attribute = "Test value"
  }
}


resource "jmon_check" "api_check" {
  name = "Check_DummyJson_Api"

  environment = "default"
  interval    = 300

  steps = <<EOF
- goto:
    url: https://dummyjson.com/products/1
    headers:
      Content-Type: application/json
    body: {id: 1}
    method: get
- check:
    json:
      selector: '$.id'
      equals: 1
- check:
    json:
      selector: '$.images[0]'
      contains: 1.jpg
- check:
    json:
      # Check entire response - can be provided as JSON value or YAML
      equals: {"id":1,"title":"iPhone 9","description":"An apple mobile which is nothing like apple","price":549,"discountPercentage":12.96,"rating":4.69,"stock":94,"brand":"Apple","category":"smartphones","thumbnail":"https://i.dummyjson.com/data/products/1/thumbnail.jpg","images":["https://i.dummyjson.com/data/products/1/1.jpg","https://i.dummyjson.com/data/products/1/2.jpg","https://i.dummyjson.com/data/products/1/3.jpg","https://i.dummyjson.com/data/products/1/4.jpg","https://i.dummyjson.com/data/products/1/thumbnail.jpg"]}
EOF
}

provider "jmon" {
  url     = "http://localhost:5000"
}

terraform {
  required_providers {
    jmon = {
      # Use provider version from Hashicorp registry
      source  = "dockstudios/jmon"
      # Other parameters...

      # Use for building locally
      #source  = "github.com/dockstudios/jmon"
    }
  }
}

