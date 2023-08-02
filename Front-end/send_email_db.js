document.getElementById("sendEmails").addEventListener("click", function (event) {
    event.preventDefault(); 
    sendEmails();
});

function sendEmails() {
    fetch("http://localhost:8080/send-emails")
        .then(response => {
            if (!response.ok) {
                throw new Error("Error al enviar correos");
            }
            return response.text();
        })
        .then(data => {
            console.log(data);
            showSuccessAlert();
        })
        .catch(error => {
            console.error("Error: ", error);
            showErrorAlert();
        });
}

function showSuccessAlert() {
    window.alert("Correos enviados correctamente");
}

function showErrorAlert() {
    window.alert("Error al enviar correos");
}
