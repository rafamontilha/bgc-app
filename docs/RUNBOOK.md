# BGC Platform Runbook

**Last Updated:** 2026-01-09
**Version:** v0.4.0
**Status:** ✅ Production Ready (E2E Validated)

## Production Metrics (v0.4.0)

**Achieved Performance (2026-01-09):**
- ✅ API Response Time: **22-92ms** (target: <200ms) - **10x better!**
- ✅ Rate Limiting: Working correctly (5 req/day for free tier)
- ✅ Uptime: API, PostgreSQL, Redis, Integration Gateway all operational
- ✅ Test Coverage: 2 NCMs validated (17011400, 02013000)

**Infrastructure:**
- API: bgc-api:v0.4.0
- Database: 50 countries + 16 export records
- Cache: Redis L2 (512MB, allkeys-lru)
- Gateway: 2 replicas, 3 connectors

## Quick Reference

### Health Checks

```bash
# All pods status
kubectl get pods -n data

# API health
curl http://api.bgc.local/healthz

# Redis ping
kubectl exec -n data $(kubectl get pod -n data -l app=redis -o jsonpath='{.items[0].metadata.name}') -- redis-cli ping

# PostgreSQL connection
kubectl exec -n data $(kubectl get pod -n data -l app=postgres -o jsonpath='{.items[0].metadata.name}') -- psql -U bgc bgc -c "SELECT 1;"

# Integration Gateway health
kubectl run test-curl --image=curlimages/curl:latest --rm -i --restart=Never -n data -- curl -s http://integration-gateway:8081/health
```

## Common Issues

### Simulator Returns 0 Results

**Cause:** NCM not in seed data

**Solution:**
1. Check available NCMs: `kubectl exec -n data $(kubectl get pod -n data -l app=postgres -o jsonpath='{.items[0].metadata.name}') -- psql -U bgc bgc -c "SELECT DISTINCT co_ncm FROM stg.exportacao;"`
2. Use seed NCMs: 17011400, 26011200, 12010090
3. Or populate more data via ComexStat integration

### Rate Limit Hit Immediately

**Cause:** Shared IP (NAT, office)

**Solution:**
1. Authenticate user (get user_id for individual tracking)
2. Or upgrade to Premium tier (unlimited)

### Cache Not Working

**Cause:** Redis not running

**Solution:**
1. Check: `kubectl get pods -n data | grep redis`
2. Restart: `kubectl rollout restart deployment redis -n data`
3. Verify: Integration Gateway logs should show "Redis connected"

### PostgreSQL Restarts Frequently

**Cause:** Liveness probe too aggressive or memory pressure

**Solution Applied (2026-01-09):**
- Increased liveness probe timeout: 1s → 10s
- Increased initial delay: 30s → 60s
- Increased period: 10s → 30s
- Added memory limits: 2Gi
- Failure threshold: 3 → 5

**Monitor:** `kubectl get pod -n data $(kubectl get pod -n data -l app=postgres -o jsonpath='{.items[0].metadata.name}') -o jsonpath='{.status.containerStatuses[0].restartCount}'`

## Deployment

### Deploy New Version

```bash
# 1. Build new image
docker build -t bgc-api:v0.4.1 api/

# 2. Tag and import to k3d
k3d image import bgc-api:v0.4.1 -c bgc

# 3. Update deployment
kubectl set image deployment/bgc-api bgc-api=bgc-api:v0.4.1 -n data

# 4. Monitor rollout
kubectl rollout status deployment/bgc-api -n data

# 5. Verify
curl http://api.bgc.local/healthz
```

### Rollback

```bash
# Rollback to previous version
kubectl rollout undo deployment/bgc-api -n data

# Check rollout history
kubectl rollout history deployment/bgc-api -n data
```

## Monitoring

### Key Metrics

- **API Response Time P95:** < 200ms (target)
- **Error Rate:** < 0.1% (target)
- **Cache Hit Rate:** > 80% (target)
- **Database Connections:** < 50 (stable)
- **PostgreSQL Restarts:** 0 per week (healthy)

### Check Metrics

