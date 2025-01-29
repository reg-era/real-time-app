import { GetData } from "../libs/post.js";
import { BASE } from "./_BASE.js";

export class Home extends BASE {
    constructor(params) {
        super(params);
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
        this.setupAuthNav();
        //this.setupNavigation();
        this.setupSidebar();
    }
}