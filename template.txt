provider_installation {
  network_mirror {
    url = "https://{{.}}/"
    include = ["registry.terraform.io/*/*"]
  }
  direct {
    exclude = ["registry.terraform.io/*/*"]
  }
}