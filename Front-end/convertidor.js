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
                    document.getElementById('items').style.display = 'table';
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

    // Limpiar la tabla antes de agregar nuevos datos
    tableBody.innerHTML = '';

    // Iterar sobre los datos y agregar filas a la tabla
    data.forEach((item) => {
        const row = document.createElement('tr');
        row.innerHTML = `
    <td>${item.Nombre}</td>
    <td>${item.Cuenta}</td>
    <td>${item['Asignatura']}</td>
    <td>${item['Parcial 1']}</td>
    <td>${item['Parcial 2']}</td>
    <td>${item['Parcial 3']}</td>
    <td>${item['Nota Final']}</td>
    <td>${item.Correo}</td>
     `;
        tableBody.appendChild(row);
    });
}
