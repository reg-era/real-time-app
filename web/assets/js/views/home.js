import { GetData } from "../libs/post.js";
// import { app } from "../main.js";
import { BASE } from "./_BASE.js";

export class Home extends BASE {
    constructor(app) {
        super();
        this.base = app;
        this.setTitle("Home");
    }

    async getPosts() {
        const posts = await GetData();
        const container = document.querySelector('.posts');
        posts.forEach(post => container.appendChild(post));
    }

    async renderHtml() {
        return `
            ${this.getHtmlBase()}
            <main>
                ${this.getSidebar()}
                <section class="posts">
                </section>
            </main>
            ${this.getPosts()}
        `;
    }

    afterRender() {
        this.setupAuthNav(this.base);
        this.setupNavigation(this.base);
        this.setupSidebar();
    }
}