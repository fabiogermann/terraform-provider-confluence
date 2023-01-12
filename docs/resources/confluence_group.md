---
layout: "confluence"
page_title: "Confluence: confluence_group"
sidebar_current: "docs-confluence-resource-group"
description: |-
  Provides group in Confluence
---

# confluence_group

Provides a access control group on your Confluence site.

## Example Usage

```hcl
resource confluence_group "test_group" {
  name = "a_group_to_test"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the group

## Import

Currenty content can not be imported.
