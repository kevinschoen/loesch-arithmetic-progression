# loesch-arithmetic-progression

Solver for the [IBM Ponder This — September 2026](https://research.ibm.com/blog/ponder-this-september-2026) challenge: find a 35-term arithmetic progression among **Loesch numbers** with the smallest possible last term.

## Problem

[Loesch numbers](https://oeis.org/A003136) (Loeschian numbers) are integers of the form `x² + xy + y²` for integers `x, y`. The first few are:

`0, 1, 3, 4, 7, 9, 12, 13, 16, 19, 21, ...`

**Goal:** Find `(start, step)` such that `start + k·step` for `k = 0, 1, …, 34` are all Loesch numbers, and the last term `start + 34·step` is as small as possible.

Example (4 terms): `1, 7, 13, 19` → start=`1`, step=`6`.

## Approach

1. **Membership test** — A number is Loeschian iff every prime `p ≡ 2 (mod 3)` appears with an even exponent in its factorization ([OEIS A003136](https://oeis.org/A003136)).
2. **Precomputed set** — Build a lookup table of Loesch numbers up to `-max-end` for O(1) checks.
3. **Minimal-end search** — Scan `end = 0, 1, 2, …` in ascending order. For each Loeschian `end`, try all step sizes `d` with `start = end − 34·d ≥ 0`. Return the first valid progression (guaranteed optimal).
4. **Parallelism** — Step candidates for each `end` are checked in parallel across CPU cores.

## Requirements

- Go 1.21+

## Run

```bash
go test ./...
go run ./cmd/solver -terms 35 -verbose
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-terms` | `35` | Number of AP terms |
| `-max-end` | `50000000` | Maximum last term to search |
| `-verbose` | `false` | Print search progress to stderr |

### Example output

```
start=...
step=...
end=...
terms=35
```

## Submit

Send your `(start, step)` to `ponder@il.ibm.com`.

## License

MIT
