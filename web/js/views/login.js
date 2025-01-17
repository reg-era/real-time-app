import { BASE } from "./_BASE.js";

export class Login extends BASE {
    constructor(params) {
        super(params);
        this.setTitle("Login");
        this.setStyle("http://localhost:8080/assets/css/base.css")
        this.setStyle("http://localhost:8080/assets/css/login.css")
    }

    setAttribute() { }

    setListners() {
        document.getElementById("login-form").addEventListener("submit", async function (event) {
            event.preventDefault();
            const username = document.getElementById("login-username").value;
            const password = document.getElementById("login-password").value;
            const messageElement = document.getElementById("responseMessage");

            try {
                const response = await fetch("http://localhost:8080/api/login", {
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
    }

    async getHtml() {
        const html = `
        ${this.getHtmlBase()}
        <main>
            <div class="container">
                <section class="login">
                    <h2>Login</h2>
                    <form id="login-form">
                        <div class="form-group">
                            <label for="login-username">Username:</label>
                            <input type="text" id="login-username" name="username" placeholder="Enter your username"
                                minlength="5" maxlength="30" required>
                        </div>
                        <div class="form-group">
                            <label for="login-password">Password:</label>
                            <input type="password" id="login-password" name="password" placeholder="Enter your password"
                                minlength="8" maxlength="64" required>
                        </div>
                        <button type="submit">Login</button>
                        <p class="signup-link">Don't have an account? <a href="/register" data-link>Signup</a></p>
                        <p id="responseMessage"></p>
                    </form>
                </section>
            </div>
        </main>
        <footer>
            <p>&copy Regera, Yhajjaoui</p>
        </footer>
        `

        setTimeout(this.setListners, 0)
        return html
    }
}