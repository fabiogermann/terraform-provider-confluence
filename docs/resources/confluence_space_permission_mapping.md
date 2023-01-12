[confluence_space_permission_mapping.md](confluence_space_permission_mapping.md)---
layout: "confluence"
page_title: "Confluence: confluence_space_permission_mapping"
sidebar_current: "docs-confluence-resource-space_permission_mapping"
description: |-
  Provides space permission mappings in Confluence
---

# confluence_space_permission_mapping

Provides a mapping between space permissions and a access control group on your Confluence site.

## Example Usage

```hcl
resource confluence_space_permission "read_permission" {
  key = "MYSPACE"
  operations = ["read:space", "create:comment"]
  group = confluence_group.test_group.name
}
```

## Argument Reference

The following arguments are supported:

* `key` - (Required) The space key of the confluence space

* `operations` - (Required) The permissions that should be associated with the access control group

* `group` - (Required) The group associated with the permission set

### Available operations
NOTE: the `"read:space"` operation is mandatory for most constellations and (due to a bug) may not be listed in the
first position in the operations list.
* `"create:page"`, `"create:blogpost"`, `"create:comment"`, `"create:attachment"`
* `"read:space"`
* `"delete:space"`, `"delete:page"`, `"delete:blogpost"`, `"delete:comment"`, `"delete:attachment"`
* `"export:space"`
* `"administer:space"`
* `"archive:page"`
* `"restrict_content:space"`

## Import

Currenty content can not be imported.

