# BEH-SEQUENCE — Interaction Sequence

- Project: Projeto Korp DevOps Challenge
- Operational revision: 6
- Classification: `declared/proposed/inferred` pre-codebase state
- Projection: `sequenceDiagram`

```mermaid
sequenceDiagram
    actor act_evaluator as Korp Evaluator
    participant ctr_nginx as NGINX Reverse Proxy
    participant ctr_go_service as http-server-projeto-korp
    act_evaluator->>ctr_nginx: GET /projeto-korp on host port 80
    ctr_nginx->>ctr_go_service: Forward GET /projeto-korp to port 8080
    ctr_go_service-->>ctr_nginx: Return JSON with project name and current UTC time
    ctr_nginx-->>act_evaluator: Return HTTP JSON response
```

## Operational summary

Purpose: Modelar e governar a arquitetura pre-codebase do desafio Korp, preservando rastreabilidade entre requisitos, criterios de aceite, decisoes, componentes e validacao.
Compiled 7 projected items from operational revision 6.
