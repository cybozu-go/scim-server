scim-server
===========

This repository contains a reference server implementation for `github.com/cybozu-go/scim`
which can be used as a simple SCIM server.

# Features

## Implemented features

* Users
  * Retrieve a user by its ID
  * Replace a user by its ID
  * Delete a user by its ID
  * Search users using a query
    * Select attributes to include
    * Select attributes to exclude
* Groups
  * Retrieve a group by its ID
  * Replace a group by its ID
  * Delete a group by its ID
  * Search groups using a query
    * Select attributes to include
    * Select attributes to exclude
* Search both groups and users

## Unimplemented features

* Patch operations
* Bulk operations
* ETags
* /Me endpoint

## Miscellaneous

* More configuration flexibility

# Developing

Tests are defined in a separate repository, `github.com/cybozu-go/scim`, under the
`test` directory.

Much of the `ent` code is generated, and definitions stored in `github.com/cybozu-go/scim`
are used. If there are changes in the SCIM resource definitions in that repository,
you will need to run `go generate`

`ent` Edge definitions are defined in manually maintained files. For example
edges from `User` resources are declared in `ent/schema/user.go`
