terraform {
  required_providers {
    openwrt = {
      source = "hashicorp.com/booya/openwrt"
    }
  }
}

provider "openwrt" {}

data "openwrt_example" "example" {}
