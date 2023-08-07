// Escuchador de evento para el formulario 'studentBD'
document.getElementById('studentBD').addEventListener('submit', function(event) {
    event.preventDefault();
    // Obtener los valores de los campos del formulario
    const subject = document.getElementById('subject').value;
    const customMessage = document.getElementById('customMessage').value;
    const professorName = document.getElementById('professorName').value;
    const c_name = document.getElementById('c_name').value;
    const c_subject = document.getElementById('c_subject').value;
    const desp = document.getElementById('desp').value;
    // Llamar a la función sendEmails con el endpoint 
    sendEmails('http://localhost:8080/send_emails/students', subject, customMessage, professorName, c_name, c_subject, desp);
});

// Escuchador de evento para el formulario 'studentEx'
document.getElementById('studentEx').addEventListener('submit', function(event) {
    event.preventDefault();
    // Obtener los valores de los campos del formulario
    const subject = document.getElementById('subject1').value;
    const customMessage = document.getElementById('customMessage1').value;
    const professorName = document.getElementById('professorName1').value;
    const c_name = document.getElementById('c_name1').value;
    const c_subject = document.getElementById('c_subject1').value;
    const desp = document.getElementById('desp1').value;
    // Llamar a la función sendEmails con el endpoint 
    sendEmails('http://localhost:8080/send_emails/students_excel', subject, customMessage, professorName, c_name, c_subject, desp);
});

// Función para enviar correos electrónicos
function sendEmails(endpoint, subject, customMessage, professorName, c_name, c_subject, desp) {
    fetch(endpoint, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        // Convertir los datos del formulario en una cadena JSON y enviarlos como cuerpo de la solicitud
        body: JSON.stringify({
            subject: subject,
            customMessage: customMessage,
            professorName: professorName,
            c_name: c_name,
            c_subject: c_subject,
            desp: desp
        })
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        // Analizar la respuesta JSON del servidor y mostrar un mensaje de alerta
        return response.json();
    })
    .then(data => {
        alert(data.message);
    })
    .catch(error => {
        console.error('Error sending emails:', error);
        alert('Error sending emails. Please try again.');
    });
}






