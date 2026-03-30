"""
TDD - Testes de validação de resource limits nos manifestos Kubernetes.

Motivação: todos os serviços usavam o mesmo template 128Mi/512Mi independente
da natureza do processo. Estes testes codificam as restrições corretas por
categoria de serviço para impedir regressão.

Categorias:
  - GO_SERVICE: Go/Gin binários eficientes (~20-40MB em runtime)
  - NEXTJS_SERVICE: Node.js + Next.js (startup ~200-250MB com MUI + Clerk)
  - NGINX_STATIC: Nginx servindo build estático (~3-15MB)
  - REDIS: Cache com maxmemory configurado explicitamente
  - GO_JOB: Jobs efêmeros em Go (sem carga contínua)

Execução:
    pip install pytest pyyaml
    pytest tests/k8s/ -v
"""

import pytest
import yaml
from pathlib import Path

ROOT = Path(__file__).parent.parent.parent
K8S = ROOT / "k8s"


# ---------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------

def parse_mi(value: str) -> int:
    """Converte strings de memória K8s (ex: '256Mi', '1Gi') para MiB (int)."""
    value = str(value)
    if value.endswith("Gi"):
        return int(value[:-2]) * 1024
    if value.endswith("Mi"):
        return int(value[:-2])
    if value.endswith("Ki"):
        return int(value[:-2]) // 1024
    return int(value) // (1024 * 1024)


def parse_m(value: str) -> int:
    """Converte strings de CPU K8s (ex: '100m', '1') para millicores (int)."""
    value = str(value)
    if value.endswith("m"):
        return int(value[:-1])
    return int(float(value) * 1000)


def load_yaml_all(path: Path) -> list[dict]:
    """Carrega todos os documentos YAML de um arquivo (multi-doc)."""
    with open(path) as f:
        return [doc for doc in yaml.safe_load_all(f) if doc]


def extract_resources(yaml_path: Path, deployment_name: str) -> dict:
    """Extrai o bloco resources do container principal de um Deployment."""
    docs = load_yaml_all(yaml_path)
    for doc in docs:
        if doc.get("kind") == "Deployment" and doc["metadata"]["name"] == deployment_name:
            containers = doc["spec"]["template"]["spec"]["containers"]
            return containers[0]["resources"]
    raise ValueError(f"Deployment '{deployment_name}' não encontrado em {yaml_path}")


def extract_hpa_min_replicas(yaml_path: Path, hpa_name: str) -> int:
    """Extrai minReplicas de um HPA."""
    docs = load_yaml_all(yaml_path)
    for doc in docs:
        if doc.get("kind") == "HorizontalPodAutoscaler" and doc["metadata"]["name"] == hpa_name:
            return doc["spec"]["minReplicas"]
    raise ValueError(f"HPA '{hpa_name}' não encontrado em {yaml_path}")


# ---------------------------------------------------------------------------
# Fixtures
# ---------------------------------------------------------------------------

@pytest.fixture
def api_resources():
    return extract_resources(K8S / "api.yaml", "bgc-api")


@pytest.fixture
def web_next_resources():
    return extract_resources(K8S / "web-next.yaml", "bgc-web-next")


@pytest.fixture
def web_public_resources():
    return extract_resources(K8S / "web-public.yaml", "bgc-web-public")


@pytest.fixture
def redis_resources():
    return extract_resources(K8S / "redis.yaml", "redis")


@pytest.fixture
def gateway_resources():
    return extract_resources(K8S / "integration-gateway" / "deployment.yaml", "integration-gateway")


@pytest.fixture
def populate_job_resources():
    docs = load_yaml_all(K8S / "jobs" / "populate-countries-job.yaml")
    for doc in docs:
        if doc.get("kind") == "Job":
            containers = doc["spec"]["template"]["spec"]["containers"]
            return containers[0]["resources"]
    raise ValueError("Job populate-countries não encontrado")


