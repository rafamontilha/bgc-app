.PHONY: help docker-up docker-down docker-restart docker-logs docker-ps docker-build docker-clean \
        k8s-setup k8s-up k8s-down k8s-restart k8s-logs k8s-status k8s-build k8s-open k8s-clean \
        seed restore-backup test-docker test-k8s

# Detectar sistema operacional
ifeq ($(OS),Windows_NT)
    SHELL := pwsh.exe
    .SHELLFLAGS := -NoProfile -Command
    PWSH := pwsh.exe
else
    PWSH := pwsh
endif

##@ Geral

help: ## Mostrar esta mensagem de ajuda
	@echo "BGC App - Sistema de Analytics de Exportação"
	@echo ""
	@echo "Uso: make <comando>"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"; printf "\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Docker Compose

docker-up: ## Iniciar serviços Docker Compose
	$(PWSH) -File ./scripts/docker.ps1 up

docker-down: ## Parar serviços Docker Compose
	$(PWSH) -File ./scripts/docker.ps1 down

docker-restart: ## Reiniciar serviços Docker Compose
	$(PWSH) -File ./scripts/docker.ps1 restart

docker-logs: ## Ver logs Docker Compose
	$(PWSH) -File ./scripts/docker.ps1 logs

docker-ps: ## Status dos containers Docker
	$(PWSH) -File ./scripts/docker.ps1 ps

docker-build: ## Rebuildar imagens Docker
	$(PWSH) -File ./scripts/docker.ps1 build

docker-clean: ## Limpar tudo Docker Compose (remove volumes)
	$(PWSH) -File ./scripts/docker.ps1 clean

##@ Kubernetes

k8s-setup: ## Setup inicial Kubernetes (cluster + deploy)
	$(PWSH) -File ./scripts/k8s.ps1 setup

k8s-up: ## Deploy serviços Kubernetes
	$(PWSH) -File ./scripts/k8s.ps1 up

k8s-down: ## Remover deployments Kubernetes
	$(PWSH) -File ./scripts/k8s.ps1 down

k8s-restart: ## Reiniciar pods Kubernetes
	$(PWSH) -File ./scripts/k8s.ps1 restart

k8s-logs: ## Ver logs Kubernetes
	$(PWSH) -File ./scripts/k8s.ps1 logs

k8s-status: ## Status do cluster Kubernetes
	$(PWSH) -File ./scripts/k8s.ps1 status

k8s-build: ## Rebuildar e reimportar imagens Kubernetes
	$(PWSH) -File ./scripts/k8s.ps1 build

k8s-open: ## Abrir aplicação no browser
	$(PWSH) -File ./scripts/k8s.ps1 open

k8s-clean: ## Deletar cluster Kubernetes
	$(PWSH) -File ./scripts/k8s.ps1 clean

##@ Dados

seed: ## Carregar dados de exemplo
	$(PWSH) -File ./scripts/seed.ps1

restore-backup: ## Restaurar backup do PostgreSQL (Kubernetes)
	$(PWSH) -File ./scripts/restore-backup.ps1

##@ Testes

test-docker: ## Testar ambiente Docker Compose
	@echo "Testando Docker Compose..."
	@curl -f http://localhost:8080/healthz || echo "API não está respondendo"
	@curl -f http://localhost:3000 || echo "Web não está respondendo"

test-k8s: ## Testar ambiente Kubernetes
	@echo "Testando Kubernetes..."
	@kubectl get pods -n data
	@kubectl get hpa -n data
	@curl -f http://web.bgc.local/healthz || echo "Ingress não está respondendo"

##@ Desenvolvimento

dev-api: ## Executar API localmente (requer Postgres)
	@cd api && go run cmd/api/main.go

dev-watch: ## Watch de mudanças (requer air)
	@cd api && air

lint: ## Executar linter Go
	@cd api && golangci-lint run

fmt: ## Formatar código Go
	@cd api && go fmt ./...

##@ Atalhos Rápidos

up: docker-up ## Alias para docker-up
down: docker-down ## Alias para docker-down
logs: docker-logs ## Alias para docker-logs
status: k8s-status ## Ver status geral

.DEFAULT_GOAL := help
