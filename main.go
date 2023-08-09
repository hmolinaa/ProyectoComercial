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
	r.HandleFunc("/send-emails", sendEmailsToStudents).Methods("POST")
	r.HandleFunc("/send-emails_ex", sendEmailsToStudentExcel).Methods("POST")

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

// Funcion para enviiar correos desde la base de datos, extrae los datos de la tabla estudiante

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

// Enviar correos en base a la base de datos, tabla students
func sendEmailsToStudents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Obtener la plantilla del correo electrónico desde la solicitud POST
	emailTemplate := r.PostFormValue("emailTemplate")
	subject := r.PostFormValue("subject") // Obtener el asunto del correo

	// Establecer una función para obtener los datos de los estudiantes desde la base de datos
	studentsData, err := getStudentsData("students")
	if err != nil {
		http.Error(w, "Error al obtener los datos de los estudiantes", http.StatusInternalServerError)
		return
	}

	for _, student := range studentsData {
		// Reemplazar los marcadores de la plantilla con los valores de los estudiantes
		personalizedContent := strings.Replace(emailTemplate, "<<Nombre>>", student.Name, -1)
		personalizedContent = strings.Replace(personalizedContent, "<<Parcial1>>}}", strconv.Itoa(student.First_partial), -1)
		personalizedContent = strings.Replace(personalizedContent, "<<Parcial2>>", strconv.Itoa(student.Second_partial), -1)
		personalizedContent = strings.Replace(personalizedContent, "<<Parcial3>>", strconv.Itoa(student.Third_partial), -1)
		personalizedContent = strings.Replace(personalizedContent, "<<Nota Final>>", strconv.Itoa(student.Final_score), -1)
		personalizedContent = strings.Replace(personalizedContent, "<<Asignatura>>", student.Subject, -1)
		personalizedContent = strings.Replace(personalizedContent, "\n", "<br>", -1)

		htmlContent := personalizedContent
		err := sendEmail(student.Email, subject, htmlContent)
		if err != nil {
			log.Printf("Error sending mail to %s (%s): %v", student.Name, student.Email, err)
		} else {
			log.Printf("send email to %s (%s)", student.Name, student.Email)
		}
	}

	// Enviar una respuesta al frontend para indicar que los correos electrónicos se enviaron correctamente
	w.Write([]byte("Correos electrónicos enviados correctamente"))
}

// Función para obtener los datos de los estudiantes desde la base de datos
func getStudentsData(table string) ([]Student, error) {
	established_connection := conectionBD()
	var query string

	if table == "students" {
		query = "SELECT * FROM students"
	} else if table == "students_ex" {
		query = "SELECT * FROM students_ex"
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

func printTableData_ex(table string) {
	established_connection := conectionBD()
	var query string

	query = "SELECT * FROM " + table

	rows, err := established_connection.Query(query)
	if err != nil {
		log.Printf("Error al obtener los datos de la tabla: %v", err)
		return
	}
	defer rows.Close()

	columnNames, err := rows.Columns()
	if err != nil {
		log.Printf("Error al obtener los nombres de las columnas: %v", err)
		return
	}

	var data []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columnNames))
		valuePtrs := make([]interface{}, len(columnNames))
		for i := range columnNames {
			valuePtrs[i] = &values[i]
		}

		err := rows.Scan(valuePtrs...)
		if err != nil {
			log.Printf("Error al escanear los valores: %v", err)
			return
		}

		entry := make(map[string]interface{})
		for i, col := range columnNames {
			val := values[i]
			entry[col] = val
		}
		data = append(data, entry)
	}

	// Iterar a través de los datos y mostrarlos en la consola
	for i, row := range data {
		log.Printf("Fila %d:", i+1)
		for colName, value := range row {
			log.Printf("%s: %s", colName, string(value.([]byte)))
		}
	}
}

