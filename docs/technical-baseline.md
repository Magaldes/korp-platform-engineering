# Korp — levantamento técnico e checkpoint de preparação

Status: proposta para revisão humana. Data do levantamento: 2026-08-07.

## Estado encontrado

O diretório inicial estava vazio. Não foram encontrados arquivos Git,
documentação, decisões arquiteturais, código Go, Dockerfile, Compose,
playbooks Ansible, NGINX, Prometheus ou Grafana.

No host foram observados:

| Item | Estado observado |
| --- | --- |
| Sistema | Ubuntu 24.04.4 LTS (Noble) |
| Arquitetura | x86_64/amd64 |
| Git | 2.43.0 |
| Docker Engine | 29.6.0 |
| Docker Compose | v5.2.0, plugin `docker compose` |
| Go | não instalado |
| Ansible | não instalado |
| GNU Make | 4.3 |
| curl | 8.5.0 |

## Dependências propostas

As versões abaixo são recomendações operacionais para reproduzibilidade.
Devem ser revalidadas no momento da implementação e não alteram os contratos
arquiteturais em revisão.

| Tecnologia | Recomendação | Motivo/origem |
| --- | --- | --- |
| Linux | Ubuntu 24.04 LTS amd64 | Ambiente observado; o Docker documenta Noble 24.04 e amd64 como suportados. Necessidade operacional. |
| Go | 1.26.x, fixando o patch vigente no ambiente de build | Release estável atual consultada; compatibilidade de linguagem/runtime e build do serviço. Necessidade operacional. |
| Docker Engine | 29.6.x, fixando o pacote no provisionamento | O host já possui essa linha e a documentação mostra pacote estável para Noble. Necessidade operacional. |
| Docker Compose | plugin v5.2.x, usando a Compose Specification | Já instalado; a especificação é o formato recomendado. Necessidade operacional. |
| Ansible | `ansible-core` 2.21.x em ambiente Python isolado | Linha com ciclo de segurança vigente no levantamento e compatível com Python 3.12–3.14 no control node. Necessidade operacional. |
| NGINX | `nginx:1.30.4` (ou variante Debian explícita) | Versão estável numerada indicada pelo projeto NGINX no levantamento; a imagem oficial deve ser validada antes do pin final. A rota/escopo ainda está em revisão. Necessidade operacional. |
| Prometheus | `prom/prometheus:v3.12.0` | Tag numerada da imagem oficial, multi-arquitetura; evita `latest`. O contrato de observabilidade ainda está em revisão. Necessidade operacional. |
| Grafana | `grafana/grafana:13.0.4-ubuntu` | Tag numerada da imagem oficial, multi-arquitetura; evita `latest`. DEC-005 exige provisionamento, mas não fixa esta versão. Necessidade operacional. |

### Alternativas consideradas

- Ubuntu 22.04 LTS também é suportado pelo Docker; Ubuntu 24.04 é preferível
  aqui por ser o host disponível e permanecer uma base LTS suportada.
- `ansible` completo ou `ansible-core` são opções válidas. Recomenda-se
  `ansible-core` em virtualenv para reduzir dependências; coleções adicionais
  só devem ser adicionadas quando o playbook exigir.
- Imagens `latest`, tags móveis (`v3`, `13.0`) e variantes distroless/alpine
  são possíveis. Recomenda-se tag numerada explícita para facilitar auditoria;
  a variante base ainda pode ser ajustada após validar necessidades de debug.

## Compatibilidade e requisitos operacionais

- Ubuntu 24.04 amd64 é suportado pelo Docker Engine. As imagens consultadas
  de NGINX, Prometheus e Grafana publicam manifestos para amd64; isso não
  substitui a validação por `docker buildx imagetools inspect` no momento da
  implementação.
- Compose v5.2 usa a especificação atual e é compatível com a declaração de
  redes externas prevista por DEC-004. O nome da rede e seu ciclo de vida
  continuam pertencendo ao desenho aprovado, não a este checkpoint.
- O Ansible control node precisa de Python dentro da faixa suportada pelo
  `ansible-core` escolhido. O host Ubuntu 24.04 deve fornecer Python 3.12;
  a presença efetiva deverá ser verificada na máquina que executará a
  automação. O alvo precisa de Python funcional para os módulos usuais e
  acesso ao daemon Docker.
- Para o host: acesso root/sudo para instalar pacotes e manipular Docker,
  `ca-certificates`, `curl`, Python 3, `python3-venv`, `python3-pip`, Docker
  Engine/CLI, Buildx e Compose plugin. Para o build Go: toolchain Go e acesso
  ao proxy/módulos (ou cache equivalente). A lista de pacotes não é ainda um
  contrato de instalação automatizada.
- Portas a reservar/confirmar: serviço Go (porta interna ainda dependente do
  contrato), NGINX 80/443 conforme escopo aprovado, Prometheus 9090 e Grafana
  3000 para demonstração local. A exposição pública/privada de métricas não
  foi decidida; não deve ser inferida desta lista.
- DEC-004 exige que a bridge exista antes do Compose. A ordem de execução e
  o nome concreto devem ser implementados somente no playbook aprovado.

## Estrutura proposta

```text
app/
  cmd/                 entradas executáveis do serviço Go
  internal/            pacotes internos do serviço, sem framework obrigatório
infra/
  docker/              Dockerfile e arquivos específicos da imagem do serviço
  compose/             Compose e arquivos de execução local
  nginx/               configuração do NGINX, após aprovação das rotas
  prometheus/          configuração/regras, após aprovação do contrato
  grafana/              datasource e dashboards versionados (DEC-005)
automation/
  ansible/             inventário, playbooks, roles e defaults da automação
tests/
  unit/                testes unitários do serviço
  integration/         testes de integração do conjunto executável
docs/
  architecture/        decisões, diagramas e registros arquiteturais
  technical-baseline.md levantamento deste checkpoint
evidence/              evidências reproduzíveis de execução/demonstração
```

Essa é uma proposta de organização, não uma decisão sobre a modularização
definitiva. Ela separa artefatos por responsabilidade e deixa os contratos
em revisão fora da implementação inicial.

## Decisões preservadas e limites deste checkpoint

Foram preservadas DEC-004, DEC-005 e DEC-006. Não foram definidos endpoint,
nomes ou labels de métricas; disponibilidade; exposição de `/metrics`; rotas
NGINX; contrato de erros HTTP; endpoints adicionais; persistência, banco,
frontend, autenticação, Kubernetes, Terraform, cloud provider ou CI/CD.

Também permanecem para decisão humana o contrato de observabilidade, escopo
de roteamento NGINX, contrato HTTP/comportamento de erros, estrutura modular
definitiva e estratégia de testes.

## Fontes consultadas

- [Docker Engine no Ubuntu](https://docs.docker.com/engine/install/ubuntu/)
- [Compose Specification](https://docs.docker.com/reference/compose-file/)
- [Instalação do Compose](https://docs.docker.com/compose/install/)
- [Release notes do Go 1.26](https://go.dev/doc/go1.26)
- [Ciclo de manutenção do ansible-core](https://docs.ansible.com/projects/ansible-core/devel/reference_appendices/release_and_maintenance.html)
- [Imagem oficial NGINX](https://hub.docker.com/_/nginx)
- [Download e versões estáveis do NGINX](https://nginx.org/en/download.html)
- [Tags da imagem Prometheus](https://hub.docker.com/r/prom/prometheus/tags)
- [Tags da imagem Grafana](https://hub.docker.com/r/grafana/grafana/tags)
