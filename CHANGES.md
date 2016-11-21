# Changes

## v0.2

New features:

* You can now choose to specify Terraform variables on the command line using `--terraform-variable` (if reading them from a file using `--terraform-variables-from` is undesirable)

Breaking changes:

* The `--terraform-variables` command-line argument has been renamed to `--terraform-variables-from`.

Bug fixes:

* Fetching source from a directory now works correctly (#1)

## v0.1

* Initial skeleton
