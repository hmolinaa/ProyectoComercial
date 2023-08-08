// Declarar la variable contentDiv
const contentDiv = document.getElementById('contentDiv');

// Función para mostrar el contenido del div
function mostrarContenidoDiv() {
    contentDiv.style.display = 'block';
}

async function enviarArchivo() {
    const fileInput = document.getElementById('fileInput');
    const file = fileInput.files[0];


    if (file) {
        const reader = new FileReader();

        reader.onload = async (event) => {
            try {
                const jsonData = await procesarArchivoExcel(event.target.result);

                console.log('JSON generado en el frontend:', jsonData); // Imprimir el JSON en la consola

                // Convertir el JSON a una cadena
                const jsonString = JSON.stringify(jsonData);


                const requestOptions = {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: jsonString,
                };

                // Realizar la solicitud Fetch POST
                const response = await fetch('http://localhost:8080/excel', requestOptions);

                if (response.ok) {
                    const data = await response.json();
                    mostrarTabla(data);
                    mostrarContenidoDiv();
                    // Obtén una referencia al contenedor de la tabla
                    const tablaContainer = document.getElementById('tabla-container');

                    // Crea el elemento de tabla
                    const tabla = document.createElement('table');
                    tabla.id = 'items';

                    // Inserta la tabla dentro del contenedor
                    tablaContainer.appendChild(tabla);

                    mostrarContenidoDiv();
                } else {
                    console.log('Error en la solicitud Fetch:', response.status, response.statusText);
                }
            } catch (error) {
                console.log('Error al procesar el archivo Excel:', error.message);
            }
        };

        reader.readAsBinaryString(file);
    } else {
        console.log('Por favor, selecciona un archivo.');
    }


}

async function procesarArchivoExcel(fileContent) {
    return new Promise((resolve, reject) => {
        const workbook = XLSX.read(fileContent, { type: 'binary' });

        // Supongamos que los datos se encuentran en la primera hoja (Sheet) del archivo Excel
        const sheetName = workbook.SheetNames[0];
        const worksheet = workbook.Sheets[sheetName];

        // Convertir el archivo Excel a un arreglo de objetos JSON
        const jsonData = XLSX.utils.sheet_to_json(worksheet);

        resolve(jsonData);
    });
}


function mostrarTabla(data) {
    const tableBody = document.querySelector('#items tbody');
    const tableHeader = document.querySelector('#items thead');

    // Limpiar la tabla antes de agregar nuevos datos
    tableBody.innerHTML = '';
    tableHeader.innerHTML = '';

    if (data.length === 0) {
        return; // No hay datos para mostrar
    }

    // Obtener los nombres de columna del primer elemento del JSON
    const columnNames = Object.keys(data[0]);

    // Crear la fila de encabezado de la tabla
    const headerRow = document.createElement('tr');
    columnNames.forEach(columnName => {
        const th = document.createElement('th');
        th.textContent = columnName;
        headerRow.appendChild(th);
    });
    tableHeader.appendChild(headerRow);

    // Iterar sobre los datos y agregar filas a la tabla
    data.forEach(item => {
        const row = document.createElement('tr');
        columnNames.forEach(columnName => {
            const td = document.createElement('td');
            td.textContent = item[columnName];
            row.appendChild(td);
        });
        tableBody.appendChild(row);
    });
}

