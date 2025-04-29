package main

import (
	"database/sql" //pacote para conectar e trabalhar com bancos de dados.
	//"encoding/json" //permite codificar/decodificar JSON (entrada e saída dos dados via API).
	"fmt" //imprime menssagens no terminal
	//"log"           // usado para registrar erros e parar o programa se algo der errado.
	//"net/http"      //cria o servidor web e lida com as rotas da API.
	_ "github.com/mattn/go-sqlite3" // importa o driver do banco SQLite, o _ é porque só queremos registrar o driver, sem usar diretamente.
	//Tambbém instalei a dependencia via terminal: go get github.com/mattn/go-sqlite3
	//permite que a linguagem Go consiga se conectar e fazer consultas em um banco de dados SQLite.
)

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
	// Mostra no terminal que a tabela foi criada (ou já existia)
	fmt.Println("Tabela criada ou já existia:", resultadoTabela)

	// Insere uma operação na tabela: 2 + 3 = 5
	// Os "?" são placeholders que evitam SQL injection
	resultadoInsert, err := db.Exec(
		"INSERT INTO operacoes (operando1, operando2, operacao, resultado) VALUES (?, ?, ?, ?)",
		2, 3, "+", 5,
	)
	if err != nil {
		// Se houver erro ao inserir, o programa para
		panic(err)
	}

	// Verifica quantas linhas foram modificadas (inseridas)
	linhasModificadas, err := resultadoInsert.RowsAffected()
	if err != nil {
		// Se não conseguir contar as linhas, mostra erro e sai da função
		fmt.Println("Erro ao obter linhas modificadas:", err)
		return
	}
	// Mostra o número de linhas modificadass (deve ser 1)
	fmt.Println("Inserção feita. Linhas modificadas:", linhasModificadas)

	// Consulta todos os dados da tabela "operacoes"
	linhas, err := db.Query("SELECT id, operando1, operando2, operacao, resultado FROM operacoes")
	if err != nil {
		// Se der erro na consulta, o programa para
		panic(err)
	}
	// Garante que os resultados da consulta serão fechados no final
	defer linhas.Close()

	// Cabeçalho para os resultados
	fmt.Println("\n📋 Resultados:")

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

/*
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

// Estrutura do Tipo HistoricoItem que representa cada operação salva no banco de dados.
// A ideia é que o ID, seja gerado automaticamente.
type HistoricoItem struct {
	ID        int     `json:"id"`
	Operando1 float64 `json:"operando1"`
	Operando2 float64 `json:"operando2"`
	Operacao  string  `json:"operacao"`
	Resultado float64 `json:"resultado"`
}

// Funcao principal que inicia o Servidor
func main() {
	// variavel err que cria o banco de dados SQLite
	var err error
	db, err = sql.Open("sqlite3", "./calculadora.db") //Configura o acesso para quando vocêo cliente fizer a primeira operação.

	//Se der erro, ele é guardado na variável err criado na linha 41.
	if err != nil {
		log.Fatal("Erro ao abrir o banco:", err)
		//Se hoouver erro ao abrir o banco, para o programa e mostra o erro.
	}
	defer db.Close() //Garante que o banco de dados seja encerrado quando o programa terminar.

	http.HandleFunc("/soma", somaHandler) //rota de cada api
	http.HandleFunc("/subtracao", subtracaoHandler)
	http.HandleFunc("/multiplicacao", multiplicacaoHandler)
	http.HandleFunc("/divisao", divisaoHandler)

	fmt.Println("Servidor rodando na porta 8080...") // funcao para mostrar que o servidor está rodando na porta 8080
	http.ListenAndServe(":8080", nil)
	// funcao para mostrar que o servidor está rodando na porta 8080
}

var db *sql.DB // Criar uma variável global db que vai armazenar a conexão aberta com o banco de dados."

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
}*/
