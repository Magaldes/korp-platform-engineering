# BEH-ERROR — Error Handling

- Project: Projeto Korp DevOps Challenge
- Operational revision: 6
- Classification: `declared/proposed/inferred` pre-codebase state
- Projection: `flowchart`

```mermaid
flowchart TB
    err_start(["Client sends HTTP request to NGINX"])
    err_ingress{"NGINX evaluates ingress policy"}
    err_metrics_404["NGINX returns 404 for external metrics route"]
    err_upstream{"NGINX attempts upstream communication"}
    err_502["NGINX returns 502 for unavailable or invalid upstream"]
    err_504["NGINX returns 504 for upstream timeout"]
    err_go_router{"Go router evaluates application route and method"}
    err_go_404["Go returns 404 for unknown application route"]
    err_go_405["Go returns 405 with Allow header for unsupported method"]
    err_handler["Go project handler processes valid request"]
    err_result{"Go evaluates processing result"}
    err_go_200["Go returns 200 application response"]
    err_go_500["Go returns 500 for detected application failure"]
    err_pass["NGINX passes Go response through unchanged"]
    err_end(["Client receives HTTP response"])
    err_start --> err_ingress
    err_ingress -->|"External metrics route"| err_metrics_404
    err_ingress -->|"Public application traffic"| err_upstream
    err_metrics_404 --> err_end
    err_upstream -->|"Unavailable or invalid upstream"| err_502
    err_upstream -->|"Upstream timeout"| err_504
    err_upstream -->|"Connected"| err_go_router
    err_502 --> err_end
    err_504 --> err_end
    err_go_router -->|"Unknown route"| err_go_404
    err_go_router -->|"Unsupported method"| err_go_405
    err_go_router -->|"GET or HEAD project route"| err_handler
    err_go_404 --> err_pass
    err_go_405 --> err_pass
    err_handler --> err_result
    err_result -->|"Success"| err_go_200
    err_result -->|"Application failure"| err_go_500
    err_go_200 --> err_pass
    err_go_500 --> err_pass
    err_pass --> err_end
```

## Operational summary

Purpose: Modelar e governar a arquitetura pre-codebase do desafio Korp, preservando rastreabilidade entre requisitos, criterios de aceite, decisoes, componentes e validacao.
Compiled 35 projected items from operational revision 6.
