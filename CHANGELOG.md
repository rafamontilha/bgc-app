# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- CHANGELOG.md para rastreamento de mudanças
- Health probes (readiness/liveness) no deployment WEB
- HorizontalPodAutoscaler (HPA) para API e WEB
- CronJob de backup automático do PostgreSQL
- Makefile como wrapper unificado dos scripts PowerShell
- Documentação de observabilidade e resiliência

### Changed
- Script k8s.ps1 atualizado para aplicar HPA e CronJobs
- README.md expandido com informações sobre HPA e backups

## [0.2.5.1] - 2025-01-15

### Changed
- Migração para PostgreSQL oficial (postgres:16) substituindo Bitnami
- Infraestrutura Kubernetes estabilizada
- Correção de secrets do banco de dados

## [0.2.5] - 2025-01-14

### Added
- API/Web estáveis em produção simulada
- Kubernetes deployments com Traefik Ingress
- Documentação completa de deployment
- Métricas de observabilidade

### Changed
- Sprint 2 finalizada com infraestrutura consolidada

## [0.1-sprint1] - 2025-01-10

### Added
- Infraestrutura inicial com k3d + PostgreSQL (Helm)
- Serviço de ingestão CSV/XLSX (bgc-ingest)
- Materialized Views para agregação de dados (rpt.*)
- API REST read-only com endpoints /metrics/*
- Manifests Kubernetes e scripts de automação
- Sistema de proveniência de dados (ingest_source, ingest_batch)
- Documentação de arquitetura e post-mortem Sprint 1

### Features
- Clean Architecture na API Go
- Endpoints: /market/size (TAM/SAM/SOM) e /routes/compare
- Docker Compose para desenvolvimento local
- Scripts PowerShell para gerenciamento (docker.ps1, k8s.ps1)

---

## Formato de Versionamento

- **MAJOR**: Mudanças incompatíveis na API
- **MINOR**: Novas funcionalidades mantendo compatibilidade
- **PATCH**: Correções de bugs e melhorias

## Tipos de Mudanças

- **Added**: Novas features
- **Changed**: Mudanças em funcionalidades existentes
- **Deprecated**: Features que serão removidas
- **Removed**: Features removidas
- **Fixed**: Correções de bugs
- **Security**: Correções de vulnerabilidades
