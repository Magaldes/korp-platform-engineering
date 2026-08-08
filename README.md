# Projeto Korp DevOps Challenge

Projeto de serviço HTTP que inclui um proxy NGINX, observabilidade de métricas do Prometheus em um dashboard do Grafana, uma rede Docker Bridge e provisionamento automatizado via Ansible.

## Estado deste checkpoint

Este repositório contém a implementação validada, os arquivos de execução, testes, documentação técnica e evidências finais.

Decisões aprovadas:

- DEC-004: a bridge Docker é criada pelo Ansible e consumida pelo Compose como rede externa.
- DEC-005: datasource e dashboard do Grafana são artefatos versionados e provisionados automaticamente.
- DEC-006: `horario` usa timestamp UTC em RFC3339.

## Documentação e artefatos locais

A documentação arquitetural do projeto está em `docs/architecture/`.
Evidências de execução e artefatos operacionais internos permanecem somente no workspace local.
