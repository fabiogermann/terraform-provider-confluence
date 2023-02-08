resource "confluence_space_permission" "test_space_permissions" {
  key = "TST"
  operations = [
    "create:comment",
    "read:space",
  ]
  group = "group_name"
}