// Función para obtener y enviar correos electrónicos basados en los datos de una tabla desconocida
func sendEmailsToStudents_ex(email string, subject1 string) error {
	established_connection := conectionBD()

	rows, err := established_connection.Query("SELECT * FROM students_ex")
	if err != nil {
		return err
	}
	defer rows.Close()

	columnNames, err := rows.Columns()
	if err != nil {
		return err
	}

	var data []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columnNames))
		valuePtrs := make([]interface{}, len(columnNames))
		for i := range columnNames {
			valuePtrs[i] = &values[i]
		}

		err := rows.Scan(valuePtrs...)
		if err != nil {
			return err
		}

		entry := make(map[string]interface{})
		for i, col := range columnNames {
			val := values[i]
			entry[col] = val
		}
		data = append(data, entry)
	}

	// Iterar a través de los datos y enviar correos electrónicos
	for _, row := range data {

		for colName, value := range row {
			// Convertir valor de bytes a cadena de texto
			valueStr := ""
			if byteValue, ok := value.([]byte); ok {
				valueStr = string(byteValue)
			}
			log.Printf("%s: %s", colName, valueStr)
		}

		// Reemplazar marcadores en la plantilla de correo con los valores de la fila
		personalizedContent := email
		for colName, value := range row {
			personalizedContent = strings.Replace(personalizedContent, "<<"+colName+">>", fmt.Sprintf("%s", value), -1)
			personalizedContent = strings.Replace(personalizedContent, "<<"+colName+">>", fmt.Sprintf("%s", value), -1)
			personalizedContent = strings.Replace(personalizedContent, "<<"+colName+">>", fmt.Sprintf("%s", value), -1)
			personalizedContent = strings.Replace(personalizedContent, "<<"+colName+">>", fmt.Sprintf("%s", value), -1)
			personalizedContent = strings.Replace(personalizedContent, "<<"+colName+">>", fmt.Sprintf("%s", value), -1)
			personalizedContent = strings.Replace(personalizedContent, "<<"+colName+">>", fmt.Sprintf("%s", value), -1)
			personalizedContent = strings.Replace(personalizedContent, "<<"+colName+">>", fmt.Sprintf("%s", value), -1)
			personalizedContent = strings.Replace(personalizedContent, "<<"+colName+">>", fmt.Sprintf("%s", value), -1)
			personalizedContent = strings.Replace(personalizedContent, "<<"+colName+">>", fmt.Sprintf("%s", value), -1)
			personalizedContent = strings.Replace(personalizedContent, "<<"+colName+">>", fmt.Sprintf("%s", value), -1)
			personalizedContent = strings.Replace(personalizedContent, "<<"+colName+">>", fmt.Sprintf("%s", value), -1)
			personalizedContent = strings.Replace(personalizedContent, "<<"+colName+">>", fmt.Sprintf("%s", value), -1)
			personalizedContent = strings.Replace(personalizedContent, "<<"+colName+">>", fmt.Sprintf("%s", value), -1)
			personalizedContent = strings.Replace(personalizedContent, "\n", "<br>", -1)

		}
		personalizedContent = strings.Replace(personalizedContent, "\n", "<br>", -1)

		htmlContent := personalizedContent
		for colName, value := range row {
			subject1 = strings.Replace(subject1, "<<"+colName+">>", fmt.Sprintf("%s", value), -1)
		}

		emailValueBytes, ok := row["Correo"].([]byte)
		if !ok {
			log.Printf("Tipo de dato no válido en columna Correo: %T", row["Correo"])
			continue
		}

		emailValue := string(emailValueBytes)
		if emailValue == "" {
			log.Printf("Valor de email vacío en fila")
			continue
		}

		err := sendEmail(emailValue, subject1, htmlContent)
		if err != nil {
			log.Printf("Error al enviar correo a %s: %v", emailValue, err)
		} else {
			log.Printf("Correo enviado a %s", emailValue)
		}

	}

	return nil
}

func sendEmailsToStudentExcel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtener la plantilla del correo electrónico y el asunto desde la solicitud POST
	email := r.PostFormValue("email")
	subject1 := r.PostFormValue("subject1")

	// Enviar correos electrónicos basados en los datos de la tabla "students_ex"
	err := sendEmailsToStudents_ex(email, subject1)
	if err != nil {
		http.Error(w, "Error al enviar correos electrónicos: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Enviar una respuesta al frontend para indicar que los correos electrónicos se enviaron correctamente
	w.Write([]byte("Correos electrónicos enviados correctamente"))
}
