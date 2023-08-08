// main.go
package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

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
	r.HandleFunc("/send_emails/{table}", Email_Student).Methods("POST")
	r.HandleFunc("/send-emails", sendEmailsToStudents).Methods("POST")

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

	err = insertDataIntoTable(jsonData, "unah", "students_ex")
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

func insertDataIntoTable(data []map[string]interface{}, dbName string, tableName string) error {
	// Conexión con la base de datos
	db := conectionBD()
	defer db.Close()

	// Selecciona la base de datos
	_, err := db.Exec("USE " + dbName)
	if err != nil {
		return err
	}

	// Elimina la tabla si existe
	_, err = db.Exec("DROP TABLE IF EXISTS " + tableName)
	if err != nil {
		return err
	}

	// Infiera la estructura de la tabla a partir del primer elemento del JSON
	firstItem := data[0]
	createTableQuery := "CREATE TABLE IF NOT EXISTS " + tableName + " (id INT AUTO_INCREMENT PRIMARY KEY"

	for key, value := range firstItem {
		var dataType string

		switch value.(type) {
		case int, int64:
			dataType = "INT"
		default:
			dataType = "VARCHAR(255)"
		}

		// Rodea los nombres de columna con comillas inversas para manejar caracteres especiales
		createTableQuery += ", `" + key + "` " + dataType
	}

	createTableQuery += ")"

	_, err = db.Exec(createTableQuery)
	if err != nil {
		return err
	}

	// Inserta los datos en la tabla y imprime en la consola
	for _, item := range data {
		insertColumns := []string{}
		insertValues := []interface{}{}

		for key, value := range item {
			insertColumns = append(insertColumns, key)
			insertValues = append(insertValues, value)
		}

		// Construye la cadena de columnas
		columnStr := "`" + strings.Join(insertColumns, "`, `") + "`"

		// Construye la cadena de marcadores de posición (?)
		valuePlaceholders := strings.Repeat("?, ", len(insertColumns)-1) + "?"

		// Construye la consulta INSERT
		insertQuery := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, columnStr, valuePlaceholders)

		// Inserta los datos en la tabla
		_, err := db.Exec(insertQuery, insertValues...)
		if err != nil {
			log.Println("Error al insertar en la base de datos:", err)
		} else {
			fmt.Println("Datos insertados:", insertValues) // Imprime los datos insertados en la consola
		}
	}

	return nil
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

