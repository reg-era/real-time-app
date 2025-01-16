import { GetData } from "../libs/post.js";
import BASE from "./_BASE.js";

export default class extends BASE {
    constructor(params) {
        super(params);

        this.statusError;
        this.statusMsg;
        this.errorMsg

        this.setTitle("Home");
        this.setStyle("http://localhost:8080/assets/css/base.css")
        this.setStyle("http://localhost:8080/assets/css/posts.css")
    }

    async getPosts() {
        const posts = await GetData()
        const container = document.body.querySelector('.posts')
        posts.forEach(post => container.appendChild(post))
    }

    async getHtml() {
        const html = `
        ${this.getHtmlBase()}
        <main>
            <section class="container">
                <div class="error-message">
                    <h1>Oops! ${this.statusMsg} (${this.statusError})</h1>
                    <p>${this.errorMsg}</p>
                    <button class="err-button" onclick="window.location.href='/'">Go to Home</button>
                </div>
            </section>
        </main>
        `
        return html
    }
}