```bash
# Redis cache stats
kubectl exec -n data $(kubectl get pod -n data -l app=redis -o jsonpath='{.items[0].metadata.name}') -- redis-cli INFO stats | grep -E "(total_commands_processed|keyspace_hits|keyspace_misses)"

# PostgreSQL connections
kubectl exec -n data $(kubectl get pod -n data -l app=postgres -o jsonpath='{.items[0].metadata.name}') -- psql -U bgc bgc -c "SELECT count(*) FROM pg_stat_activity WHERE state = 'active';"

# Pod resource usage
kubectl top pods -n data
```

## Backup & Restore

### Backup PostgreSQL

```bash
# Create backup
kubectl exec -n data $(kubectl get pod -n data -l app=postgres -o jsonpath='{.items[0].metadata.name}') -- pg_dump -U bgc bgc > backup-$(date +%Y-%m-%d).sql

# Compress
gzip backup-$(date +%Y-%m-%d).sql

# Verify integrity
gzip -t backup-$(date +%Y-%m-%d).sql.gz
```

**Location:** Backups are stored in project root (add to .gitignore)

**Schedule:** Automated backups run via CronJob daily at 2 AM

### Restore from Backup

```bash
# Decompress
gunzip backup-2026-01-09.sql.gz

# Restore
kubectl exec -i -n data $(kubectl get pod -n data -l app=postgres -o jsonpath='{.items[0].metadata.name}') -- psql -U bgc bgc < backup-2026-01-09.sql
```

## Troubleshooting

### Slow API Responses

1. Check database: `kubectl top pod -n data $(kubectl get pod -n data -l app=postgres -o jsonpath='{.items[0].metadata.name}')`
2. Check API: `kubectl top pod -n data $(kubectl get pod -n data -l app=bgc-api -o jsonpath='{.items[0].metadata.name}')`
3. Check cache hit rate: Redis INFO stats
4. Analyze slow queries: Check PostgreSQL logs

### Memory Issues

1. Check resource usage: `kubectl top pods -n data`
2. Increase limits if needed: `kubectl edit deployment xxx -n data`
3. Check for memory leaks: Monitor over 24h

### Image Pull Errors

**Cause:** Image not imported to k3d cluster

**Solution:**
```bash
# Import image to k3d
k3d image import <image-name>:latest -c bgc

# Restart deployment
kubectl rollout restart deployment/<deployment-name> -n data
```

## Database Migrations

### Run Migrations

```bash
# Apply migration
kubectl exec -i -n data $(kubectl get pod -n data -l app=postgres -o jsonpath='{.items[0].metadata.name}') -- psql -U bgc bgc < db/migrations/XXXX_migration_name.sql

# Verify
kubectl exec -n data $(kubectl get pod -n data -l app=postgres -o jsonpath='{.items[0].metadata.name}') -- psql -U bgc bgc -c "\dt"
```

### Migration Status

```bash
# Check tables
kubectl exec -n data $(kubectl get pod -n data -l app=postgres -o jsonpath='{.items[0].metadata.name}') -- psql -U bgc bgc -c "\d+ countries_metadata"

# Check data count
kubectl exec -n data $(kubectl get pod -n data -l app=postgres -o jsonpath='{.items[0].metadata.name}') -- psql -U bgc bgc -c "SELECT COUNT(*) FROM countries_metadata;"
```

## Secrets Management

### View Secrets (Non-Sensitive)

```bash
# List secrets
kubectl get secrets -n data

# View secret keys (not values)
kubectl describe secret <secret-name> -n data
```

### Create/Update Secrets

```bash
# Create generic secret
kubectl create secret generic <name> \
  --from-literal=key=value \
  --namespace=data

# Update secret (delete and recreate)
kubectl delete secret <name> -n data
kubectl create secret generic <name> \
  --from-literal=key=new-value \
  --namespace=data
```

**Note:** For ComexStat credentials, use sealed secrets (see k8s/integration-gateway/README-SECRETS.md)

## Network Debugging

### Test Pod-to-Pod Communication

