# WP-07 baseline before provisioning

- Timestamp: 2026-08-07 (execution host timezone: America/Sao_Paulo)
- Architecture: `ARCH-KORP-001`, Revision `5`
- Architecture SHA-256: `72cd89047e7edd8aad5f6dc3572936daf1cf7deb977f59036e5901244abf6567`
- Ansible: `ansible-core 2.21.2` (`/home/magaldes/.local/bin/ansible-playbook`)
- Docker Engine: `29.6.0`
- Docker Compose: `5.2.0`
- Go tests: `cd app && go test ./...` passed
- Host TCP port 80: free (`ss -ltnp '( sport = :80 )'` returned no listener)
- Traefik: stopped (`docker ps -a --filter name=traefik` showed `Exited (0)`); no action taken by WP-07
- `korp-project-network`: absent before provisioning
- Existing Docker workloads: preserved; no Korp containers existed before provisioning

## Pre-existing working-tree status

```text
?? .gitignore
?? README.md
?? agent-packages/
?? app/
?? automation/
?? docs/
?? evidence/
?? infra/
?? tests/
```

The repository is initially untracked as a whole; this baseline is retained to distinguish WP-07 evidence from pre-existing files.
