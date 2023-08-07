# Sistema de Gestión de Calificaciones

Este proyecto es un sistema de envío de correos para estudiantes basado en Go, MySQL y Excel. Permite visualizar calificaciones almacenadas en una base de datos y Excel.

## Funcionalidades
- Visualizar calificaciones almacenadas en la base de datos.
- Envío de correos electrónicos personalizados a estudiantes con información relevante sobre sus calificaciones.
- Procesamiento de datos en formato JSON para su almacenamiento en una base de datos MySQL.
- Uso de la librería `gomail` para enviar correos electrónicos desde una cuenta de Gmail.

## Uso de Excel

1. Abre `inicio.html` en el navegador web para acceder a la interfaz del sistema.
2. Selecciona un archivo Excel y haz clic en "Enviar Archivo" para cargar las calificaciones.
3. Utiliza la tabla para ver las calificaciones almacenadas en la base de datos.
4. Llena el formulario de envío de correos para enviar notificaciones personalizadas a los estudiantes.

## Tecnologia usada en el Frontend

- HTML5
- Bootstrap 4
- XLSX.js (para manipular archivos Excel)
- JavaScript
## Requisitos en Go

- Go (Golang)
- MySQL
- Librería `github.com/go-gomail/gomail`
- Librería `github.com/gorilla/mux`
- Librería `github.com/rs/cors`

## Dependencias de Go
- go get -u `github.com/go-gomail/gomail`
- go get -u  `github.com/gorilla/mux`
- go get -u `github.com/rs/cors`


## Contenido en Go

- `inicio.html`: Página principal que contiene la interfaz de usuario.
- `convertidor.js`: Script que convierte los datos del archivo Excel en formato JSON.
- `inicio.js`: Script que carga y muestra las calificaciones de la base de datos en la interfaz.
- `send_email_excel.js`: Script que envía correos electrónicos personalizados a estudiantes con calificaciones de un archivo Excel y Mysql.



## Créditos

Este proyecto fue creado por Henrry Molina como parte de la clase Programacion Comercial. Puedes encontrar más detalles y contactarme en [mi perfil de GitHub](https://github.com/hmolinaa).



![](https://wallpaperaccess.com/full/1262277.jpg)

> UNAH CU

                    
$$\sin(\alpha)^{\theta}=\sum_{i=0}^{n}(x^i + \cos(f))$$
                


"Las matemáticas son el lenguaje son el idioma que uso Dios para escribir el mundo": Galileo Galilei
