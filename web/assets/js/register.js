document.getElementById("signup-form").addEventListener("submit", async function (event) {
    event.preventDefault();

    const username = document.getElementById("signup-username").value;
    const email = document.getElementById("signup-email").value;
    const password = document.getElementById("signup-password").value;
    const confirmPassword = document.getElementById("signup-confirm-password").value;
    const messageElement = document.getElementById("responseMessage");

    if (!validateSignup(password, confirmPassword, messageElement)) {
        return;
    }

    try {
        const response = await fetch("http://localhost:8080/register", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ username, email, password, confirmPassword }),
        });

        if (response.ok) {
            messageElement.textContent = "Registration successful!";
            messageElement.style.color = "green";
            window.location.href = '/';
        } else {
            const errorData = await response.text();
            messageElement.textContent = `Error: ${errorData}`;
            messageElement.style.color = "red";
        }
    } catch (error) {
        messageElement.textContent = "An error occurred during registration.";
        messageElement.style.color = "red";
    }
});

function validateSignup(password, confirmPassword, messageElement) {
    if (password !== confirmPassword) {
        messageElement.textContent = "Passwords do not match.";
        messageElement.style.color = "red";
        return false;
    }
    return true;
}
