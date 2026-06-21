# Security Policy

## Reporting a vulnerability

Please report security issues **privately** via GitHub Security Advisories
("Report a vulnerability" on the repository's *Security* tab) rather than opening
a public issue. We aim to acknowledge reports within a few business days.

## Supported versions

The latest `v0.0.x` tag receives security fixes. Because Go module versions are
immutable, fixes ship as a new tag (never by moving an existing one).

## Security model & configurable controls

rkit favors **secure defaults that stay configurable**. The controls most
relevant to security:

| Area | Control | Secure default | Notes |
|------|---------|----------------|-------|
| Repository field names (SQL/NoSQL injection) | `PostgresRepository.FieldPolicy` / `MongoRepository.FieldPolicy` (`criteria.FieldPolicy`) | Strict identifier validation (`criteria.IsValidIdentifier`) | Set `FieldPolicy.Allowed` to a column allow-list when field names come from client input. |
| HTTP response size (memory DoS) | `http.Config.MaxResponseBytes` | 10 MiB (`DefaultMaxResponseBytes`) | Reads beyond the limit fail with `ErrResponseTooLarge`. `<0` opts into unlimited. |
| HTTP timeout (resource exhaustion) | `http.Config.Timeout` | 30s (`DefaultTimeoutSeconds`) | A bare `Config{}` is never left without a timeout. `<0` opts into none. |
| Loki log shipping (memory DoS / blocking) | `hooks.LokiBufferedHook.MaxBufferSize`, `Timeout`, `OnError` | buffer cap 10000, 5s timeout | Network POST happens outside the lock; oldest entries are dropped (and counted via `Dropped()`) when Loki is unreachable. |

## Handling guidance for consumers

- `kerrors.KError.Error()` (and `infrastructure/http.Fail`) is **client-facing**.
  Keep secrets, SQL and stack detail out of `Message`; put internal context in
  `Metadata`/`Cause`, which surface through `Detail()` in logs only.
- Loki/loguru fields become Loki **labels**. Keep them low-cardinality and free
  of credentials/PII.
- When exposing generic filter/sort endpoints, configure `FieldPolicy.Allowed`.
