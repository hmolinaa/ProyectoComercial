package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-gomail/gomail"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// Database connection
func conectionBD() (conection *sql.DB) {
	Driver := "mysql"
	User := "root"
	Password := ""
	Name := "unah"

	conection, err := sql.Open(Driver, User+":"+Password+"@tcp(127.0.0.1)/"+Name)
	if err != nil {
		panic(err.Error())
	}
	return conection
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/excel", procesarJSON).Methods(http.MethodPost)
	r.HandleFunc("/inicio", Home)
	r.HandleFunc("/send-emails", Email_Student)
	r.HandleFunc("/send", Email_Student_E)

	corsHandler := cors.Default().Handler(r)
	fmt.Println("Servidor en ejecución en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}

// Structure of each student
type Student struct {
	Id             int
	Name           string
	Account        int
	Subject        string
	First_partial  int
	Second_partial int
	Third_partial  int
	Final_score    int
	Email          string
}

// Process JSON from excel sent from frontend
func procesarJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var jsonData []map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&jsonData)
	if err != nil {
		http.Error(w, "Error reading the JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = save_in_BD(jsonData)
	if err != nil {
		http.Error(w, "Failed to save to database: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a response to the frontend
	w.WriteHeader(http.StatusOK)

	// Send the JSON back to the frontend
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jsonData)
}

func save_in_BD(data []map[string]interface{}) error {
	// Database connection
	db := conectionBD()
	defer db.Close()

	// Clear the table "students_excel"
	_, err := db.Exec("DELETE FROM students_excel")
	if err != nil {
		log.Println("Error cleaning table students_excel:", err)
		return err
	}

	// Insert the data in the table "students_excel"
	for _, item := range data {
		_, err := db.Exec("INSERT INTO students_excel (name, account, subject, first_partial, second_partial, third_partial, final_score, email) VALUES ( ?, ?, ?, ?, ?, ?, ?, ?)",
			item["Nombre"], item["Cuenta"], item["Asignatura"], item["Parcial 1"], item["Parcial 2"], item["Parcial 3"], item["Nota Final"], item["Correo"])
		if err != nil {
			log.Println("Error inserting into database:", err)
		}
	}
	return nil
}

func Home(w http.ResponseWriter, req *http.Request) {
	// Database connection
	established_connection := conectionBD()
	records, err := established_connection.Query("SELECT * FROM students")

	if err != nil {
		panic(err.Error())
	}

	student := Student{}
	ArrayStudent := []Student{}

	for records.Next() {
		var id, account, first_partial, second_partial, third_partial, final_score int
		var name, subject, email string
		err = records.Scan(&id, &name, &account, &subject, &first_partial, &second_partial, &third_partial, &final_score, &email)
		if err != nil {
			panic(err.Error())
		}
		student.Id = id
		student.Name = name
		student.Account = account
		student.Subject = subject
		student.First_partial = first_partial
		student.Second_partial = second_partial
		student.Third_partial = third_partial
		student.Final_score = final_score
		student.Email = email

		ArrayStudent = append(ArrayStudent, student)

	}
	// Convert the array to JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ArrayStudent)
}

func sendEmail(email, subject, body string) error {
	// Configure sending emails
	dialer := gomail.NewDialer("smtp.gmail.com", 587, "henrrymolina100@gmail.com", "mfhbgbzqdqpqpbtj")

	// Create the custom email
	message := gomail.NewMessage()
	message.SetHeader("From", "henrrymolina100@gmail.com")
	message.SetHeader("To", email)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", body)

	// send email
	err := dialer.DialAndSend(message)
	if err != nil {
		return err
	}

	return nil
}

func Email_Student(w http.ResponseWriter, r *http.Request) {
	established_connection := conectionBD()
	rows, err := established_connection.Query("SELECT * FROM students") // we access the table students

	if err != nil {
		panic(err.Error())
	}

	for rows.Next() {
		var student Student
		err := rows.Scan(&student.Id, &student.Name, &student.Account, &student.Subject, &student.First_partial, &student.Second_partial, &student.Third_partial, &student.Final_score, &student.Email)
		if err != nil {
			log.Printf("Failed to scan registry: %v", err)
			continue
		}

		htmlContent := generateEmailContent(student)

		err = sendEmail(student.Email, "Calificaciones de la clase MM-520", htmlContent)
		if err != nil {
			log.Printf("Error sending mail to %s (%s): %v", student.Name, student.Email, err)
		} else {
			log.Printf("send email to %s (%s)", student.Name, student.Email)
		}
	}

	fmt.Println("Mailing process completed")
}

// Access the table students_excel
func Email_Student_E(w http.ResponseWriter, r *http.Request) {
	established_connection := conectionBD()
	rows, err := established_connection.Query("SELECT * FROM students_excel") //we access the table students_excel

	if err != nil {
		panic(err.Error())
	}

	for rows.Next() {
		var student Student
		err := rows.Scan(&student.Id, &student.Name, &student.Account, &student.Subject, &student.First_partial, &student.Second_partial, &student.Third_partial, &student.Final_score, &student.Email)
		if err != nil {
			log.Printf("Failed to scan registry: %v", err)
			continue
		}

		htmlContent := generateEmailContent(student)

		err = sendEmail(student.Email, "Calificaciones", htmlContent)
		if err != nil {
			log.Printf("Error sending mail to %s (%s): %v", student.Name, student.Email, err)
		} else {
			log.Printf("sen email to %s (%s)", student.Name, student.Email)
		}
	}

	fmt.Println("Mailing process completed")

}

// Falta cambiar, esto lo debe editar el usuarios desde el frontend
func generateEmailContent(student Student) string {

	P1 := strconv.Itoa(student.First_partial)
	P2 := strconv.Itoa(student.Second_partial)
	P3 := strconv.Itoa(student.Third_partial)
	Final_score := strconv.Itoa(student.Final_score)

	htmlContent := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Correo Electrónico Personalizado</title>
	</head>
	<body>
		<h1>Hola, ` + student.Name + `</h1>
		<p>A continuación se muestran sus datos y calificaciones obtenidas:</p>
		<p>Nombre: ` + student.Name + `</p>
		<p>Correo Electrónico: ` + student.Email + `</p>
		<p>Asignatura: ` + student.Subject + `</p>
		<p>Su nota del primer parcial es: ` + P1 + `</p>
		<p>Su nota del segundo parcial es: ` + P2 + `</p>
		<p>Su nota del tercer parcial es: ` + P3 + `</p>
		<p>Su nota final es: ` + Final_score + `</p>
		
		<p>Saludos</p>
	</body>
	</html>
	`

	return htmlContent
}
