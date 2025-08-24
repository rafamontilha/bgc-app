param(
  [Parameter(Position=0)]
  [ValidateSet("help","kubeinfo","psql-health")]
  [string]$cmd = "help"
)

switch ($cmd) {
  "help" {
    Write-Host "Comandos:" -ForegroundColor Cyan
    Write-Host "  dev.ps1 kubeinfo     -> info do cluster/nodes/pods"
    Write-Host "  dev.ps1 psql-health  -> SELECT version() via client dentro do cluster"
  }

  "kubeinfo" {
    kubectl cluster-info
    kubectl get nodes -o wide
    kubectl get pods -A
  }

  "psql-health" {
    $secretB64 = kubectl get secret -n data pg-postgresql -o jsonpath="{.data.postgres-password}"
    if (-not $secretB64) { throw "Secret 'pg-postgresql' não encontrado no ns 'data'." }
    $pwd = [Text.Encoding]::UTF8.GetString([Convert]::FromBase64String($secretB64))

    kubectl run pg-postgresql-client `
      --rm -it --restart='Never' -n data `
      --image docker.io/bitnami/postgresql:17.6.0-debian-12-r0 `
      --env="PGPASSWORD=$pwd" --env="PAGER=cat" `
      --command -- psql -h pg-postgresql -U postgres -d postgres -p 5432 `
      -c "SELECT version();"
  }
}
