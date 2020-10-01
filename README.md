This directory contains spannertest, an in-memory fake Cloud Spanner. A sibling
directory, spansql, contains types and parser for the Cloud Spanner SQL dialect.

spansql is reusable for anything that interacts with Cloud Spanner on a
syntactic basis, such as tools for handling Spanner schema (DDL).

spannertest builds on spansql for testing code that uses Cloud Spanner client
libraries.

Neither of these packages aims to be performant nor exact replicas of the
production Cloud Spanner. They are reasonable for building tools, or writing
unit or integration tests. Full-scale performance testing or serious workloads
should use the production Cloud Spanner instead.

Here's a list of features that are missing or incomplete. It is roughly ordered
by ascending esotericism:

- expression functions
- more aggregation functions
- more joins types (INNER, CROSS, FULL, RIGHT)
- INSERT/UPDATE DML statements
- SELECT HAVING
- case insensitivity
- alternate literal types (esp. strings)
- STRUCT types
- transaction simulation
- expression type casting, coercion
- subselects
- set operations (UNION, INTERSECT, EXCEPT)
- partition support
- conditional expressions
- table sampling (implementation)
