### LIKE wildcards
`_` matches exactly 1 character.
`%` matches 0 or more characters.

### Transaction isolation, read phenomena
`Lost Update` transaction 2 overwrites changes made by transaction 1 (last write wins).
`Dirty Read` transaction 1 makes changes, and transaction 2 reads them before transaction 1 commits.
`Non-repeatable Read` (read skew) transaction 1 reads data, transaction 2 modifies and commits a row, and transaction 1 reads the same data again and sees the change.
`Phantom Read` transaction 1 reads data, transaction 2 inserts and commits new data, and transaction 1 repeats the same query and sees the new data.

### Transaction isolation levels
`Read Uncommitted` prevents lost updates but does not prevent dirty reads.
`Read Committed` prevents lost updates and dirty reads.
`Repeatable Read` prevents non-repeatable reads.
`Serializable` prevents phantom reads.


### Preventing lost updates
`Atomic operations` for different data types: counters, queues, and JSON documents.
Locks using `SELECT FOR UPDATE`.
Automatic detection of `lost updates` and `transaction rollback`.
`Compare-and-set`.

## SQL standard
SQL was first officially adopted as an `ANSI` standard in `1986`.

## SQL is an abstraction that hides:
- complex `on-disk and` `in-memory` `data structures`;
- `concurrent requests` from other clients;
- `inconsistencies` after crashes.

## Redis
Redis is an acronym that stands for `REmote DIctionary Server`.

### Redis persistence
The `RDB` (`Redis Database`) persistence performs point-in-time snapshots of your dataset at specified intervals.
The `AOF` (`Append Only File`) persistence logs every write operation received by the server, that will be played again at server startup, reconstructing the original dataset.