@pytest.fixture
def redis_maxmemory_mb() -> int:
    """Extrai o valor de maxmemory do ConfigMap do Redis (em MB)."""
    docs = load_yaml_all(K8S / "redis.yaml")
    for doc in docs:
        if doc.get("kind") == "ConfigMap" and doc["metadata"]["name"] == "redis-config":
            conf = doc["data"]["redis.conf"]
            for line in conf.splitlines():
                if line.strip().startswith("maxmemory "):
                    val = line.strip().split()[1]
                    if val.endswith("mb"):
                        return int(val[:-2])
                    if val.endswith("gb"):
                        return int(val[:-2]) * 1024
    raise ValueError("maxmemory não encontrado no redis-config ConfigMap")


# ---------------------------------------------------------------------------
# Invariante global: limit >= request para qualquer serviço
# ---------------------------------------------------------------------------

class TestLimitGeqRequest:
    """Limite deve sempre ser >= request. Violação causa OOMKill imediato."""

    def test_api_memory_limit_gte_request(self, api_resources):
        req = parse_mi(api_resources["requests"]["memory"])
        lim = parse_mi(api_resources["limits"]["memory"])
        assert lim >= req, f"API: limit {lim}Mi < request {req}Mi"

    def test_web_next_memory_limit_gte_request(self, web_next_resources):
        req = parse_mi(web_next_resources["requests"]["memory"])
        lim = parse_mi(web_next_resources["limits"]["memory"])
        assert lim >= req, f"web-next: limit {lim}Mi < request {req}Mi"

    def test_web_public_memory_limit_gte_request(self, web_public_resources):
        req = parse_mi(web_public_resources["requests"]["memory"])
        lim = parse_mi(web_public_resources["limits"]["memory"])
        assert lim >= req, f"web-public: limit {lim}Mi < request {req}Mi"

    def test_redis_memory_limit_gte_request(self, redis_resources):
        req = parse_mi(redis_resources["requests"]["memory"])
        lim = parse_mi(redis_resources["limits"]["memory"])
        assert lim >= req, f"Redis: limit {lim}Mi < request {req}Mi"

    def test_gateway_memory_limit_gte_request(self, gateway_resources):
        req = parse_mi(gateway_resources["requests"]["memory"])
        lim = parse_mi(gateway_resources["limits"]["memory"])
        assert lim >= req, f"Gateway: limit {lim}Mi < request {req}Mi"


# ---------------------------------------------------------------------------
# Serviços Go: request lean, sem sobreprovisão
# ---------------------------------------------------------------------------

class TestGoServiceMemory:
    """
    Go compila para binário estático eficiente.
    Requests acima de 100Mi para serviços Go simples é sobreprovisão pura.
    """

    MAX_GO_REQUEST_MI = 100   # teto para requests de serviços Go
    MAX_GO_LIMIT_MI = 384     # limite razoável com headroom generoso

    def test_api_memory_request_is_lean(self, api_resources):
        req = parse_mi(api_resources["requests"]["memory"])
        assert req <= self.MAX_GO_REQUEST_MI, (
            f"bgc-api memory request {req}Mi está acima de {self.MAX_GO_REQUEST_MI}Mi. "
            "Go/Gin em produção usa ~20-40MB; reduza o request."
        )

    def test_api_memory_limit_not_excessive(self, api_resources):
        lim = parse_mi(api_resources["limits"]["memory"])
        assert lim <= self.MAX_GO_LIMIT_MI, (
            f"bgc-api memory limit {lim}Mi está acima de {self.MAX_GO_LIMIT_MI}Mi. "
            "Sem operações in-memory de larga escala, este limite é sobreprovisão."
        )

    def test_gateway_memory_request_is_lean(self, gateway_resources):
        req = parse_mi(gateway_resources["requests"]["memory"])
        assert req <= self.MAX_GO_REQUEST_MI, (
            f"integration-gateway memory request {req}Mi > {self.MAX_GO_REQUEST_MI}Mi. "
            "Go + Ristretto L1 cache usa ~60-80MB em operação."
        )

    def test_gateway_memory_limit_not_excessive(self, gateway_resources):
        lim = parse_mi(gateway_resources["limits"]["memory"])
        assert lim <= self.MAX_GO_LIMIT_MI, (
            f"integration-gateway memory limit {lim}Mi > {self.MAX_GO_LIMIT_MI}Mi."
        )

    def test_populate_job_memory_request_is_lean(self, populate_job_resources):
        req = parse_mi(populate_job_resources["requests"]["memory"])
        assert req <= self.MAX_GO_REQUEST_MI, (
            f"populate-countries job memory request {req}Mi > {self.MAX_GO_REQUEST_MI}Mi. "
            "Job Go efêmero de seed de dados não precisa de mais que 100Mi de request."
        )


