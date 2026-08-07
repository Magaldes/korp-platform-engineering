# Projeto Korp

Repositório preparado para a implementação do serviço, infraestrutura local,
automação e observabilidade após a aprovação do gate arquitetural.

## Estado deste checkpoint

Este commit contém somente a estrutura mínima, o `.gitignore` e o levantamento
técnico em [`docs/technical-baseline.md`](docs/technical-baseline.md). A lógica
do serviço, os contratos ainda em revisão e os arquivos de execução não foram
implementados.

Decisões já aprovadas que deverão ser respeitadas na implementação:

- DEC-004: a bridge Docker será criada pelo Ansible e consumida pelo Compose
  como rede externa.
- DEC-005: datasource e dashboard do Grafana serão artefatos versionados e
  provisionados automaticamente.
- DEC-006: `horario` usará timestamp UTC em RFC3339.
