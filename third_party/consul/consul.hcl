datacenter = "yone"
data_dir = "/data/consul/data"

server = true
bootstrap_expect = 1
bind_addr = "0.0.0.0"
client_addr = "0.0.0.0"

ui_config {
  enabled = true
}

node_name = "node"
bootstrap_expect = 1

connect {
  enabled = false
}

acl = {
  enabled = true
  default_policy = "deny"
  enable_token_persistence = true
  tokens = {
    agent = "923a6a81-dd6f-5b22-b0f2-0f7da6fda489"
  }
}