// bgc-api: expõe /health, /metrics/resumo e /metrics/pais lendo views em Postgres. // Comentário geral do serviço.

package main // Define o pacote principal (ponto de entrada do binário).

import ( // Bloco de imports das bibliotecas usadas.
	"context"        // Contextos para controlar tempo de execução/cancelamentos.
	"encoding/json"  // Serialização de respostas JSON.
	"fmt"            // Formatação de strings (para DSN, por exemplo).
	"log"            // Log simples para inicialização/erros do servidor.
	"net/http"       // Servidor HTTP padrão do Go.
	"os"             // Acesso a variáveis de ambiente (credenciais/host).
	"time"           // Timeouts e políticas do pool/servidor.

	"github.com/jackc/pgx/v5/pgxpool" // Driver nativo do Postgres com pool de conexões.
)

// -------- infra de conexão -------- // Delimita a seção de utilitários de conexão.

// getenv devolve valor da variável de ambiente ou um default se estiver vazia. // Documenta a função.
func getenv(k, def string) string { // Assinatura: recebe a chave e um default.
	if v := os.Getenv(k); v != "" { // Lê a env; se não vazio, usa o valor.
		return v // Retorna o valor definido na env.
	}
	return def // Senão, retorna o default informado.
}

// dsnFromEnv monta a string DSN a partir das variáveis de ambiente. // Descreve a função.
func dsnFromEnv() string { // Início da função.
	host := getenv("PGHOST", "pg-postgresql.data.svc.cluster.local") // Host do Postgres (service do cluster).
	port := getenv("PGPORT", "5432")                                 // Porta padrão do Postgres.
	user := getenv("PGUSER", "postgres")                             // Usuário do banco.
	pass := os.Getenv("PGPASSWORD")                                  // Senha (vem de Secret no K8s).
	db := getenv("PGDATABASE", "postgres")                           // Banco de dados alvo.
	ssl := getenv("PGSSLMODE", "disable")                            // Modo SSL (disable no cluster interno).
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, pass, host, port, db, ssl) // Concatena a DSN final.
}

// newPool cria um pool de conexões pgx com parâmetros básicos e timeout. // Explica a função.
func newPool() (*pgxpool.Pool, error) { // Retorna um *pgxpool.Pool ou erro.
	cfg, err := pgxpool.ParseConfig(dsnFromEnv()) // Converte a DSN em configuração do pool.
	if err != nil {                               // Se falhar o parse…
		return nil, err // Propaga o erro.
	}
	cfg.MaxConns = 8                      // Limita número máximo de conexões simultâneas.
	cfg.MinConns = 0                      // Não mantém conexões ociosas obrigatórias.
	cfg.MaxConnLifetime = 2 * time.Hour   // Recicla conexões a cada 2 horas.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // Contexto com timeout para criar o pool.
	defer cancel()                                                            // Garante cancelamento do contexto.
	return pgxpool.NewWithConfig(ctx, cfg) // Cria e retorna o pool com a config.
}

// -------- modelos de resposta -------- // Seção de tipos de resposta auxiliares.

type apiError struct {  // Estrutura JSON para erros de API.
	Error string `json:"error"` // Campo "error" com a mensagem.
}

// writeJSON centraliza cabeçalhos + status + encoding JSON. // Explica a função.
func writeJSON(w http.ResponseWriter, code int, payload any) { // Recebe writer, status HTTP e payload arbitrário.
	w.Header().Set("Content-Type", "application/json; charset=utf-8") // Define content-type como JSON.
	w.WriteHeader(code)                                               // Escreve o status HTTP.
	_ = json.NewEncoder(w).Encode(payload)                            // Serializa o payload em JSON (ignora erro de escrita).
}

// -------- handlers -------- // Seção com os manipuladores HTTP.

// healthHandler executa SELECT 1 para verificar conectividade. // Explica o handler.
func healthHandler(pool *pgxpool.Pool) http.HandlerFunc { // Recebe o pool e devolve um handler.
	return func(w http.ResponseWriter, r *http.Request) { // Função que atende a requisição.
		var one int                                         // Variável para receber o resultado.
		if err := pool.QueryRow(r.Context(), "SELECT 1").Scan(&one); err != nil { // Executa SELECT 1 e escaneia para 'one'.
			writeJSON(w, http.StatusInternalServerError, apiError{Error: err.Error()}) // Em erro, responde 500 com JSON de erro.
			return                                                                    // Sai do handler.
		}
		writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "select": one}) // Em sucesso, responde 200 com status/valor.
	}
}

// Resumo representa uma linha da view rpt.vw_exportacao_resumo. // Descreve o tipo.
type Resumo struct { // Estrutura do JSON retornado pelo endpoint /metrics/resumo.
	Ano        int     `json:"ano"`         // Ano do registro agregado.
	Setor      string  `json:"setor"`       // Setor agregado.
	ValorTotal float64 `json:"valor_total"` // Soma de valor.
	QtdeTotal  float64 `json:"qtde_total"`  // Soma de quantidade.
}

