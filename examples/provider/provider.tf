provider "jmon" {
  url = "https://my-jmon-instance.example.com"

  # Local instance
  #url = "http://localhost:5000"
}

terraform {
  required_providers {
    jmon = {
      source  = "dockstudios/jmon"
    }
  }
}
