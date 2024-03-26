terraform {
  required_providers {
    openwrt = {
      source = "hashicorp.com/booya/openwrt"
    }
  }
}

provider "openwrt" {
  host         = "192.168.1.1"
  username     = "root"
  password     = "root"
  insecure_tls = true
}


data "openwrt_network_interface" "lan" {
  name = "lan"
}

data "openwrt_board_info" "board" {
}