# ---------------------------------------------------------------------------
# Next.js: request MÍNIMO para evitar OOMKill no startup
# ---------------------------------------------------------------------------

class TestNextJsMemory:
    """
    Next.js 15 + React 19 + MUI v7 + Clerk consome ~200-250MB no startup.
    Request abaixo de 200Mi causa OOMKill antes do pod ficar Ready.
    """

    MIN_NEXTJS_REQUEST_MI = 200   # mínimo para startup seguro
    MAX_NEXTJS_LIMIT_MI = 768     # teto razoável para 10 páginas

    def test_web_next_memory_request_adequate_for_startup(self, web_next_resources):
        req = parse_mi(web_next_resources["requests"]["memory"])
        assert req >= self.MIN_NEXTJS_REQUEST_MI, (
            f"bgc-web-next memory request {req}Mi é insuficiente para startup do Next.js. "
            f"Node.js + Next.js + MUI + Clerk consome ~200-250MB. "
            f"Mínimo recomendado: {self.MIN_NEXTJS_REQUEST_MI}Mi."
        )

    def test_web_next_memory_limit_reasonable_for_current_scope(self, web_next_resources):
        lim = parse_mi(web_next_resources["limits"]["memory"])
        assert lim <= self.MAX_NEXTJS_LIMIT_MI, (
            f"bgc-web-next memory limit {lim}Mi > {self.MAX_NEXTJS_LIMIT_MI}Mi. "
            "Para 10 páginas no escopo atual, este limite é excessivo."
        )


# ---------------------------------------------------------------------------
# Nginx estático: footprint mínimo
# ---------------------------------------------------------------------------

class TestNginxStaticMemory:
    """
    Nginx servindo arquivos estáticos (build Vite) usa 3-15MB.
    Requests de 128Mi+ são 10x o necessário.
    """

    MAX_NGINX_REQUEST_MI = 32   # nginx estático não precisa de mais
    MAX_NGINX_LIMIT_MI = 96     # headroom generoso para eventual nginx + lua

    def test_web_public_memory_request_minimal(self, web_public_resources):
        req = parse_mi(web_public_resources["requests"]["memory"])
        assert req <= self.MAX_NGINX_REQUEST_MI, (
            f"bgc-web-public memory request {req}Mi é excessivo para nginx estático. "
            f"Nginx em produção usa 3-15MB. Máximo razoável: {self.MAX_NGINX_REQUEST_MI}Mi."
        )

    def test_web_public_memory_limit_minimal(self, web_public_resources):
        lim = parse_mi(web_public_resources["limits"]["memory"])
        assert lim <= self.MAX_NGINX_LIMIT_MI, (
            f"bgc-web-public memory limit {lim}Mi > {self.MAX_NGINX_LIMIT_MI}Mi. "
            "Nginx estático não justifica mais do que 96Mi de limite."
        )

    def test_web_public_cpu_request_minimal(self, web_public_resources):
        req = parse_m(web_public_resources["requests"]["cpu"])
        assert req <= 20, (
            f"bgc-web-public CPU request {req}m é excessivo para nginx estático. "
            "Nginx idle usa <5m. Máximo razoável: 20m."
        )


# ---------------------------------------------------------------------------
# Redis: limit deve respeitar headroom sobre maxmemory
# ---------------------------------------------------------------------------

