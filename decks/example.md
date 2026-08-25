# Example Deck

## AI

### RLHF

`Reinforcement Learning from Human Feedback`

## MCP

Протокол MCP `Model Context Protocol` позволяет подключать ИИ к коду, документации, таск-трекерам, Google Docs, GitHub и другим реальным источникам данных.


## Network

### FQDN

`Fully Qualified Domain Name`


## Databases

### Transactions Isolation Levels

**Read Uncommitted**
Защищает от Lost Update, не защищает от Dirty Read.

**Read Committed**
Защищает от Lost Update и Dirty Read. Используется по-умолчанию в PostgreSQL и большинстве других БД.

**Repeatable Read**
Защищает от Non-repeatable Reads.

**Serializable**
Защищает от Phantom Reads.


## PostgreSQL

### PostgreSQL upsert

```sql
INSERT INTO my_table (col1, col2)
VALUES ($1, $2)
ON CONFLICT (id) DO
UPDATE SET col1 = $1, col2 = $2;
```

### Селективность условия

- выше, если в выборку попадает мало строк;
- ниже, если в выборку попадает много (почти все) сроки.


## ClickHouse 

### Sparse index

Instead of indexing every row, the primary index has one index entry (mark) per group of rows (granule).

### Select actual value from duplicated rows

Если дубликаты с одинаковым ключом (order by):
```sql
select uid, status
from my_table final
where ...
```

Если дубликаты с разным ключом:
```sql
select uid, 
argMax(status, updated_at) as last_status
from my_table
where ...
group by uid
```

### MergeTree

Семейство движков `MergeTree` позволяют схлопывать данные в фоновом режиме.

Их 7 штук, из них 3 позволяют реализовать update.
`ReplacingMergeTree`
`CollapsingMergeTree`
`VersionedCollapsingMergeTree`

