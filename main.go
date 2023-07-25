package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/rs/cors"
	"github.com/tealeg/xlsx"

	"github.com/go-gomail/gomail"

	_ "github.com/go-sql-driver/mysql"
)

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
	mux := http.NewServeMux()
	mux.HandleFunc("/inicio", Home)
	mux.HandleFunc("/send-emails", Email_Student)
	//mux.HandleFunc("/subirexcel", handleUpload)

	//mux.HandleFunc("/headers", headers)

	fmt.Println("Servidor en ejecución en http://localhost:8080")

	// cors.Default() setup the middleware with default options being
	// all origins accepted with simple methods (GET, POST). See
	// documentation below for more options.
	handler := cors.Default().Handler(mux)
	http.ListenAndServe(":8080", handler)

}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error al obtener el archivo", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Crear un archivo temporal para guardar el archivo subido
	tempFile, err := os.CreateTemp("", "tempfile.xlsx")
	if err != nil {
		http.Error(w, "Error al crear el archivo temporal", http.StatusInternalServerError)
		return
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Copiar el contenido del archivo subido al archivo temporal
	_, err = io.Copy(tempFile, file)
	if err != nil {
		http.Error(w, "Error al guardar el archivo temporal", http.StatusInternalServerError)
		return
	}

	// Leer el archivo de Excel
	xlFile, err := xlsx.OpenFile(tempFile.Name())
	if err != nil {
		http.Error(w, "Error al abrir el archivo de Excel", http.StatusInternalServerError)
		return
	}

	// Procesar el archivo de Excel y obtener los datos en una tabla
	var tabla [][]string
	for _, sheet := range xlFile.Sheets {
		for _, row := range sheet.Rows {
			var fila []string
			for _, cell := range row.Cells {
				fila = append(fila, cell.String())
			}
			tabla = append(tabla, fila)
		}
	}

	// Convertir la tabla a formato JSON
	jsonData, err := json.Marshal(tabla)
	if err != nil {
		http.Error(w, "Error al convertir los datos a JSON", http.StatusInternalServerError)
		return
	}

	// Impresión para verificar el contenido del JSON
	fmt.Println(string(jsonData))

	// Responder con el archivo.json
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func Delete_student(w http.ResponseWriter, r *http.Request) {
	id_student := r.URL.Query().Get("id")
	//fmt.Println(id_student)

	established_connection := conectionBD()
	delete_records, err := established_connection.Prepare("DELETE FROM students WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delete_records.Exec(id_student)
	http.Redirect(w, r, "/", 301)

}

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

type EmailContent struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func Home(w http.ResponseWriter, req *http.Request) {

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

	//fmt.Println(ArrayStudent)

}

func sendEmail(email, subject, body string) error {
	// Configurar el envío de correos electrónicos
	dialer := gomail.NewDialer("smtp.gmail.com", 587, "henrrymolina100@gmail.com", "mfhbgbzqdqpqpbtj")

	// Crear el mensaje de correo personalizado
	message := gomail.NewMessage()
	message.SetHeader("From", "henrrymolina100@gmail.com")
	message.SetHeader("To", email)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", body)

	// Enviar el correo electrónico
	err := dialer.DialAndSend(message)
	if err != nil {
		return err
	}

	return nil
}

func Edit_Student(w http.ResponseWriter, r *http.Request) {
	idstudent := r.URL.Query().Get("id")
	//fmt.Println(idstudent)

	established_connection := conectionBD()
	record, err := established_connection.Query("SELECT * FROM estudiantes WHERE id=?", idstudent)

	student := Student{}
	for record.Next() {
		var id, account, first_partial, second_partial, third_partial, final_score int
		var name, subject, email string
		err = record.Scan(&id, &name, &account, &subject, &first_partial, &second_partial, &third_partial, &final_score, &email)
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

	}

	//fmt.Println(student)

}

func Email_Student(w http.ResponseWriter, r *http.Request) {

	established_connection := conectionBD()
	rows, err := established_connection.Query("SELECT * FROM students")

	if err != nil {
		panic(err.Error())
	}

	for rows.Next() {
		var student Student
		err := rows.Scan(&student.Id, &student.Name, &student.Account, &student.Subject, &student.First_partial, &student.Second_partial, &student.Third_partial, &student.Final_score, &student.Email)
		if err != nil {
			log.Printf("Error al escanear el registro: %v", err)
			continue
		}

		htmlContent := generateEmailContent(student)

		err = sendEmail(student.Email, "Calificaciones de la clase MM-520", htmlContent)
		if err != nil {
			log.Printf("Error al enviar el correo a %s (%s): %v", student.Name, student.Email, err)
		} else {
			log.Printf("Correo enviado a %s (%s)", student.Name, student.Email)
		}
	}

	fmt.Println("Proceso de envío de correos completado")

}

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

func Create_Student(w http.ResponseWriter, r *http.Request) {

}

func Insert_Student(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("name")
		account := r.FormValue("account")
		subject := r.FormValue("subject")
		first_partial := r.FormValue("first_partial")
		second_partial := r.FormValue("second_partial")
		third_partial := r.FormValue("third_partial")
		final_score := r.FormValue("final_score")
		email := r.FormValue("email")

		established_connection := conectionBD()
		insert_record, err := established_connection.Prepare("INSERT INTO estudiantes( name, account, subject, first_partial, second_partial, third_partial, final_score, email) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			panic(err.Error())
		}
		insert_record.Exec(name, account, subject, first_partial, second_partial, third_partial, final_score, email)
		http.Redirect(w, r, "/", 301)

	}

}

func Update_Student(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("name")
		account := r.FormValue("account")
		subject := r.FormValue("subject")
		first_partial := r.FormValue("first_partial")
		second_partial := r.FormValue("second_partial")
		third_partial := r.FormValue("third_partial")
		final_score := r.FormValue("final_score")
		email := r.FormValue("email")

		established_connection := conectionBD()
		modify_record, err := established_connection.Prepare("UPDATE students SET  name = ?, account, subject = ?, first_partial = ?, second_partial = ?, third_partial = ?, final_score = ?, email = ? WHERE id=? ")
		if err != nil {
			panic(err.Error())
		}
		modify_record.Exec(name, account, subject, first_partial, second_partial, third_partial, final_score, email)
		http.Redirect(w, r, "/", 301)

	}

}
