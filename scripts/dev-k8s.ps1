param(
  [Parameter(Position=0)]
  [string]$Cmd = "help"
)

$ns = "data"
$cluster = "bgc"

function Usage {
@"
Comandos (K8s):
  dev-k8s.ps1 help              -> esta ajuda
  dev-k8s.ps1 info              -> cluster-info + pods/ingress no ns data
  dev-k8s.ps1 import-api        -> importa imagem local 'bgc-api:latest' no k3d: $cluster
  dev-k8s.ps1 apply             -> aplica manifests k8s (api/web)
  dev-k8s.ps1 cm-web            -> cria/atualiza ConfigMap bgc-web com web/index.html e web/routes.html
  dev-k8s.ps1 restart-api       -> rollout restart api
  dev-k8s.ps1 open              -> abre http://api.bgc.local/healthz e http://web.bgc.local
"@ | Write-Host
}

switch ($Cmd.ToLower()) {
  "help" { Usage; break }
  "info" {
    kubectl cluster-info
    kubectl -n $ns get pods -o wide
    kubectl -n $ns get svc
    kubectl -n $ns get ingress
    break
  }
  "import-api" {
    k3d image import bgc-api:latest -c $cluster
    break
  }
  "apply" {
    kubectl apply -f k8s/namespace.yaml
    kubectl apply -f k8s/api-deploy.yaml
    kubectl apply -f k8s/api-svc.yaml
    kubectl apply -f k8s/api-ingress.yaml
    kubectl -n $ns create configmap bgc-web `
      --from-file=web/index.html `
      --from-file=web/routes.html `
      --dry-run=client -o yaml | kubectl apply -f -
    kubectl apply -f k8s/web-deploy.yaml
    kubectl apply -f k8s/web-svc.yaml
    kubectl apply -f k8s/web-ingress.yaml
    break
  }
  "cm-web" {
    kubectl -n $ns create configmap bgc-web `
      --from-file=web/index.html `
      --from-file=web/routes.html `
      --dry-run=client -o yaml | kubectl apply -f -
    break
  }
  "restart-api" {
    kubectl -n $ns rollout restart deploy/bgc-api
    break
  }
  "open" {
    Start-Process "http://api.bgc.local/healthz"
    Start-Process "http://web.bgc.local"
    break
  }
  default { Usage; break }
}
