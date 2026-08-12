# Projeto Korp DevOps Challenge

Projeto de serviço HTTP em Go com NGINX como reverse proxy, observabilidade com Prometheus e Grafana, rede Docker Bridge e provisionamento automatizado via Ansible.

## Baseline do desafio

O baseline original do Projeto Korp permanece preservado:

- serviço HTTP em Go com `GET /projeto-korp` e horário UTC dinâmico;
- NGINX como única entrada HTTP publicada no host, na porta 80;
- aplicação Go disponível apenas na rede Docker interna, na porta 8080;
- métricas Prometheus expostas internamente;
- Prometheus e Grafana para monitoramento e visualização;
- Docker Compose para declaração da stack;
- Ansible como ponto único de provisionamento e validação.

`GET /projeto-korp` continua independente das extensões de dados e não depende de PostgreSQL ou Redis para responder.

## Extensões adicionais do projeto

Após a conclusão do baseline, o projeto recebeu duas camadas adicionais, tratadas explicitamente como **extensões voluntárias** e não como requisitos originais do desafio Korp.

### Persistência — PostgreSQL

PostgreSQL é a fonte persistente e autoritativa dos eventos de auditoria das requisições e das estatísticas agregadas.

A extensão inclui:

- persistência de `request_audit_events`;
- endpoints `/audit-events` e `/request-statistics`;
- migrations versionadas;
- volume persistente do PostgreSQL;
- banco acessível somente pela rede Docker interna, sem porta publicada no host;
- retenção de eventos por 7 dias, executada de forma assíncrona e best-effort;
- falhas de persistência não alteram o contrato de `GET /projeto-korp`.

### Cache opcional — Redis

Redis é uma camada de cache-aside **opcional** para `GET /request-statistics`.

A extensão segue estes princípios:

- PostgreSQL continua sendo a fonte de verdade;
- Redis não armazena dados autoritativos;
- cache com TTL configurável, com default de 60 segundos;
- Redis é volátil, sem volume persistente e sem porta publicada no host;
- cache miss ou indisponibilidade do Redis faz fallback para PostgreSQL;
- um cache hit válido pode responder mesmo se PostgreSQL estiver temporariamente indisponível;
- se não houver cache válido nem PostgreSQL disponível, `/request-statistics` retorna indisponibilidade da dependência;
- `GET /projeto-korp` não utiliza Redis.

## Arquitetura

Canonical Architecture: `ARCH-KORP-001 Revision 7`

A arquitetura canônica atual é:

- Package: `ARCH-KORP-001`
- Revision: `7`
- Fonte canônica: `korp/korp-architecture-v7.json`
- Change set de evolução: `CHG-KORP-006`

As projeções públicas derivadas da arquitetura canônica estão em `docs/architecture/`.

A Revision 7 incorpora o estado as-built validado de PostgreSQL, Redis, retenção, contratos HTTP, automação de readiness/migrations e observabilidade das dependências.

## Estado de entrega

A implementação das extensões de persistência e cache foi validada antes da promoção arquitetural. O repositório público contém código, infraestrutura, automação, testes e documentação necessários para reproduzir e explicar a solução.

Evidências operacionais internas, pacotes de agentes e arquivos com valores locais de ambiente permanecem fora do versionamento público.
