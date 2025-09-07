// main.go — bgc-ingest
// Este binário implementa 3 comandos:
//   1) `health`         -> testa conexão no Postgres com SELECT 1
//   2) `insert-sample`  -> insere uma linha de exemplo em stg.exportacao
//   3) `load-csv`       -> carrega um CSV para stg.exportacao (opções --path/--sep/--dec/--header)
//
// Observação: usamos pgx/pool (driver nativo) e variáveis de ambiente
// (PGHOST, PGPORT, PGUSER, PGPASSWORD, PGDATABASE, PGSSLMODE) para configurar a conexão.

package main

import (
	"bufio"          // buffer de leitura de arquivo (eficiente p/ CSVs grandes)
	"context"        // contextos com timeout/cancel para operações no banco
	"encoding/csv"   // parser de CSV
	"errors"         // tratamento de erros
	"flag"           // parsing de flags do comando load-csv
	"fmt"            // impressão/format
	"io"             // io.EOF para detectar final do arquivo
	"os"             // acesso a arquivos e variáveis de ambiente
	"strconv"        // conversão string->int/float
	"strings"        // utilidades de string (trim, replace, etc.)
	excelize "github.com/xuri/excelize/v2"
	"time"
	"github.com/jackc/pgx/v5/pgxpool" // pool de conexões com Postgres
)

////////////////////////////////////////////////////////////////////////////////
// Utilidades de configuração/conexão
////////////////////////////////////////////////////////////////////////////////

// getenv devolve o valor da variável de ambiente k, ou def se estiver vazia.
// Mantemos defaults sensatos para rodar dentro do cluster.
func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

// dsnFromEnv monta a string de conexão (DSN) a partir das envs esperadas.
// Ex.: postgres://user:pass@host:port/db?sslmode=disable
func dsnFromEnv() string {
	host := getenv("PGHOST", "pg-postgresql.data.svc.cluster.local") // service do chart Bitnami
	port := getenv("PGPORT", "5432")                                 // porta padrão Postgres
	user := getenv("PGUSER", "postgres")                             // usuário padrão
	pass := os.Getenv("PGPASSWORD")                                  // senha (vem do Secret do K8s)
	db := getenv("PGDATABASE", "postgres")                           // banco padrão
	ssl := getenv("PGSSLMODE", "disable")                            // em dev, costumamos desabilitar
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, pass, host, port, db, ssl)
}

// connect cria um *pgxpool.Pool com timeout e parâmetros básicos.
// O pool gerencia conexões reutilizáveis para eficiência.
func connect() (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dsnFromEnv()) // parse da DSN -> config
	if err != nil {
		return nil, err
	}
	// Parâmetros razoáveis para dev; em prod, ajuste conforme carga:
	cfg.MaxConns = 4                  // limita conexões simultâneas do pool
	cfg.MinConns = 0                  // não mantém conexões ociosas
	cfg.MaxConnLifetime = time.Hour   // recicla conexões a cada 1h

	// Criamos o pool com um contexto de criação com timeout (10s):
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return pgxpool.NewWithConfig(ctx, cfg)
}

////////////////////////////////////////////////////////////////////////////////
// Comando: health
////////////////////////////////////////////////////////////////////////////////

