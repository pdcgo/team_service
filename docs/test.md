# Testing

Integration/unit tests use **Go + the `moretest` harness against real Postgres**
(`moretest_mock.MockPostgresDatabase(&scenario)` + `DbScenario`), following the repo rule
(never sqlite). Each `moretest.Suite` scenario runs inside **one transaction**, so:

- operations expected to error (e.g. a duplicate-`team_code` insert) go **last** — a
  failed statement aborts the shared transaction;
- run a module with `-p 1` (`cd team_service && go test ./... -p 1`).

Tests live beside the code in `team_v1/` (internal `package team_v1`, so the testable core
`createTeam` is reachable). Handlers are exercised by calling the impl directly
(`NewTeamService(tx)`); the access interceptor is not in the unit path, so no token is needed.

## Coverage

### `team_create_test.go` — `TeamCreate` core (`createTeam`)
- Creates the team (with the `team_code` **uppercased**), a linked `team_infos` row, and a
  `user_team_roles` row for the caller with `ROLE_TEAM_OWNER` / alias `own` — the three
  documented steps, in one transaction.
- Duplicate `team_code` violates the unique index → error (kept last).

### `team_crud_test.go` — `TeamList` / `TeamDetail` / `TeamUpdate` / `TeamDelete`
- **List**: excludes soft-deleted teams; filters by `team_type`; keyword `q` matches
  name / team_code (`ILIKE`); newest-first, paged.
- **Detail**: returns the team with its `TeamInfo`; missing id → `CodeNotFound`.
- **Update**: changes name + description only; `type` and `team_code` stay immutable;
  missing id → `CodeNotFound`.
- **Delete**: soft-deletes (`deleted = true`), which drops the team from the list;
  re-deleting an already-deleted team → `CodeNotFound`.

## Not covered here
- Authorization (the `(role_base.v1.request_policy)` enforcement) is owned and tested by
  `user_service/access_interceptors`; team_service only declares the policies.
- The streaming `TeamCreate` wrapper is a thin adapter over `createTeam` (which is fully
  tested); its progress streaming is a `slog`→`ServerStream` shim.