// resumoHandler consulta a view de resumo e devolve um array JSON. // Explica o handler.
func resumoHandler(pool *pgxpool.Pool) http.HandlerFunc { // Recebe o pool e devolve um handler.
	const q = `SELECT ano,setor,valor_total,qtde_total FROM rpt.vw_exportacao_resumo ORDER BY ano,setor` // Query SQL fixa.
	return func(w http.ResponseWriter, r *http.Request) { // Função do handler.
		rows, err := pool.Query(r.Context(), q) // Executa a consulta no contexto da request.
		if err != nil {                         // Se erro ao executar…
			writeJSON(w, http.StatusInternalServerError, apiError{Error: err.Error()}) // Responde 500 com erro.
			return                                                                     // Sai.
		}
		defer rows.Close() // Garante fechamento do cursor.

		var out []Resumo            // Slice para acumular as linhas.
		for rows.Next() {           // Itera as linhas da consulta.
			var x Resumo                                                  // Estrutura temporária para a linha.
			if err := rows.Scan(&x.Ano, &x.Setor, &x.ValorTotal, &x.QtdeTotal); err != nil { // Lê colunas para a struct.
				writeJSON(w, http.StatusInternalServerError, apiError{Error: err.Error()}) // Em erro de scan, 500.
				return                                                                     // Sai.
			}
			out = append(out, x) // Adiciona a linha ao slice de saída.
		}
		if err := rows.Err(); err != nil { // Checa erro de iteração (I/O, etc.).
			writeJSON(w, http.StatusInternalServerError, apiError{Error: err.Error()}) // Em erro, 500.
			return                                                                     // Sai.
		}
		writeJSON(w, http.StatusOK, out) // Sucesso: responde 200 com o slice JSON.
	}
}

// PorPais representa uma linha da view rpt.vw_exportacao_por_pais. // Descreve o tipo.
type PorPais struct { // Estrutura do JSON retornado pelo endpoint /metrics/pais.
	Ano        int     `json:"ano"`         // Ano do registro agregado.
	Pais       string  `json:"pais"`        // País agregado.
	ValorTotal float64 `json:"valor_total"` // Soma de valor.
	QtdeTotal  float64 `json:"qtde_total"`  // Soma de quantidade.
}

// porPaisHandler consulta a view por país e devolve um array JSON. // Explica o handler.
func porPaisHandler(pool *pgxpool.Pool) http.HandlerFunc { // Recebe o pool e devolve um handler.
	const q = `SELECT ano,pais,valor_total,qtde_total FROM rpt.vw_exportacao_por_pais ORDER BY ano,pais` // SQL fixa.
	return func(w http.ResponseWriter, r *http.Request) { // Função do handler.
		rows, err := pool.Query(r.Context(), q) // Executa a consulta.
		if err != nil {                         // Em erro…
			writeJSON(w, http.StatusInternalServerError, apiError{Error: err.Error()}) // 500 com erro.
			return                                                                     // Sai.
		}
		defer rows.Close() // Fecha o cursor ao final.

		var out []PorPais         // Slice de saída.
		for rows.Next() {         // Itera as linhas.
			var x PorPais                                                 // Struct temporária.
			if err := rows.Scan(&x.Ano, &x.Pais, &x.ValorTotal, &x.QtdeTotal); err != nil { // Faz o scan de colunas.
				writeJSON(w, http.StatusInternalServerError, apiError{Error: err.Error()}) // 500 se falhar.
				return                                                                     // Sai.
			}
			out = append(out, x) // Acumula a linha.
		}
		if err := rows.Err(); err != nil { // Checa erro de iteração.
			writeJSON(w, http.StatusInternalServerError, apiError{Error: err.Error()}) // 500 em erro.
			return                                                                     // Sai.
		}
		writeJSON(w, http.StatusOK, out) // Sucesso: responde com JSON.
	}
}

// -------- main / servidor HTTP -------- // Seção principal (bootstrap do servidor).

func main() { // Função main: inicia pool e servidor HTTP.
	pool, err := newPool() // Cria pool de conexões ao Postgres.
	if err != nil {        // Se falhar…
		log.Fatalf("connect: %v", err) // Aborta a aplicação com log do erro.
	}
	defer pool.Close() // Garante fechar o pool ao encerrar.

	mux := http.NewServeMux()                    // Cria um multiplexer HTTP simples.
	mux.HandleFunc("/health", healthHandler(pool))          // Registra handler de saúde.
	mux.HandleFunc("/metrics/resumo", resumoHandler(pool))  // Registra handler de resumo por setor.
	mux.HandleFunc("/metrics/pais", porPaisHandler(pool))   // Registra handler de resumo por país.

	srv := &http.Server{                    // Configura o servidor HTTP.
		Addr:              ":8080",         // Porta de escuta dentro do container.
		Handler:           mux,             // Multiplexer com as rotas registradas.
		ReadHeaderTimeout: 5 * time.Second, // Timeout para cabeçalhos (protege contra slowloris).
		ReadTimeout:       10 * time.Second,// Timeout total de leitura.
		WriteTimeout:      30 * time.Second,// Timeout de escrita de resposta.
		IdleTimeout:       60 * time.Second,// Timeout de conexões ociosas.
	}

	log.Println("bgc-api listening on :8080")          // Log informando que o servidor subiu.
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed { // Inicia servidor e checa erro não-graceful.
		log.Fatalf("server: %v", err) // Aborta em erro crítico do servidor.
	}
}
