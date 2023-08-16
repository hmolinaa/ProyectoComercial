// main.go
package main

//Dependens
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

// main
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/inicio", Home)
	r.HandleFunc("/excel", procesarJSON).Methods(http.MethodPost)
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
	// Database connection
	db := conectionBD()
	defer db.Close()

	// Select the database
	_, err := db.Exec("USE " + dbName)
	if err != nil {
		return err
	}

	// Delete the table if it exists
	_, err = db.Exec("DROP TABLE IF EXISTS " + tableName)
	if err != nil {
		return err
	}

	// Infer the structure of the table from the first element of the JSON
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

		// Surround column names with back quotes to handle special characters
		createTableQuery += ", `" + key + "` " + dataType
	}

	createTableQuery += ")"

	_, err = db.Exec(createTableQuery)
	if err != nil {
		return err
	}

	// Insert the data into the table and print to the console
	for _, item := range data {
		insertColumns := []string{}
		insertValues := []interface{}{}

		for key, value := range item {
			insertColumns = append(insertColumns, key)
			insertValues = append(insertValues, value)
		}

		// Build the column chain
		columnStr := "`" + strings.Join(insertColumns, "`, `") + "`"

		// Constructs the string of placeholders (?)
		valuePlaceholders := strings.Repeat("?, ", len(insertColumns)-1) + "?"

		// Build the INSERT query
		insertQuery := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, columnStr, valuePlaceholders)

		// Insert the data into the table
		_, err := db.Exec(insertQuery, insertValues...)
		if err != nil {
			log.Println("Error inserting into database:", err)
		} else {
			fmt.Println("inserted data:", insertValues) // Print the data inserted in the console
		}
	}

	return nil
}

// Send emails based on the database, students table
func sendEmailsToStudents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get email template from POST request
	password := r.PostFormValue("password")
	Femail := r.PostFormValue("Femail")
	emailTemplate := r.PostFormValue("emailTemplate")
	subject := r.PostFormValue("subject") // Get the subject of the email

	//Set up a function to get student data from the database
	studentsData, err := getStudentsData("students")
	if err != nil {
		http.Error(w, "Error getting student data", http.StatusInternalServerError)
		return
	}

	for _, student := range studentsData {
		// Replace template markers with student values
		personalizedContent := strings.Replace(emailTemplate, "<<Nombre>>", student.Name, -1)
		personalizedContent = strings.Replace(personalizedContent, "<<Parcial-1>>}}", strconv.Itoa(student.First_partial), -1)
		personalizedContent = strings.Replace(personalizedContent, "<<Parcial-2>>", strconv.Itoa(student.Second_partial), -1)
		personalizedContent = strings.Replace(personalizedContent, "<<Parcial-3>>", strconv.Itoa(student.Third_partial), -1)
		personalizedContent = strings.Replace(personalizedContent, "<<Nota Final>>", strconv.Itoa(student.Final_score), -1)
		personalizedContent = strings.Replace(personalizedContent, "<<Asignatura>>", student.Subject, -1)
		personalizedContent = strings.Replace(personalizedContent, "\n", "<br>", -1)

		htmlContent := personalizedContent
		err := sendEmail(Femail, password, student.Email, subject, htmlContent)
		if err != nil {
			log.Printf("Error sending mail to %s (%s): %v", student.Name, student.Email, err)
		} else {
			log.Printf("send email to %s (%s)", student.Name, student.Email)
		}
	}

	// Send a response to the frontend to indicate that the emails were sent successfully
	w.Write([]byte("Emails sent successfully"))
}

// Function to get student data from the database
func getStudentsData(table string) ([]Student, error) {
	established_connection := conectionBD()
	var query string

	if table == "students" {
		query = "SELECT * FROM students"
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

// function to send emails to the excel table that we do not know its data
func sendEmailsToStudentExcel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get email template and subject from POST request
	password1 := r.PostFormValue("password1")
	Femail1 := r.PostFormValue("Femail1")
	email := r.PostFormValue("email")
	subject1 := r.PostFormValue("subject1")

	// Send emails based on the data in the "students_ex" table
	err := sendEmailsToStudents_ex(Femail1, password1, email, subject1)
	if err != nil {
		http.Error(w, "Error sending emails: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a response to the frontend to indicate that the emails were sent successfully
	w.Write([]byte("Emails sent successfully"))
}

// function to send emails to the excel table that we do not know its data
func sendEmailsToStudents_ex(Femail1, password1, email string, subject1 string) error {
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

	// Iterate through the data and send emails
	for _, row := range data {

		for colName, value := range row {
			// Convert byte value to text string
			valueStr := ""
			if byteValue, ok := value.([]byte); ok {
				valueStr = string(byteValue)
			}
			log.Printf("%s: %s", colName, valueStr)
		}

		// Replace bookmarks in mail template with row values
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
			log.Printf("Invalid data type in Email column: %T", row["Correo"])
			continue
		}

		emailValue := string(emailValueBytes)
		if emailValue == "" {
			log.Printf("Empty email value in row")
			continue
		}

		err := sendEmail(Femail1, password1, emailValue, subject1, htmlContent)
		if err != nil {
			log.Printf("Error sending mail to %s: %v", emailValue, err)
		} else {
			log.Printf("Email sent to %s", emailValue)
		}

	}

	return nil
}

// function to send office 365 emails
func sendEmail(Femail, password, email, subject, body string) error {
	// Configure sending emails
	dialer := gomail.NewDialer("smtp.office365.com", 587, Femail, password)

	// Create the custom email
	message := gomail.NewMessage()
	message.SetHeader("From", Femail)
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
