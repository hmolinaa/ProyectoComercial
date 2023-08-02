fetch('http://localhost:8080/inicio')
.then(response => response.json())
.then(data => {
  const tableBody = document.querySelector('#items-table tbody');
  data.forEach(item => {
    const row = document.createElement('tr');
    row.innerHTML = `
              <td>${item.Id}</td>
              <td>${item.Name}</td>
              <td>${item.Account}</td>
              <td>${item.Subject}</td>
              <td>${item.First_partial}</td>
              <td>${item.Second_partial}</td>
              <td>${item.Third_partial}</td>
              <td>${item.Final_score}</td>
              <td>${item.Email}</td>
          `;
    tableBody.appendChild(row);
  });
})
.catch(error => {
  console.error('Error:', error);
});
