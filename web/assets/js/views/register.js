// import { app } from "../main.js";
import { BASE } from "./_BASE.js";

export class Register extends BASE {
    constructor(app) {
        super();
        this.base = app;
        this.setTitle("Register");
        this.setStyle("http://localhost:8080/api/css/base.css")
        this.setStyle("http://localhost:8080/api/css/register.css")
    }

    setListners() {
        const self = this;
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
                    self.base.loged = true;

                    history.pushState(null, null, "/");
                    self.base.router.handleRoute();
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

    async renderHtml() {
        const html = `
        ${this.getHtmlBase()}
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
        `

        return html
    }

    afterRender() {
        this.setupAuthNav(this.base);
        this.setupSidebar();
        // this.setupNavigation();
        this.setListners();
    }
}

function validateSignup(password, confirmPassword, messageElement) {
    if (password !== confirmPassword) {
        messageElement.textContent = "Passwords do not match.";
        messageElement.style.color = "red";
        return false;
    }
    return true;
}