// cmdHealth conecta no banco e executa SELECT 1 para checar saúde.
func cmdHealth() error {
	pool, err := connect()
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	defer pool.Close() // devolve conexões do pool

	var one int
	// QueryRow -> uma linha; Scan lê a primeira coluna para a variável one
	if err := pool.QueryRow(context.Background(), "SELECT 1").Scan(&one); err != nil {
		return fmt.Errorf("query: %w", err)
	}
	// imprime JSON simples (fácil de ler no kubectl logs)
	fmt.Println(`{"status":"ok","select":`, one, `}`)
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// Comando: insert-sample
////////////////////////////////////////////////////////////////////////////////

// cmdInsertSample insere uma linha de exemplo em stg.exportacao, retornando o id.
func cmdInsertSample() error {
	pool, err := connect()
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	defer pool.Close()

	// SQL com RETURNING id para capturar o identificador gerado
	sql := `
		INSERT INTO stg.exportacao (ano, setor, pais, ncm, valor, qtde)
		VALUES ($1,$2,$3,$4,$5,$6)
		RETURNING id;
	`

	var id int64
	// QueryRow + Scan executa o insert e lê a coluna "id" retornada
	err = pool.QueryRow(context.Background(), sql,
		2024, "Teste", "Chile", "0101.10.10", 100.00, 1.00,
	).Scan(&id)
	if err != nil {
		return fmt.Errorf("insert: %w", err)
	}
	fmt.Println(`{"inserted_id":`, id, `}`)
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// Comando: load-csv
////////////////////////////////////////////////////////////////////////////////

// csvOpts define as opções aceitas por load-csv (via flags).
type csvOpts struct {
	path   string // caminho do arquivo dentro do pod (ex.: /data/sample.csv)
	sep    rune   // separador de campos (',' ou ';', etc.)
	dec    string // separador decimal ('.' ou ',')
	hasHdr bool   // primeira linha é cabeçalho?
	table  string // tabela de destino (mantido p/ futuro; hoje usamos stg.exportacao)
}

// nopWriter "silencia" a saída padrão do flag.FlagSet (evita poluir logs)
type nopWriter struct{}

func (n *nopWriter) Write(p []byte) (int, error) { return len(p), nil }

// parseCSVOpts lê as flags do comando load-csv e valida entradas.
func parseCSVOpts(args []string) (*csvOpts, error) {
	fs := flag.NewFlagSet("load-csv", flag.ContinueOnError)

	// Define flags com defaults:
	path := fs.String("path", "", "caminho do CSV dentro do pod (ex.: /data/sample.csv)")
	sepStr := fs.String("sep", ",", "separador de campos (ex.: ',' ou ';')")
	dec := fs.String("dec", ".", "separador decimal ('.' ou ',')")
	hasHdr := fs.Bool("header", true, "primeira linha é cabeçalho?")
	table := fs.String("table", "stg.exportacao", "tabela de destino")

	// Evita que o FlagSet escreva help no stdout (deixa erros limpos nos logs)
	fs.SetOutput(new(nopWriter))

	// Faz o parse do array de argumentos (após "load-csv")
	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	// Valida obrigatórios:
	if *path == "" {
		return nil, errors.New("obrigatório: --path")
	}

	// Converte o separador de string para rune (precisa 1 caractere)
	runes := []rune(*sepStr)
	if len(runes) != 1 {
		return nil, errors.New("--sep deve ter 1 caractere")
	}

	return &csvOpts{
		path:   *path,
		sep:    runes[0],
		dec:    *dec,
		hasHdr: *hasHdr,
		table:  *table,
	}, nil
}

// cmdLoadCSV abre o arquivo CSV, interpreta colunas e insere na stg.exportacao.
func cmdLoadCSV(args []string) error {
	// Lê e valida opções
	opts, err := parseCSVOpts(args)
	if err != nil {
		return err
	}

	// Abre o arquivo para leitura
	f, err := os.Open(opts.path)
	if err != nil {
		return fmt.Errorf("abrir csv: %w", err)
	}
	defer f.Close()

	// csv.Reader com separador configurável; FieldsPerRecord=-1 permite linhas com
	// contagem de campos variável (útil para CSVs “imperfeitos”)
	r := csv.NewReader(bufio.NewReader(f))
	r.Comma = opts.sep
	r.FieldsPerRecord = -1
	r.ReuseRecord = true // reaproveita slice internamente (menos GC)

	// Índices das colunas; tentaremos mapear automaticamente via cabeçalho
	var idxAno, idxSetor, idxPais, idxNcm, idxValor, idxQtde int

	// autoIndex tenta localizar colunas por nome (case-insensitive)
	autoIndex := func(header []string) {
		find := func(name string) int {
			for i, h := range header {
				if strings.EqualFold(strings.TrimSpace(h), name) {
					return i
				}
			}
			return -1
		}
		idxAno = find("ano")
		idxSetor = find("setor")
		idxPais = find("pais")
		idxNcm = find("ncm")
		idxValor = find("valor")
		idxQtde = find("qtde")
	}

	// Se o CSV tem cabeçalho, lê a 1ª linha e tenta mapear colunas por nome.
	if opts.hasHdr {
		hdr, err := r.Read()
		if err != nil {
			return fmt.Errorf("ler cabeçalho: %w", err)
		}
		autoIndex(hdr)
	}

	// Se alguma coluna não foi encontrada por nome, assume ordem padrão (0..5).
	defaultIfNeg := func(i, def int) int {
		if i < 0 {
			return def
		}
		return i
	}
	idxAno = defaultIfNeg(idxAno, 0)
	idxSetor = defaultIfNeg(idxSetor, 1)
	idxPais = defaultIfNeg(idxPais, 2)
	idxNcm = defaultIfNeg(idxNcm, 3)
	idxValor = defaultIfNeg(idxValor, 4)
	idxQtde = defaultIfNeg(idxQtde, 5)

	// Abre pool e inicia transação (melhor para múltiplos inserts)
	pool, err := connect()
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	defer pool.Close()

	tx, err := pool.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("begin: %w", err)
	}
	// Garante rollback se der erro ou se sair antes do commit
	defer func() { _ = tx.Rollback(context.Background()) }()

	// Statement simples; depois podemos evoluir para COPY/COPY FROM p/ performance
	stmt := `INSERT INTO stg.exportacao (ano, setor, pais, ncm, valor, qtde) VALUES ($1,$2,$3,$4,$5,$6)`

	// decToDot converte decimal com vírgula para ponto, se necessário
	decToDot := func(s string) string {
		if opts.dec == "," {
			return strings.ReplaceAll(s, ",", ".")
		}
		return s
	}

	count := 0 // contador de linhas carregadas
	for {
		// Lê próxima linha do CSV
		rec, err := r.Read()

		// io.EOF -> acabou o arquivo (sai do loop)
		if err == io.EOF {
			break
		}
		// Se houve erro e NÃO é diferença de quantidade de campos, aborta
		if err != nil && !errors.Is(err, csv.ErrFieldCount) {
			return fmt.Errorf("csv read: %w", err)
		}
		// Linhas vazias: segue para a próxima
		if rec == nil || len(rec) == 0 {
			continue
		}

		// Parse e normalização dos campos:
		ano, _ := strconv.Atoi(strings.TrimSpace(rec[idxAno]))
		setor := strings.TrimSpace(rec[idxSetor])
		pais := strings.TrimSpace(rec[idxPais])
		ncm := strings.TrimSpace(rec[idxNcm])
		valor, _ := strconv.ParseFloat(decToDot(strings.TrimSpace(rec[idxValor])), 64)
		qtde, _ := strconv.ParseFloat(decToDot(strings.TrimSpace(rec[idxQtde])), 64)

		// Executa INSERT na transação
		if _, err := tx.Exec(context.Background(), stmt, ano, setor, pais, ncm, valor, qtde); err != nil {
			return fmt.Errorf("insert linha %d: %w", count+1, err)
		}
		count++
	}

	// Fecha a transação gravando as mudanças
	if err := tx.Commit(context.Background()); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	// Log amigável em JSON com o total carregado
	fmt.Printf(`{"loaded": %d, "table":"%s"}`+"\n", count, opts.table)
	return nil
}
// ===== XLSX loader =====

type xlsxOpts struct {
    path   string // /data/sample.xlsx
    sheet  string // nome da planilha (ex.: "Sheet1"); se vazio, usa a primeira
    dec    string // separador decimal: "." ou ","
    hasHdr bool   // primeira linha é cabeçalho?
    table  string // tabela destino (mantém para futura flexibilidade)
}

func parseXLSXOpts(args []string) (*xlsxOpts, error) {
    fs := flag.NewFlagSet("load-xlsx", flag.ContinueOnError)
    path := fs.String("path", "", "caminho do XLSX dentro do pod (ex.: /data/sample.xlsx)")
    sheet := fs.String("sheet", "", "nome da planilha (se vazio, usa a primeira)")
    dec := fs.String("dec", ".", "separador decimal ('.' ou ',')")
    hasHdr := fs.Bool("header", true, "primeira linha é cabeçalho?")
    table := fs.String("table", "stg.exportacao", "tabela de destino")
    fs.SetOutput(new(nopWriter))
    if err := fs.Parse(args); err != nil {
        return nil, err
    }
    if *path == "" {
        return nil, errors.New("obrigatório: --path")
    }
    return &xlsxOpts{
        path:   *path,
        sheet:  *sheet,
        dec:    *dec,
        hasHdr: *hasHdr,
        table:  *table,
    }, nil
}

func cmdLoadXLSX(args []string) error {
    opts, err := parseXLSXOpts(args)
    if err != nil { return err }

    f, err := excelize.OpenFile(opts.path)
    if err != nil { return fmt.Errorf("abrir xlsx: %w", err) }
    defer f.Close()

    sheet := opts.sheet
    if sheet == "" {
        sheets := f.GetSheetList()
        if len(sheets) == 0 {
            return errors.New("arquivo xlsx sem planilhas")
        }
        sheet = sheets[0]
    }

    rows, err := f.Rows(sheet)
    if err != nil { return fmt.Errorf("abrir linhas da sheet %q: %w", sheet, err) }
    defer rows.Close()

    // índices das colunas (tentamos mapear por cabeçalho)
    var idxAno, idxSetor, idxPais, idxNcm, idxValor, idxQtde int
    autoIndex := func(header []string) {
        find := func(name string) int {
            for i, h := range header {
                if strings.EqualFold(strings.TrimSpace(h), name) { return i }
            }
            return -1
        }
        idxAno   = find("ano")
        idxSetor = find("setor")
        idxPais  = find("pais")
        idxNcm   = find("ncm")
        idxValor = find("valor")
        idxQtde  = find("qtde")
    }

    // lê cabeçalho, se houver
    if opts.hasHdr && rows.Next() {
        hdr, err := rows.Columns()
        if err != nil { return fmt.Errorf("ler cabeçalho: %w", err) }
        autoIndex(hdr)
    }

    defaultIfNeg := func(i, def int) int {
        if i < 0 { return def }
        return i
    }
    idxAno = defaultIfNeg(idxAno, 0)
    idxSetor = defaultIfNeg(idxSetor, 1)
    idxPais = defaultIfNeg(idxPais, 2)
    idxNcm = defaultIfNeg(idxNcm, 3)
    idxValor = defaultIfNeg(idxValor, 4)
    idxQtde = defaultIfNeg(idxQtde, 5)

    decToDot := func(s string) string {
        if opts.dec == "," { return strings.ReplaceAll(s, ",", ".") }
        return s
    }

    pool, err := connect()
    if err != nil { return fmt.Errorf("connect: %w", err) }
    defer pool.Close()

    tx, err := pool.Begin(context.Background())
    if err != nil { return fmt.Errorf("begin: %w", err) }
    defer func() { _ = tx.Rollback(context.Background()) }()

    stmt := `INSERT INTO stg.exportacao (ano, setor, pais, ncm, valor, qtde) VALUES ($1,$2,$3,$4,$5,$6)`

    count := 0
    for rows.Next() {
        rec, err := rows.Columns()
        if err != nil { return fmt.Errorf("ler linha: %w", err) }
        // pular linhas totalmente vazias
        empty := true
        for _, c := range rec { if strings.TrimSpace(c) != "" { empty = false; break } }
        if empty { continue }

        get := func(idx int) string {
            if idx >= 0 && idx < len(rec) { return strings.TrimSpace(rec[idx]) }
            return ""
        }
        anoStr   := get(idxAno)
        setor    := get(idxSetor)
        pais     := get(idxPais)
        ncm      := get(idxNcm)
        valorStr := get(idxValor)
        qtdeStr  := get(idxQtde)

        ano, _ := strconv.Atoi(anoStr)
        valor, _ := strconv.ParseFloat(decToDot(valorStr), 64)
        qtde, _ := strconv.ParseFloat(decToDot(qtdeStr), 64)

        if _, err := tx.Exec(context.Background(), stmt, ano, setor, pais, ncm, valor, qtde); err != nil {
            return fmt.Errorf("insert linha %d: %w", count+1, err)
        }
        count++
    }
    if err := rows.Error(); err != nil { return fmt.Errorf("iterar linhas: %w", err) }

    if err := tx.Commit(context.Background()); err != nil {
        return fmt.Errorf("commit: %w", err)
    }
    fmt.Printf(`{"loaded": %d, "table":"%s","source":"xlsx","sheet":"%s"}`+"\n", count, opts.table, sheet)
    return nil
}

////////////////////////////////////////////////////////////////////////////////
// Função principal (roteia para o subcomando)
////////////////////////////////////////////////////////////////////////////////

func main() {
	// Verifica se pelo menos 1 argumento foi passado (o nome do comando)
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: bgc-ingest <health|insert-sample|load-csv|load-xlsx>")
		os.Exit(2) // código 2 -> uso incorreto
	}

	// Seleciona o subcomando pela 1ª palavra e delega para a função correspondente.
	// Os argumentos específicos de `load-csv` (flags) são passados como os.Args[2:].
	var err error
	switch os.Args[1] {
	case "health":
		err = cmdHealth()
	case "insert-sample":
		err = cmdInsertSample()
	case "load-csv":
		err = cmdLoadCSV(os.Args[2:])
	case "load-xlsx":
    	err = cmdLoadXLSX(os.Args[2:])
	default:
		fmt.Fprintln(os.Stderr, "unknown command:", os.Args[1])
		os.Exit(2)
	}

	// Se alguma função retornou erro, imprime no stderr e sai com código 1.
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
