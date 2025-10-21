# feat: Add observability, resilience and automation features

## 📊 Summary

This PR adds comprehensive observability, resilience, and automation features to improve production readiness and operational excellence.

## ✨ New Features

### Observability & Resilience
- ✅ **Health Probes**: Readiness and liveness probes for WEB deployment (API already had them)
- ✅ **HPA**: Horizontal Pod Autoscaling for API (1-5 pods) and WEB (1-3 pods)
- ✅ **Resource Limits**: CPU and memory requests/limits for all deployments
  - API: 100m-500m CPU, 128Mi-512Mi RAM
  - WEB: 50m-200m CPU, 64Mi-256Mi RAM

### Automation & Backup
- ✅ **PostgreSQL Backup CronJob**: Daily automated backups at 02:00
  - Compressed backups (.sql.gz)
  - Keeps last 7 backups
  - Stored in persistent PVC
- ✅ **Backup Restore Script**: PowerShell script for disaster recovery
- ✅ **MView Refresh CronJob**: Integrated into deployment workflow

### Developer Experience
- ✅ **Makefile**: Cross-platform wrapper for all PowerShell scripts
- ✅ **CHANGELOG.md**: Semantic versioning and release history
- ✅ **Enhanced Documentation**: Comprehensive README updates

## 🔧 Changes

### Kubernetes Manifests
- `k8s/api.yaml`: Added resource limits and requests
- `k8s/web.yaml`: Added probes and resource limits
- `k8s/api-hpa.yaml`: NEW - HPA configuration for API
- `k8s/web-hpa.yaml`: NEW - HPA configuration for WEB
- `k8s/postgres-backup-cronjob.yaml`: NEW - Backup automation

### Scripts
- `scripts/k8s.ps1`: Updated to deploy HPA and CronJobs automatically
- `scripts/restore-backup.ps1`: NEW - Backup restore utility

### Documentation
- `README.md`: Added observability section, updated features list
- `CHANGELOG.md`: NEW - Release history tracking
- `Makefile`: NEW - Cross-platform command wrapper

## 🧪 Testing

All environments tested and verified:

### Docker Compose
- ✅ API health check: OK
- ✅ WEB home page: 200 OK
- ✅ PostgreSQL connection: OK
- ✅ Data endpoints: OK

### Kubernetes
- ✅ Pods rollout: Successful
- ✅ HPA configuration: Working (CPU: 1-2%, Memory: 4-18%)
- ✅ Health probes: Active on all pods
- ✅ CronJobs: Created and scheduled
- ✅ Backup test: Manual job completed successfully
- ✅ API via Ingress: 200 OK
- ✅ WEB via Ingress: 200 OK

## 📈 Impact

### Before
- Manual scaling required
- No automated backups
- Limited observability
- No disaster recovery plan

### After
- Auto-scaling based on load
- Daily automated backups with 7-day retention
- Full health monitoring with auto-restart
- One-command backup restore
- Cross-platform development support

## 🎯 User Stories Completed

✅ **US-1 (Dados)**: Enhanced with automated backup and refresh
✅ **US-2 (API)**: Added resource limits and HPA
✅ **US-3 (K8s)**: Health probes and HPA implemented
✅ **US-4 (Ops)**: CronJobs, docs, and automation complete

## 📋 Checklist

- [x] Code follows project standards
- [x] All tests passing (Docker Compose + Kubernetes)
- [x] Documentation updated
- [x] CHANGELOG.md updated
- [x] No breaking changes
- [x] Backward compatible with existing deployments

## 🚀 Deployment Notes

Existing deployments will be automatically updated with:
- Health probes (5-10s longer startup time)
- Resource limits (ensures stable performance)
- HPA (automatic scaling based on load)
- Backup CronJob (runs at 02:00 daily)

No manual intervention required. Scripts remain compatible.

---

🤖 Generated with [Claude Code](https://claude.com/claude-code)
