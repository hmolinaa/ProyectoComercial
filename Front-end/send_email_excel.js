document.getElementById("studentEx").addEventListener("submit", function (event) {
    event.preventDefault();
    sendEmails();
});

function sendEmails() {
    const email = document.getElementById("email").value;
    const subject1 = document.getElementById("subject1").value;
    const formData = new FormData();
    formData.append("email", email);
    formData.append("subject1", subject1);

    fetch("http://localhost:8080/send-emails_ex", {
        method: "POST",
        body: formData,
    })
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






