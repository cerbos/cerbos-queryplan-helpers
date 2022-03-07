# Cerbos Query Plan helpers

---
The Query Planner API aims to address a common use case in access control scenarios: filtering a list to retain only those items that a particular user has access to. You can send a principal, action, resource kind and any known resource attributes  to `/api/x/plan/resources` and obtain a tree representation, abstract syntax tree (AST), of conditions that must be satisfied in order for that principal to be allowed to perform the action on that resource kind. See [Resources Query Plan](https://docs.cerbos.dev/cerbos/latest/api/index.html#resources-query-plan) for details about the new API.

This project contains Go modules showing how to use the Query Planner API to filter data using various ORMs and toolkits.

## Go Modules
Current list of modules:
- [ent-adapter](https://github.com/cerbos/cerbos-queryplan-helpers/tree/main/ent-adapter) is an example and helper functions using [Ent](https://entgo.io/) ORM.
- [pgx-adapter](https://github.com/cerbos/cerbos-queryplan-helpers/tree/main/pgx-adapter) is an example and helper function using [pgx](https://github.com/jackc/pgx) - PostgreSQL Driver and Toolkit.

