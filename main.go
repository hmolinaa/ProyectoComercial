package main

import (
	//"fmt"
	"database/sql"
	"fmt"

	//"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

func conexionBD() (conexion *sql.DB) {
	Driver := "mysql"
	Usuario := "root"
	Contrasenia := ""
	Nombre := "unah"

	conexion, err := sql.Open(Driver, Usuario+":"+Contrasenia+"@tcp(127.0.0.1)/"+Nombre)
	if err != nil {
		panic(err.Error())
	}
	return conexion

}

var plantillas = template.Must(template.ParseGlob("plantillas/*"))

func main() {
	http.HandleFunc("/", Principal)
	http.HandleFunc("/inicio", Inicio)
	http.HandleFunc("/crear", Crear)
	http.HandleFunc("/insertar", Insertar)
	http.HandleFunc("/eliminar", Eliminar)
	http.HandleFunc("/editar", Editar)
	http.HandleFunc("/actualizar", Actualizar)
	http.HandleFunc("/email", Email)

	fmt.Println("Hola, este servidor esta corriendo...")

	http.ListenAndServe(":8080", nil)

}

func Principal(w http.ResponseWriter, r *http.Request) {
	plantillas.ExecuteTemplate(w, "principal", nil)

}

func Eliminar(w http.ResponseWriter, r *http.Request) {
	idEstudiante := r.URL.Query().Get("id")
	//fmt.Println(idEstudiante)

	conexionEstablecida := conexionBD()
	eliminarRegistro, err := conexionEstablecida.Prepare("DELETE FROM estudiantes WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	eliminarRegistro.Exec(idEstudiante)
	http.Redirect(w, r, "/", 301)

}

type Estudiante struct {
	Id         int
	Nombre     string
	Cuenta     int
	Asignatura string
	Parcial_1  int
	Parcial_2  int
	Parcial_3  int
	NotaFinal  int
	Correo     string
}

func Inicio(w http.ResponseWriter, r *http.Request) {

	conexionEstablecida := conexionBD()
	registros, err := conexionEstablecida.Query("SELECT * FROM estudiantes")

	if err != nil {
		panic(err.Error())
	}

	estudiante := Estudiante{}
	arregloEstudiante := []Estudiante{}

	for registros.Next() {
		var id, cuenta, parcial1, parcial2, parcial3, notafinal int
		var nombre, asignatura, correo string
		err = registros.Scan(&id, &nombre, &cuenta, &asignatura, &parcial1, &parcial2, &parcial3, &notafinal, &correo)
		if err != nil {
			panic(err.Error())
		}
		estudiante.Id = id
		estudiante.Nombre = nombre
		estudiante.Cuenta = cuenta
		estudiante.Asignatura = asignatura
		estudiante.Parcial_1 = parcial1
		estudiante.Parcial_2 = parcial2
		estudiante.Parcial_3 = parcial3
		estudiante.NotaFinal = notafinal
		estudiante.Correo = correo

		arregloEstudiante = append(arregloEstudiante, estudiante)

	}
	//fmt.Println(arregloEstudiante)

	plantillas.ExecuteTemplate(w, "inicio", arregloEstudiante)
}

func Editar(w http.ResponseWriter, r *http.Request) {
	idEstudiante := r.URL.Query().Get("id")
	//fmt.Println(idEstudiante)

	conexionEstablecida := conexionBD()
	registro, err := conexionEstablecida.Query("SELECT * FROM estudiantes WHERE id=?", idEstudiante)

	estudiante := Estudiante{}
	for registro.Next() {
		var id, cuenta, parcial1, parcial2, parcial3, notafinal int
		var nombre, asignatura, correo string
		err = registro.Scan(&id, &nombre, &cuenta, &asignatura, &parcial1, &parcial2, &parcial3, &notafinal, &correo)
		if err != nil {
			panic(err.Error())
		}
		estudiante.Id = id
		estudiante.Nombre = nombre
		estudiante.Cuenta = cuenta
		estudiante.Asignatura = asignatura
		estudiante.Parcial_1 = parcial1
		estudiante.Parcial_2 = parcial2
		estudiante.Parcial_3 = parcial3
		estudiante.NotaFinal = notafinal
		estudiante.Correo = correo

	}

	//fmt.Println(estudiante)
	plantillas.ExecuteTemplate(w, "editar", estudiante)

}

func Email(w http.ResponseWriter, r *http.Request) {

	conexionEstablecida := conexionBD()
	registros, err := conexionEstablecida.Query("SELECT * FROM estudiantes")

	if err != nil {
		panic(err.Error())
	}

	estudiante := Estudiante{}
	AEstudiante := []Estudiante{}

	for registros.Next() {
		var id, cuenta, parcial1, parcial2, parcial3, notafinal int
		var nombre, asignatura, correo string
		err = registros.Scan(&id, &nombre, &cuenta, &asignatura, &parcial1, &parcial2, &parcial3, &notafinal, &correo)
		if err != nil {
			panic(err.Error())
		}
		estudiante.Id = id
		estudiante.Nombre = nombre
		estudiante.Cuenta = cuenta
		estudiante.Asignatura = asignatura
		estudiante.Parcial_1 = parcial1
		estudiante.Parcial_2 = parcial2
		estudiante.Parcial_3 = parcial3
		estudiante.NotaFinal = notafinal
		estudiante.Correo = correo

		AEstudiante = append(AEstudiante, estudiante)

	}
	//fmt.Println(AEstudiante)

	plantillas.ExecuteTemplate(w, "email", AEstudiante)

}

func Crear(w http.ResponseWriter, r *http.Request) {
	plantillas.ExecuteTemplate(w, "crear", nil)
}

func Insertar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		nombre := r.FormValue("nombre")
		cuenta := r.FormValue("cuenta")
		asignatura := r.FormValue("asignatura")
		parcial := r.FormValue("parcial1")
		parcial2 := r.FormValue("parcial2")
		parcial3 := r.FormValue("parcial3")
		notafinal := r.FormValue("notafinal")
		correo := r.FormValue("correo")

		conexionEstablecida := conexionBD()
		insertarRegistros, err := conexionEstablecida.Prepare("INSERT INTO estudiantes( nombre, cuenta, asignatura, parcial, parcial2, parcial3, notafinal, correo) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			panic(err.Error())
		}
		insertarRegistros.Exec(nombre, cuenta, asignatura, parcial, parcial2, parcial3, notafinal, correo)
		http.Redirect(w, r, "/", 301)

	}

}

func Actualizar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		nombre := r.FormValue("nombre")
		cuenta := r.FormValue("cuenta")
		asignatura := r.FormValue("asignatura")
		parcial := r.FormValue("parcial1")
		parcial2 := r.FormValue("parcial2")
		parcial3 := r.FormValue("parcial3")
		notafinal := r.FormValue("notafinal")
		correo := r.FormValue("correo")

		conexionEstablecida := conexionBD()
		modificarRegistros, err := conexionEstablecida.Prepare("UPDATE estudiantes SET  nombre = ?, cuenta = ?, asignatura = ?, parcial=?,parcial2=?,parcial3=?,notafinal=?,correo=? WHERE id=? ")
		if err != nil {
			panic(err.Error())
		}
		modificarRegistros.Exec(nombre, cuenta, asignatura, parcial, parcial2, parcial3, notafinal, correo, id)
		http.Redirect(w, r, "/", 301)

	}

}
