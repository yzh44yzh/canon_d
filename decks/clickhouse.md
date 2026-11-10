# Clickhouse

## LSM-tree
`Log-Structured Merge-Tree`

## SSTable
`Sorted String Table`

## MergeTree
A MergeTree in ClickHouse is a table engine designed for high data ingest rates and large data volumes.

It is the core storage engine in ClickHouse, providing features such as
- `columnar storage`,
- `custom partitioning`,
- `sparse primary indexes`,
- and support for `background data merges`.

7 engines, 3 of them allow update:
- `ReplacingMergeTree`
- `CollapsingMergeTree`
- `VersionedCollapsingMergeTree`

## Output formats
select * from table format Vertical

Also:
- `JSON`
- `JSONEachRow`
- `CSV`
- `Pretty`
- `PrettyCompact` (default in interactive mode)
- `Markdown`
- and many other

## SQL

### Describe Table
`show create table my_table`

### Update Row
`alter table` my_table
`update field1 = 'value1', field2 = 'value2'`
where uid = '123'

### Delete Row
`alter table` my_table
`delete` where uid = '123'

### Select with deduplication
PKs are the same:
select uid, status
`from my_table final`
where ...

PKs differ:
select uid, 
`argMax(status, updated_at) as last_status`
from my_table
where ...
`group by uid`

### Count All And Uniq Rows
`select count()` from my_table where ...
`select uniq(uid)` from my_table where ...

### Fill table with data
INSERT INTO my_table(field_1, field_2)
SELECT 
  `number`,
  `concat('prefix_', toString(rand())`
FROM `system.numbers`
LIMIT 1000000

### Table Size
SELECT 
  table, 
  `sum(bytes)` as total_bytes,
  `formatReadableSize(sum(bytes))` as readable_bytes
FROM `system.parts`
WHERE active AND database = 'my_db'
ORDER BY total_bytes DESC
GROUP BY table

### View The Table Parts And Partitions
SELECT `partition`, name, `active`
FROM `system.parts`
WHERE table = 'my_table'
