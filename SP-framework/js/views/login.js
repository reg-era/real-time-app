import BASE from "./_BASE.js";

export default class extends BASE {
    constructor(params) {
        super(params);
        this.setTitle("Home");
        this.setStyle("http://localhost:8080/assets/css/login.css")
    }

    setListners(){

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

        setTimeout(this.setListners,0)
        return html
    }
}