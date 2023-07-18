package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	//"net/smtp"

	"github.com/go-gomail/gomail"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/cors"
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

//var plantillas = template.Must(template.ParseGlob("plantillas/*"))

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/inicio", Home)
	mux.HandleFunc("/send-emails", Email_Student)
	//mux.HandleFunc("/headers", headers)

	fmt.Println("Servidor en ejecución en http://localhost:8080")

	// cors.Default() setup the middleware with default options being
	// all origins accepted with simple methods (GET, POST). See
	// documentation below for more options.
	handler := cors.Default().Handler(mux)
	http.ListenAndServe(":8080", handler)

}

func Home_website(w http.ResponseWriter, r *http.Request) {
	//plantillas.ExecuteTemplate(w, "Home_website", nil)

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

func sendEmail(name, email, subject, body string) error {
	// Configurar el cliente de envío de correos electrónicos
	dialer := gomail.NewDialer("smtp.gmail.com", 587, "henrrymolina100@gmail.com", "mfhbgbzqdqpqpbtj")

	// Crear el mensaje de correo personalizado
	message := gomail.NewMessage()
	message.SetHeader("From", "henrrymolina100@gmail.com")
	message.SetHeader("To", email)
	message.SetHeader("Subject", subject)
	message.SetBody("text/plain", body)

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
	rows, err := established_connection.Query("SELECT name, email, final_score  FROM students")

	if err != nil {
		panic(err.Error())
	}

	// Recorrer los registros y enviar correos electrónicos personalizados
	for rows.Next() {
		var student Student
		err := rows.Scan(&student.Name, &student.Email, &student.Final_score)
		if err != nil {
			log.Printf("Error al escanear el registro: %v", err)
			continue
		}

		err = sendEmail(student.Name, student.Email, "Calificaciones de la clase MM-520", "Correo enviado desde Go")
		if err != nil {
			log.Printf("Error al enviar el correo a %s (%s): %v", student.Name, student.Email, err)
		} else {
			log.Printf("Correo enviado a %s (%s)", student.Name, student.Email)
		}
	}

	fmt.Println("Proceso de envío de correos completado")

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
