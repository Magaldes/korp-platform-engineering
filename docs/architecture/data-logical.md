# DATA-LOGICAL — Logical Data Model

- Project: Projeto Korp DevOps Challenge
- Operational revision: 7
- Classification: `declared/proposed/inferred` approved canonical as-built state
- Projection: `erDiagram`

```mermaid
erDiagram
    ENT_REQUEST_AUDIT {
        opaque-id id primary
        timestamp-utc occurred_at
        string method
        string canonical_route
        integer http_status
        number duration_ms
    }
```

## Operational summary

Purpose: Modelar e governar a arquitetura implementada e validada do Projeto Korp, preservando rastreabilidade entre requisitos, decisoes, componentes, containers, dados, processos operacionais e o estado as-built da aplicacao.
Compiled 1 projected items from operational revision 7.
