---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "webdock_server Resource - terraform-provider-webdock"
subcategory: ""
description: |-
  
---

# webdock_server (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `image_slug` (String) Server image
- `location_id` (String) Location ID of the server
- `name` (String) Server name
- `profile_slug` (String) Server profile

### Optional

- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))
- `virtualization` (String) Virtualization type for your new server. container means the server will be a Webdock LXD VPS and kvm means it will be a KVM Virtual machine. If you specify a snapshotId in the request, the server type from which the snapshot belongs much match the virtualization selected. Reason being that KVM images are incompatible with LXD images and vice-versa.

### Read-Only

- `aliases` (List of String) Server description (what's installed here?) as entered by admin in Server Metadata
- `created_at` (String) Creation date/time
- `id` (String) The ID of this resource.
- `ipv4` (String) IPv4 address
- `ipv6` (String) IPv6 address
- `slug` (String) Server slug
- `snapshot_runtime` (Number) Last knows snapshot runtime (seconds)
- `ssh_password_auth_enabled` (Boolean) Whether SSH password authentication is enabled
- `status` (String) Server status
- `webserver` (String) Webserver type (apache, nginx, none)
- `wordpress_lockdown` (Boolean) Whether WordPress is in lockdown mode

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