```bash
# Test Redis from API pod
kubectl exec -n data $(kubectl get pod -n data -l app=bgc-api -o jsonpath='{.items[0].metadata.name}') -- nc -zv redis.data.svc.cluster.local 6379

# Test PostgreSQL from API pod
kubectl exec -n data $(kubectl get pod -n data -l app=bgc-api -o jsonpath='{.items[0].metadata.name}') -- nc -zv pg-postgresql.data.svc.cluster.local 5432

# Test Integration Gateway from API pod
kubectl exec -n data $(kubectl get pod -n data -l app=bgc-api -o jsonpath='{.items[0].metadata.name}') -- nc -zv integration-gateway.data.svc.cluster.local 8081
```

### DNS Resolution

```bash
# Test DNS resolution
kubectl run test-dns --image=busybox:1.28 --rm -it --restart=Never -n data -- nslookup redis.data.svc.cluster.local
```

## Contacts & Escalation

| Area | Contact | Escalation |
|------|---------|------------|
| API Issues | Backend Team | Tech Lead |
| Database Issues | DevOps | DBA |
| Frontend Issues | Frontend Team | Product Manager |
| Infrastructure | DevOps | CTO |
| Integration Gateway | Integration Team | Tech Lead |

## Service URLs

### Development (k3d)

- **API:** http://api.bgc.local
- **Web Public:** http://web.bgc.local
- **Integration Gateway:** http://integration-gateway.data.svc.cluster.local:8081 (internal)

### Kubernetes Services

```bash
# List all services
kubectl get svc -n data

# Port-forward for local access
kubectl port-forward -n data svc/redis 6379:6379
kubectl port-forward -n data svc/pg-postgresql 5432:5432
kubectl port-forward -n data svc/integration-gateway 8081:8081
```

## Epic 4 Specific

### Simulator Endpoint

```bash
# Test simulator
curl -X POST http://api.bgc.local/v1/simulator/destinations \
  -H "Content-Type: application/json" \
  -d '{
    "ncm": "17011400",
    "volume_kg": 1000,
    "max_results": 10
  }'
```

### Check Countries Data

```bash
# Count countries
kubectl exec -n data $(kubectl get pod -n data -l app=postgres -o jsonpath='{.items[0].metadata.name}') -- psql -U bgc bgc -c "SELECT COUNT(*) FROM countries_metadata;"

# Top 5 closest countries
kubectl exec -n data $(kubectl get pod -n data -l app=postgres -o jsonpath='{.items[0].metadata.name}') -- psql -U bgc bgc -c "SELECT code, name_en, distance_brazil_km FROM countries_metadata ORDER BY distance_brazil_km LIMIT 5;"
```

### Re-run Populate Countries Job

```bash
# Delete old job
kubectl delete job -n data populate-countries

# Re-apply
kubectl apply -f k8s/jobs/populate-countries-job.yaml

# Monitor
kubectl logs -n data -l app=populate-countries -f
```

## Emergency Procedures

### Complete System Restart

```bash
# 1. Scale down all deployments
kubectl scale deployment --all --replicas=0 -n data

# 2. Wait for pods to terminate
kubectl get pods -n data -w

# 3. Scale up PostgreSQL first
kubectl scale deployment postgres --replicas=1 -n data
kubectl wait --for=condition=ready pod -l app=postgres -n data

# 4. Scale up Redis
kubectl scale deployment redis --replicas=1 -n data
kubectl wait --for=condition=ready pod -l app=redis -n data

# 5. Scale up Integration Gateway
kubectl scale deployment integration-gateway --replicas=2 -n data
kubectl wait --for=condition=ready pod -l app=integration-gateway -n data

# 6. Scale up API
kubectl scale deployment bgc-api --replicas=1 -n data
kubectl wait --for=condition=ready pod -l app=bgc-api -n data

# 7. Verify all services
kubectl get pods -n data
```

### Database Corruption Recovery

1. Stop all API pods: `kubectl scale deployment bgc-api --replicas=0 -n data`
2. Create backup: See "Backup PostgreSQL" section
3. Assess corruption: Check PostgreSQL logs
4. Restore from last good backup: See "Restore from Backup" section
5. Run migrations if needed
6. Restart API: `kubectl scale deployment bgc-api --replicas=1 -n data`

## Change Log

| Date | Change | Author |
|------|--------|--------|
| 2026-01-09 | Initial runbook created | DevOps Team |
| 2026-01-09 | Added PostgreSQL restart fix | DevOps Team |
| 2026-01-09 | Added Epic 4 specific sections | Product Team |