class TestRedisMemory:
    """
    Redis limit deve ser > maxmemory para acomodar buffers de replicação,
    AOF rewrite buffer e overhead do processo.
    Regra: limit >= maxmemory + 200Mi de headroom.
    """

    MIN_HEADROOM_MI = 200

    def test_redis_limit_has_headroom_over_maxmemory(self, redis_resources, redis_maxmemory_mb):
        lim = parse_mi(redis_resources["limits"]["memory"])
        required = redis_maxmemory_mb + self.MIN_HEADROOM_MI
        assert lim >= required, (
            f"Redis memory limit {lim}Mi não tem headroom suficiente sobre "
            f"maxmemory={redis_maxmemory_mb}mb. "
            f"Mínimo recomendado: {required}Mi (maxmemory + {self.MIN_HEADROOM_MI}Mi)."
        )

    def test_redis_request_below_maxmemory(self, redis_resources, redis_maxmemory_mb):
        req = parse_mi(redis_resources["requests"]["memory"])
        assert req < redis_maxmemory_mb, (
            f"Redis memory request {req}Mi >= maxmemory {redis_maxmemory_mb}Mi. "
            "O scheduler reservaria mais memória do que o Redis usaria em idle."
        )


# ---------------------------------------------------------------------------
# HPA: minReplicas para ambiente de desenvolvimento/MVP
# ---------------------------------------------------------------------------

class TestHpaMinReplicas:
    """
    Para MVP e ambiente de desenvolvimento, minReplicas=1 economiza recursos.
    O HPA escala automaticamente conforme a carga aumenta.
    """

    def test_web_next_hpa_min_replicas_is_one(self):
        min_r = extract_hpa_min_replicas(K8S / "web-next.yaml", "bgc-web-next-hpa")
        assert min_r == 1, (
            f"bgc-web-next-hpa minReplicas={min_r}. "
            "Para MVP com HPA configurado, 1 réplica mínima é suficiente."
        )

    def test_web_public_hpa_min_replicas_is_one(self):
        min_r = extract_hpa_min_replicas(K8S / "web-public.yaml", "bgc-web-public-hpa")
        assert min_r == 1, (
            f"bgc-web-public-hpa minReplicas={min_r}. "
            "Nginx estático aguenta 1 réplica; HPA escala com tráfego."
        )

    def test_gateway_hpa_min_replicas_is_one(self):
        min_r = extract_hpa_min_replicas(
            K8S / "integration-gateway" / "deployment.yaml",
            "integration-gateway-hpa",
        )
        assert min_r == 1, (
            f"integration-gateway-hpa minReplicas={min_r}. "
            "Em MVP com poucos conectores, 1 réplica mínima é adequada."
        )


# ---------------------------------------------------------------------------
# Proporção request/limit: evita starvation e sobreprovisão simultâneos
# ---------------------------------------------------------------------------

class TestRequestLimitRatio:
    """
    Ratio limit/request muito alto indica request artificialmente baixo
    (scheduler reserva pouco mas o pod pode usar muito mais, causando noisy-neighbor).
    Ratio saudável: entre 2x e 8x para maioria dos workloads.
    """

    MAX_RATIO = 8.0

    @pytest.mark.parametrize("fixture_name,service_name", [
        ("api_resources", "bgc-api"),
        ("web_next_resources", "bgc-web-next"),
        ("web_public_resources", "bgc-web-public"),
        ("gateway_resources", "integration-gateway"),
        ("redis_resources", "redis"),
    ])
    def test_memory_request_limit_ratio(self, fixture_name, service_name, request):
        resources = request.getfixturevalue(fixture_name)
        req = parse_mi(resources["requests"]["memory"])
        lim = parse_mi(resources["limits"]["memory"])
        ratio = lim / req
        assert ratio <= self.MAX_RATIO, (
            f"{service_name}: ratio limit/request = {ratio:.1f}x (limit={lim}Mi, request={req}Mi). "
            f"Ratio acima de {self.MAX_RATIO}x indica request subprovisado — "
            "o scheduler reserva pouco mas o pod pode consumir muito mais."
        )
