# SORM

Smooth ORM, or Saft ORM, or Super ORM, what ever.

# Target

- Code first, but not SQL first.
- Type-safe, as much as possible.
- No literal strings to build SQl, as less as possible.
- No reflection in SQL building and running.

# TODO

- [ ] SQL builder
  - [ ] Generate schema from model
  - [ ] Indexes, constraints
  - [ ] Associations (has 1, has many, belongs to, many to many)
  - [ ] Transaction
  - [ ] Upsert, Locking
- [ ] Schema migration
  - [ ] Dry run with plans
  - [ ] Data migration
- [ ] Default dialects
  - [ ] SQLite
  - [ ] PostgreSQL
  - [ ] MySQL/MariaDB
- [ ] Logging
  - [ ] User-providing `slog.Logger` instance
  - [ ] Support dependency-injection framework
- [ ] Context
- [ ] Plugins
