package main

import (
	"database/sql"  //pacote para conectar e trabalhar com bancos de dados.
	"encoding/json" //permite codificar/decodificar JSON (entrada e saída dos dados via API).
	"fmt"           //imprime menssagens no terminal
	"log"           // usado para registrar erros e parar o programa se algo der errado.
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

// Estrutura do Tipo HistoricoItem que representa cada operação salva no banco de dados.
// A ideia é que o ID, seja gerado automaticamente.
type HistoricoItem struct {
	ID        int     `json:"id"`
	Operando1 float64 `json:"operando1"`
	Operando2 float64 `json:"operando2"`
	Operacao  string  `json:"operacao"`
	Resultado float64 `json:"resultado"`
}

var db *sql.DB // Criar uma variável global db que vai armazenar a conexão aberta com o banco de dados."

// Funcao principal que inicia o Servidor
func main() {
	// variavel err que cria o banco de dados SQLite
	var err error
	db, err = sql.Open("sqlite3", "./calculadora.db") //Configura o acesso para quando vocêo cliente fizer a primeira operação.
	// Se o arquivo calculadora.db ainda não existir, o driver do SQLite automaticamente cria esse arquivo na hora da primeira operação.

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

// Função que responde a um pedido (requisição) que chega no seu servidor.
func somaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	} // Checando se a requiscao é metodo POST, se nao for envia a menssagem metodo nao permitido.

	var req OperacaoRequest //req é uma variável do tipo OperacaoRequest.
	// Serve para ler informações que chegaram na requisição.
	json.NewDecoder(r.Body).Decode(&req)
	//json.NewDecoder(r.Body) --> Cria um "leitor" para os dados que chegaram na requisição HTTP.
	//.Decode(&req) --> Lê e decodifica os dados JSONm preenche a variável 'req' usando o endereço http dela.

	resultado := req.Operando1 + req.Operando2 //Crie uma variável chamada resultado e guarde dentro dela o valor que chegou como resposta na requisição.

	json.NewEncoder(w).Encode(ResultadoResponse{Resultado: resultado})
	// Serve para mandar uma resposta no formato JSON para quem chamou a sua API.
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
