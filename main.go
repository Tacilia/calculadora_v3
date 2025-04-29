package main

import (
	"database/sql"  //pacote para conectar e trabalhar com bancos de dados.
	"encoding/json" //permite codificar/decodificar JSON (entrada e saída dos dados via API).
	"fmt"           //imprime menssagens no terminal
	"net/http"      //cria o servidor web e lida com as rotas da API.

	_ "github.com/mattn/go-sqlite3" // importa o driver do banco SQLite, o _ é porque só queremos registrar o driver, sem usar diretamente.
	//Tambbém instalei a dependencia via terminal: go get github.com/mattn/go-sqlite3
	//permite que a linguagem Go consiga se conectar e fazer consultas em um banco de dados SQLite.
)

// Estrutura do Tipo OperacaoRequest para receber os dados da operação.
type OperacaoRequest struct {
	Operando1 float64 `json:"operando1"`
	Operando2 float64 `json:"operando2"`
	Operacao  string  `json:"operacao"`
}

// Estrutura do Tipo ResultadoResponse para enviar o resultado.
type ResultadoResponse struct {
	Resultado float64 `json:"resultado"`
}

// Função que responde a um pedido (requisição) que chega no seu servidor.
func somaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	} // Checando se a requiscao é metodo POST, se nao for envia a menssagem metodo nao permitido.

	var req OperacaoRequest //req é uma variável do tipo OperacaoRequest.
	// Serve para ler informações que chegaram na requisição.
	json.NewDecoder(r.Body).Decode(&req)
	// Le dados Json que chegaram na requisicao e transforma em um objeto/struct Go.

	resultado := req.Operando1 + req.Operando2 //Crie uma variável chamada resultado e guarde dentro dela o valor que chegou como resposta na requisição.

	json.NewEncoder(w).Encode(ResultadoResponse{Resultado: resultado})
	// Envia e decodifica uma resposta para o formato JSON para quem chamou a API.
}

func subtracaoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Use o método POST", http.StatusMethodNotAllowed)
		return
	}

	var req OperacaoRequest
	json.NewDecoder(r.Body).Decode(&req)

	resultado := req.Operando1 - req.Operando2

	json.NewEncoder(w).Encode(ResultadoResponse{Resultado: resultado})
}

func multiplicacaoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Use o método POST", http.StatusMethodNotAllowed)
		return
	}

	var req OperacaoRequest
	json.NewDecoder(r.Body).Decode(&req)

	resultado := req.Operando1 * req.Operando2

	json.NewEncoder(w).Encode(ResultadoResponse{Resultado: resultado})
}

func divisaoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Use o método POST", http.StatusMethodNotAllowed)
		return
	}

	var req OperacaoRequest
	json.NewDecoder(r.Body).Decode(&req)

	//Estrutura condicional que coloca uma condição: náo é permitido que o operando2 seja igual 0.
	if req.Operando2 == 0 {
		http.Error(w, "Erro: nenhum número pode ser dividido por zero", http.StatusBadRequest)
		return
	}

	resultado := req.Operando1 / req.Operando2

	json.NewEncoder(w).Encode(ResultadoResponse{Resultado: resultado})
}

func main() {
	// Abre (ou cria) o banco de dados SQLite com nome "calculadora.db"
	// A função sql.Open retorna dois valores: um ponteiro para o banco (db) e um erro (err)
	db, err := sql.Open("sqlite3", "calculadora.db")
	if err != nil {
		// Se houver erro ao abrir o banco, o programa é encerrado com panic
		panic(err)
	}
	// Garante que o banco de dados será fechado ao final da função main
	defer db.Close()

	http.HandleFunc("/soma", somaHandler) //rota de cada api
	http.HandleFunc("/subtracao", subtracaoHandler)
	http.HandleFunc("/multiplicacao", multiplicacaoHandler)
	http.HandleFunc("/divisao", divisaoHandler)

	fmt.Println("Servidor rodando na porta 8080...") // mostrar que o servidor está rodando na porta 8080
	http.ListenAndServe(":8080", nil)

	// Comando SQL para criar a tabela "operacoes" caso ela ainda não exista
	sqlTabela := `CREATE TABLE IF NOT EXISTS operacoes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,  -- ID único, gerado automaticamente
		operando1 REAL,                        -- Primeiro número da operação (float)
		operando2 REAL,                        -- Segundo número da operação (float)
		operacao TEXT,                         -- Tipo da operação: "+", "-", "*", "/"
		resultado REAL                         -- Resultado do cálculo
	);`

	// Executa o comando SQL para criar a tabela
	resultadoTabela, err := db.Exec(sqlTabela)
	if err != nil {
		// Se ocorrer erro, o programa para imediatamente
		panic(err)
	}
	// Mostra no terminal que a tabela foi criada
	fmt.Println("Tabela criada:", resultadoTabela)

	// Definindo os valores da operação
	operando1 := 2.0
	operando2 := 3.0
	operacao := "+"
	resultado := operando1 + operando2

	// Insere uma operação na tabela: 2 + 3 = 5
	// Executa o comando INSERT usando os valores acima
	resultadoInsert, err := db.Exec(
		"INSERT INTO operacoes (operando1, operando2, operacao, resultado) VALUES (?, ?, ?, ?)",
		operando1, operando2, operacao, resultado,
	)
	if err != nil {
		// Se houver erro ao inserir, o programa para
		panic(err)
	}

	// Verifica quantas linhas foram inseridas (modificadas)
	linhasModificadas, err := resultadoInsert.RowsAffected()
	if err != nil {
		// Se não conseguir contar as linhas, mostra erro e sai da função
		fmt.Println("Erro ao obter linhas modificadas:", err)
		return
	}
	// Mostra o número de linhas inseridas (deve ser 1)
	fmt.Println("Inserção feita. quantidade de linhas modificadas:", linhasModificadas)

	// Consulta todos os dados da tabela "operacoes"
	linhas, err := db.Query("SELECT id, operando1, operando2, operacao, resultado FROM operacoes")
	if err != nil {
		// Se der erro na consulta, o programa para
		panic(err)
	}
	// Garante que os resultados da consulta serão fechados no final
	defer linhas.Close()

	// Percorre cada linha retornada pelo SELECT
	for linhas.Next() {
		var id int          // Armazena o ID da operação
		var a, b, r float64 // Armazena os operandos e o resultado
		var op string       // Armazena o tipo da operação

		// Pega os valores da linha atual e coloca nas variáveis acima
		linhas.Scan(&id, &a, &b, &op, &r)

		// Imprime a operação formatada, exemplo: "1 -> 2.0 + 3.0 = 5.0"
		fmt.Printf("%d -> %.1f %s %.1f = %.1f\n", id, a, op, b, r)
	}
}
