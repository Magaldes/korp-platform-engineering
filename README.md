# korp-platform-engineering

Platform engineering project built with Go, Docker, NGINX, Prometheus, Grafana and Ansible, featuring automated provisioning, observability and architecture-driven delivery.

## Projeto Korp

Serviço HTTP em Go com proxy NGINX, métricas Prometheus, dashboard Grafana,
rede Docker bridge e provisionamento automatizado por Ansible.

## Estado deste checkpoint

Este repositório contém a implementação validada, os arquivos de execução,
testes, documentação técnica e evidências finais. A baseline arquitetural
aprovada é registrada em `agent-packages/` apenas no workspace interno e não é
publicada neste repositório.

Decisões aprovadas:

- DEC-004: a bridge Docker é criada pelo Ansible e consumida pelo Compose como rede externa.
- DEC-005: datasource e dashboard do Grafana são artefatos versionados e provisionados automaticamente.
- DEC-006: `horario` usa timestamp UTC em RFC3339.
