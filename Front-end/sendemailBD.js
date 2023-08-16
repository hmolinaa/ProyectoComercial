document.getElementById("sendEmails").addEventListener("click", function (event) {
    event.preventDefault();
    
    sendEmails();
  });

  function sendEmails() {
    const Femail = document.getElementById("Femail").value
    const password = document.getElementById("password").value
    const emailTemplate = document.getElementById("emailTemplate").value;
    const subject = document.getElementById("subject").value;
    const formData = new FormData();
    formData.append("Femail",Femail)
    formData.append("password",password)
    formData.append("emailTemplate", emailTemplate);
    formData.append("subject", subject);

    fetch("http://localhost:8080/send-emails", {
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