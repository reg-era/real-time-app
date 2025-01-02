document.getElementById("login-form").addEventListener("submit", async function (event) {
    event.preventDefault();
    // Capture form data
    const username = document.getElementById("login-username").value;
    const password = document.getElementById("login-password").value;
    const messageElement = document.getElementById("responseMessage");
    // Send data to the API
    try {
        const response = await fetch("http://localhost:8080/login", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ username, password }),
        });

        if (response.ok) {
            window.location.href = '/';
        } else {
            const errorData = await response.text();
            messageElement.textContent = `Error: ${errorData}`;
            messageElement.style.color = "red";
        }
    } catch {
        messageElement.textContent = "An error occurred during registration.";
        messageElement.style.color = "red";
    }
});
