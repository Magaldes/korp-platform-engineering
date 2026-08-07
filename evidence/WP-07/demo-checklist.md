# WP-07 demonstration checklist

- Provisionamento único: `ansible-playbook -i automation/ansible/inventory.ini automation/ansible/site.yml`
- Log completo e exit code: [ansible-run.log](ansible-run.log) — exit code `0`; o playbook exibiu o JSON HTTP 200.
- Endpoint final: `curl http://localhost:80/projeto-korp`; resposta e horários dinâmicos: [http-response.json](http-response.json).
- Topologia Docker, rede bridge e isolamento de HostPorts: [runtime-inspection.json](runtime-inspection.json).
- Prometheus target/up/counter, tráfego via NGINX e URLs internas: [observability.json](observability.json).
- Grafana para inspeção humana: `http://172.22.0.3:3000`.
- Dashboard: UID `korp-observability`, título `Projeto Korp - Observability`; painéis `Availability`, `Request Total` e `Request Rate`.
- Integração: `tests/integration/http_proxy.sh` contra `http://localhost:80` e `observability.sh` contra as URLs internas.
- Entrega GitHub: [delivery.json](delivery.json) registra que não há `origin`; não houve publicação, commit ou push.

## Decisões técnicas a explicar

1. NGINX é a única fronteira pública e publica `80:80`; Go permanece em `8080` somente na rede interna.
2. `korp-project-network` é uma rede Docker externa, `bridge`, criada/validada pelo Ansible antes do Compose.
3. Prometheus coleta diretamente o Go na rede interna; Grafana consulta Prometheus pelo datasource `http://prometheus:9090`.
4. O horário é calculado dinamicamente a cada requisição e expresso em UTC RFC3339.
5. O playbook impõe versões aprovadas, usa preflight não destrutivo da porta 80 e valida o endpoint ao final.
6. Traefik permaneceu parado e preservado; a restauração é uma ação posterior, fora deste WP.
