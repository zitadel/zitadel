# What is this

This folder contains migrations copied from `cmd/setup/` and `cmd/initialise/` that are needed to create a minimal eventstore setup on DB, which is then used to test event reduction on relational tables.

**These migrations are only for testing purposes,** they are not run on main DB but an embedded one and only as part of unit tests.