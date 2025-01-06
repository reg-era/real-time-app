import BASE from "./_BASE.js";

export default class extends BASE {
    constructor(params) {
        super(params);
        this.setStyle("http://localhost:8080/assets/css/register.css")
        this.setTitle("Home");
    }

    async getHtml() {
        return `
        ${this.getHtmlBase()}
        ${this.getSideBar()}
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
    }
}