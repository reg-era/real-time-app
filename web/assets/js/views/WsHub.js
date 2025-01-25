import { BASE } from "./_BASE.js";

export class Messg extends BASE {
    constructor(params) {
        super(params);
        this.conn = new WebSocket('ws://localhost:8080/ws');
        this.conn.onopen = () => {
            console.log('Connected to WebSocket');
        };

        this.conn.onmessage = (event) => {
            const data = JSON.parse(event.data);
            console.log('Received:', data);
            // Handle message
        };

        this.conn.onerror = (error) => {
            console.error('WebSocket error:', error);
        };

        this.conn.onclose = () => {
            console.log('WebSocket closed');
        };
    }

    sendMessage(message) {
        if (this.conn.readyState === WebSocket.OPEN) {
            this.conn.send(JSON.stringify(message));
        }
    }
    setListners() {
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
                const response = await fetch("api/register", {
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
    }
    async getHtml() {
        let html = `
        ${await this.getHtmlBase()}
        <main>
            <div class="container">
                <section class="signup">
                    <h2>Create Account</h2>
                    <form id="signup-form">
                        <div class="form-group">
                            <label for="signup-username">Username:</label>
                            <input type="text" id="signup-username" name="username" placeholder="Enter your username"
                                minlength="5" maxlength="30" required>
                        </div>
                        <div class="form-group">
                            <label for="signup-email">Email:</label>
                            <input type="email" id="signup-email" name="email" placeholder="Enter your email" required>
                        </div>
                        <div class="form-group">
                            <label for="signup-password">Password:</label>
                            <input type="password" id="signup-password" name="password" placeholder="Enter your password"
                                minlength="8" maxlength="64" required>
                        </div>
                        <div class="form-group">
                            <label for="signup-confirm-password">Confirm Password:</label>
                            <input type="password" id="signup-confirm-password" name="confirm-password"
                                placeholder="Confirm your password" minlength="8" maxlength="64" required>
                        </div>
                        <button type="submit">Sign Up</button>
                        <p class="login-link">Already have an account? <a href="/login" data-link>Login</a></p>
                        <p id="responseMessage"></p>
                    </form>
                </section>
            </div>
        </main>
        <footer>
            <p>&copy; Regera, Yhajjaoui</p>
        </footer>
        `;

        //this.setListners();
        return html
    }
}