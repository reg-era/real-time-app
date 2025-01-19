import { renderPage } from "../libs/post.js";
import { BASE } from "./_BASE.js";

export class Filter extends BASE {
    constructor(params) {
        super(params);
        this.setTitle("Filter");
    }

    async getFiltredData() {
        if (this.params.by != "created-posts" && this.params.by != "liked-posts") return
        try {
            const response = await fetch(`http://localhost:8080/api/me/${this.params.by}`);
            if (response.ok) {
                const data = await response.json()
                const posts = await renderPage(data.post_ids)
                const container = document.body.querySelector('.posts')
                posts.forEach(post => container.appendChild(post))
            } else {
                window.location.href = `http://localhost:8080/error?status=${response.status}`;
            }
        } catch (err) {
            console.error(err);
        }
    }

    async getHtml() {
        const html = `
        ${this.getHtmlBase()}
        <main>     
            ${this.getSideBar()}
            <section class="posts">
            </section>
        </main>
        `
        setTimeout(() => this.getFiltredData(), 0);
        return html
    }
}