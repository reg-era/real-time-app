import { BASE } from "./_BASE.js";
import { handleResize } from "../libs/script.js";
export class Login extends BASE {
    constructor(app) {
        super();
        this.base = app;
        this.setTitle("Login");
        this.setStyle("/api/css/login.css");
    }

    setListeners() {

        document.getElementById("login-form").addEventListener("submit", async (event) => {
            event.preventDefault();
            const username = document.getElementById("login-username").value;
            const password = document.getElementById("login-password").value;
            const messageElement = document.getElementById("responseMessage");

            try {
                const response = await fetch("/api/login", {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json"
                    },
                    credentials: 'include',
                    body: JSON.stringify({ username, password }),
                });

                if (response.ok) {
                    this.base.loged = true;
                    history.pushState(null, null, "/");
                    await this.base.router.handleRoute();
                    await this.base.initializeWebSocket();
                } else {
                    const errorData = await response.text();
                    messageElement.textContent = errorData || 'Login failed. Please try again.';
                    messageElement.style.color = "red";
                }
            } catch (error) {
                console.error('Login error:', error);
                messageElement.textContent = "Unable to connect to the server. Please try again later.";
                messageElement.style.color = "red";
            }
        });
    }

    async renderHtml() {
        return `
        ${this.getHtmlBase()}
        <main>
            <div class="container">
                <section class="login">
                    <h2>Login</h2>
                    <form id="login-form">
                        <div class="form-group">
                            <label for="login-username">Username or Email:</label>
                            <input type="text" id="login-username" name="username" 
                                placeholder="Enter your username or Email"
                                minlength="5" maxlength="30" required>
                        </div>
                        <div class="form-group">
                            <label for="login-password">Password:</label>
                            <input type="password" id="login-password" name="password" 
                                placeholder="Enter your password"
                                minlength="8" maxlength="64" required>
                        </div>
                        <button type="submit">Login</button>
                        <p class="signup-link">Don't have an account? <a href="/register"  data-link>Signup</a></p>
                        <p id="responseMessage"></p>
                    </form>
                </section>
            </div>
        </main>
        <footer>
            <p>&copy; Regera, Yhajjaoui</p>
        </footer>
        `;
    }

    afterRender() {
        this.setListeners();
        this.setupNavigation(this.base);
        this.setupSidebar();
        handleResize();

    }
}