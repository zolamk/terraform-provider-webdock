---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "webdock_shell_user Resource - terraform-provider-webdock"
subcategory: ""
description: |-
  
---

# webdock_shell_user (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `password` (String, Sensitive) shell user password
- `public_keys` (List of Number) shell user public keys
- `server_slug` (String) shell user server slug
- `username` (String) shell user username

### Optional

- `group` (String) shell user group
- `shell` (String) shell user shell

### Read-Only

- `created_at` (String) shell user creation datetime
- `id` (String) shell user id
