For this search engine, I would start with a **small rule-based date parser**
rather than a general natural-language date library. Your domain has a limited
set of useful expressions:

```text
2024
in 2024
released in 2024
before 2010
after 2020
from 2020 to 2023
between 2020 and 2023
January 2024
2024-03-08
```

A predictable parser is easier to test and prevents date expressions from being
interpreted incorrectly.

## Option 1: Custom rule-based parser — recommended for v1

Represent date filters as half-open ranges:

```go
type DateFilter struct {
    From   *time.Time // inclusive
    Before *time.Time // exclusive
    Start  int        // span in original query
    End    int
    Raw    string
}
```

For example:

```text
"released in 2024"
```

becomes:

```go
From:   2024-01-01
Before: 2025-01-01
```

Then the SQL is consistently:

```sql
release_date >= :from
AND release_date < :before
```

### Useful rules

| Input | Filter |
|---|---|
| `2024` | `>= 2024-01-01 AND < 2025-01-01` |
| `in 2024` | Same |
| `before 2010` | `< 2010-01-01` |
| `after 2020` | `>= 2021-01-01`, or define as `> 2020-01-01` |
| `from 2020 to 2023` | `>= 2020-01-01 AND < 2024-01-01` |
| `between 2020 and 2023` | Define whether the end year is inclusive |
| `March 2024` | `>= 2024-03-01 AND < 2024-04-01` |
| `2024-03-08` | Exact date, or `>=` that date and `<` next day |

For card release dates, year and month ranges are especially useful.

### Basic Go implementation

You can first handle ISO dates and four-digit years:

```go
var yearOnly = regexp.MustCompile(`\b(19|20)\d{2}\b`)
var isoDate = regexp.MustCompile(`\b\d{4}-\d{2}-\d{2}\b`)
```

Then use `time.Parse`:

```go
date, err := time.Parse("2006-01-02", value)
```

For a year:

```go
year, err := strconv.Atoi(value)

from := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
before := time.Date(year+1, time.January, 1, 0, 0, 0, 0, time.UTC)
```

The important part is that the parser returns the consumed span:

```go
type DateMatch struct {
    Filter DateFilter
    Start  int
    End    int
}
```

For:

```text
cards released in 2024
```

the parser could consume either:

```text
"2024"
```

or the full phrase:

```text
"released in 2024"
```

Consuming the full phrase is generally better because it leaves cleaner
search text.

## Option 2: Regex plus a date library

You can use regular expressions to identify the structure and `time.Parse` to
construct dates.

For example:

```go
var yearRange = regexp.MustCompile(
    `(?i)\b(?:from\s+)?(19|20)\d{2}\s+(?:to|through|-)\s+(19|20)\d{2}\b`,
)
```

Then parse the captured years yourself.

This gives you more control than a general date parser while avoiding too much
handwritten date construction.

This is probably the best middle ground:

```text
regex identifies intent
time.Time constructs the actual range
```

## Option 3: `github.com/olebedev/when`

[`when`](https://github.com/olebedev/when) is a Go natural-language date
parser. It can interpret phrases such as:

```text
today
tomorrow
next Friday
last month
```

It may be useful if you eventually support relative expressions such as:

```text
released last year
released this month
```

However, your catalog search mostly needs historical calendar expressions, not
relative dates. You would still need custom handling for domain-specific
meanings such as:

```text
released in 2024
before the base set
```

A general parser can also introduce ambiguity that is undesirable for search.

## Option 4: `github.com/araddon/dateparse`

`dateparse` is useful for parsing many date **formats**, such as:

```text
2024-01-10
01/10/2024
January 10, 2024
```

It is more of a format parser than a natural-language query parser. It will not
by itself understand the complete filter semantics of:

```text
cards released before 2010
```

It can be used after your parser extracts a date-like substring.

## Option 5: Use PostgreSQL for date parsing

You could send extracted values to PostgreSQL:

```sql
SELECT '2024-03-08'::date;
```

But I would not use PostgreSQL to parse the original natural-language query.
The API needs to know:

- Which text was consumed
- What filter was created
- How confident the parse was
- What filter chip to show to the user

Those concerns belong in the application layer.

PostgreSQL should receive the resulting typed values:

```go
from   = 2024-01-01
before = 2025-01-01
```

## Recommended parser pipeline

For this project, I would implement date parsing in this order:

### 1. Exact ISO dates

```text
2024-03-08
```

### 2. Explicit ranges

```text
from 2020 to 2023
between 2020 and 2023
2020-2023
```

### 3. Relative comparison phrases

```text
before 2010
after 2020
since 2022
```

### 4. Month and year expressions

```text
March 2024
March 2024 release
```

### 5. Bare years

```text
2024
```

Only interpret a bare four-digit number as a year if it is in a valid range,
such as:

```go
1800 <= year && year <= currentYear+1
```

That avoids interpreting card numbers or unrelated numbers as dates.

## Define boundary semantics explicitly

The biggest source of bugs is ambiguity around words like `after`, `through`,
and `between`.

For example:

```text
after 2020
```

Could mean:

```sql
release_date > '2020-01-01'
```

or:

```sql
release_date >= '2021-01-01'
```

For year-based catalog searches, I would define:

```text
in 2020       → [2020-01-01, 2021-01-01)
before 2020   → [NULL, 2020-01-01)
after 2020    → [2021-01-01, NULL)
through 2020  → [NULL, 2021-01-01)
```

Use half-open ranges consistently:

```sql
release_date >= from
AND release_date < before
```

This avoids problems with timestamps and end-of-day calculations.

## Recommended choice

Use:

```text
custom regex/token parser
+ time.Parse
+ typed DateFilter
```

Then optionally add `when` later if you want relative dates.

For v1, the parser only needs to support a documented grammar such as:

```text
YEAR
in YEAR
before YEAR
after YEAR
from YEAR to YEAR
between YEAR and YEAR
MONTH YEAR
YYYY-MM-DD
```

That will be more predictable than accepting arbitrary natural-language dates, and every parsed expression can retain its original span for removing it from the leftover search text.
