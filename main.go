package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	//_ "github.com/mattn/go-sqlite3" // Importa o driver SQLite
	_ "modernc.org/sqlite"
)

// Operacao representa um registro no histórico
type Operacao struct {
	ID        int
	Operando1 float64
	Operando2 float64
	Operacao  string
	Resultado float64
	Timestamp string
}

// registrarOperacao insere uma nova operação no banco de dados
func registrarOperacao(db *sql.DB, operando1 float64, operando2 float64, operacao string, resultado float64) error {
	stmt, err := db.Prepare("INSERT INTO historico(operando1, operando2, operacao, resultado, timestamp) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	timestamp := time.Now().Format(time.RFC3339)
	_, err = stmt.Exec(operando1, operando2, operacao, resultado, timestamp)
	return err
}

// listarHistorico recupera e exibe todas as operações do banco de dados
func listarHistorico(db *sql.DB) error {
	rows, err := db.Query("SELECT id, operando1, operando2, operacao, resultado, timestamp FROM historico ORDER BY timestamp DESC")
	if err != nil {
		return err
	}
	defer rows.Close()

	fmt.Println("\nHistórico de Operações:")
	for rows.Next() {
		var op Operacao
		err = rows.Scan(&op.ID, &op.Operando1, &op.Operando2, &op.Operacao, &op.Resultado, &op.Timestamp)
		if err != nil {
			return err
		}
		fmt.Printf("ID: %d, Operando 1: %.2f, Operando 2: %.2f, Operação: %s, Resultado: %.2f, Timestamp: %s\n",
			op.ID, op.Operando1, op.Operando2, op.Operacao, op.Resultado, op.Timestamp)
	}

	return nil
}

func main() {
	// Abre a conexão com o banco de dados SQLite (cria o arquivo se não existir)
	db, err := sql.Open("sqlite", "historico_calculadora.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Cria a tabela 'historico' se ela não existir
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS historico (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			operando1 REAL NOT NULL,
			operando2 REAL NOT NULL,
			operacao TEXT NOT NULL,
			resultado REAL NOT NULL,
			timestamp TEXT NOT NULL
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Exemplos de uso
	err = registrarOperacao(db, 5, 3, "+", 8)
	if err != nil {
		log.Println("Erro ao registrar operação:", err)
	}
	err = registrarOperacao(db, 10, 2, "-", 8)
	if err != nil {
		log.Println("Erro ao registrar operação:", err)
	}
	err = registrarOperacao(db, 4, 6, "*", 24)
	if err != nil {
		log.Println("Erro ao registrar operação:", err)
	}
	err = registrarOperacao(db, 9, 3, "/", 3)
	if err != nil {
		log.Println("Erro ao registrar operação:", err)
	}

	// Lista o histórico
	err = listarHistorico(db)
	if err != nil {
		log.Println("Erro ao listar histórico:", err)
	}
}