// Funcion para enviiar correos desde la base de datos, extrae los datos de la tabla estudiante
func Email_Student(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	table := vars["table"]

	// Parsear el cuerpo de la solicitud JSON
	var data map[string]string
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Establecer conexión a la base de datos
	established_connection := conectionBD()
	var rows *sql.Rows

	// Seleccionar la tabla y campos correctos según la información proporcionada
	if table == "students" {
		rows, err = established_connection.Query("SELECT * FROM students")
	} else if table == "students_excel" {
		rows, err = established_connection.Query("SELECT * FROM students_excel ")
	} else {
		http.Error(w, "Invalid table name", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "Error querying the database", http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		var person Student
		err := rows.Scan(&person.Id, &person.Name, &person.Account, &person.Subject, &person.First_partial, &person.Second_partial, &person.Third_partial, &person.Final_score, &person.Email)
		if err != nil {
			log.Printf("Failed to scan registry: %v", err)
			continue
		}
		subject := data["subject"]
		professorName := data["professorName"]
		c_name := data["c_name"]
		c_subject := data["c_subject"]
		desp := data["desp"]

		// Generar el contenido del correo con el contenido personalizado del usuario
		htmlContent := generateEmailContent(person, data["customMessage"], professorName, c_name, c_subject, desp)

		// Enviar el correo electrónico
		err = sendEmail(person.Email, subject, htmlContent)
		if err != nil {
			log.Printf("Error sending mail to %s (%s): %v", person.Name, person.Email, err)
		} else {
			log.Printf("send email to %s (%s)", person.Name, person.Email)
		}
	}

	// Enviar respuesta exitosa al frontend
	response := map[string]string{"message": "Emails sent successfully"}
	jsonResponse, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

}

func generateEmailContent(student Student, customContent, professorName, c_name, c_subject, desp string) string {
	P1 := strconv.Itoa(student.First_partial)
	P2 := strconv.Itoa(student.Second_partial)
	P3 := strconv.Itoa(student.Third_partial)
	Final_score := strconv.Itoa(student.Final_score)

	// Reemplazar los marcadores de posición en la plantilla HTML.
	htmlContent := strings.Replace(templateHTML, "{{Nombre}}", student.Name, -1)
	htmlContent = strings.Replace(htmlContent, "{{Asignatura}}", student.Subject, -1)
	htmlContent = strings.Replace(htmlContent, "{{P1}}", P1, -1)
	htmlContent = strings.Replace(htmlContent, "{{P2}}", P2, -1)
	htmlContent = strings.Replace(htmlContent, "{{P3}}", P3, -1)
	htmlContent = strings.Replace(htmlContent, "{{Final_score}}", Final_score, -1)
	htmlContent = strings.Replace(htmlContent, "{{CustomContent}}", customContent, -1)

	return htmlContent
}

const templateHTML = `
<!DOCTYPE html>
<html>

<body>
    <div id="customMessage" style="border: 1px solid #ccc; padding: 10px;">
         {{Nombre}}
        <br><br>
        {{CustomContent}}
        <br>
		<br>
        {{c_subject}} {{Asignatura}}
        <br>
		

		<table border="1" cellpadding="5">
        <thead >
          <tr>
            <th >Primer Parcial</th>
            <th >Segundo Parcial</th>
            <th >Tercer Parcial</th>
            <th>Nota Final</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td>{{P1}}</td>
            <td>{{P2}}</td>
            <td>{{P3}}</td>
            <td>{{Final_score}}</td>
          </tr>
        </tbody>
      </table>

	  <br>
      {{desp}}
	  <br>
	  {{Profesor}}
    </div>
</body>
</html>
`

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

func sendEmailsToStudents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Obtener la plantilla del correo electrónico desde la solicitud POST
	emailTemplate := r.PostFormValue("emailTemplate")

	// Establecer una función para obtener los datos de los estudiantes desde la base de datos
	// En este ejemplo, se asume que tienes una función getStudentsData() que obtiene los datos de los estudiantes desde la base de datos
	studentsData, err := getStudentsData("students")
	if err != nil {
		http.Error(w, "Error al obtener los datos de los estudiantes", http.StatusInternalServerError)
		return
	}

	// Establecer una función para enviar correos electrónicos aquí
	// La función sendEmail() podría ser parte de una biblioteca de envío de correos electrónicos o una API de envío de correos electrónicos

	// Supongamos que tienes una función sendEmail() que puede enviar correos electrónicos
	// Puedes llamar a esta función para enviar los correos electrónicos a los estudiantes
	// Por ejemplo, podrías recorrer los datos de los estudiantes y enviar correos electrónicos uno por uno con el mensaje personalizado

	for _, student := range studentsData {
		// Reemplazar los marcadores de la plantilla con los valores de los estudiantes
		personalizedContent := strings.Replace(emailTemplate, "{{Nombre}}", student.Name, -1)
		personalizedContent = strings.Replace(personalizedContent, "{{P1}}", strconv.Itoa(student.First_partial), -1)
		personalizedContent = strings.Replace(personalizedContent, "{{P2}}", strconv.Itoa(student.Second_partial), -1)
		personalizedContent = strings.Replace(personalizedContent, "{{P3}}", strconv.Itoa(student.Third_partial), -1)
		personalizedContent = strings.Replace(personalizedContent, "{{Final_score}}", strconv.Itoa(student.Final_score), -1)
		personalizedContent = strings.Replace(personalizedContent, "{{Asignatura}}", student.Subject, -1)

		htmlContent := personalizedContent
		err := sendEmail(student.Email, "Calificaciones de la clase MM-520", htmlContent)
		if err != nil {
			log.Printf("Error sending mail to %s (%s): %v", student.Name, student.Email, err)
		} else {
			log.Printf("send email to %s (%s)", student.Name, student.Email)
		}
	}

	// Enviar una respuesta al frontend para indicar que los correos electrónicos se enviaron correctamente
	w.Write([]byte("Correos electrónicos enviados correctamente"))
}

func getStudentsData(table string) ([]Student, error) {
	established_connection := conectionBD()
	var query string

	if table == "students" {
		query = "SELECT * FROM students"
	} else if table == "students_excel" {
		query = "SELECT * FROM students_excel"
	} else {
		return nil, errors.New("Invalid table name")
	}

	rows, err := established_connection.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []Student
	for rows.Next() {
		var student Student
		err := rows.Scan(&student.Id, &student.Name, &student.Account, &student.Subject, &student.First_partial, &student.Second_partial, &student.Third_partial, &student.Final_score, &student.Email)
		if err != nil {
			return nil, err
		}
		students = append(students, student)
	}

	return students, nil
}
