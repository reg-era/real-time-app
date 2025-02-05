import { GetData, renderPage } from "../libs/post.js";
import { debounce } from "../libs/script.js";
import { BASE } from "./_BASE.js";

export class Home extends BASE {
    constructor(app) {
        super();
        this.base = app;
        this.posts = null;
        this.listener = null;
        this.debouncedRenderPage = null;
        this.setTitle("Home");
        this.init();
    }

    async init() {
        this.posts = await GetData();
        this.debouncedRenderPage = debounce(this.getPosts, 1000);
    }

    async getPosts() {
        const posts = await renderPage(this.posts);
        const container = document.querySelector('.posts');
        posts.forEach(post => container.appendChild(post));
    }

    renderHtml() {
        return `
            ${this.base.getHtmlBase()}
            <main>
                ${this.base.getSidebar()}
                <section class="posts">
                </section>
            </main>
        `;
    }

    async listenerstoscroll() {
        console.log(this.debouncedRenderPage);

        if (this.listener === null) {
            this.listener = window.addEventListener('scroll', async () => {
                const scrollPosition = window.scrollY;
                const documentHeight = document.documentElement.scrollHeight;
                const windowHeight = window.innerHeight;
                if (scrollPosition + windowHeight >= documentHeight - 10) {
                    this.debouncedRenderPage();
                }
            });
        }
    }

    async afterRender() {
        await this.getPosts(this.posts);
        await this.listenerstoscroll();
        this.setupmssglistner(this.base);
        this.base.renderSidebar()
        this.setupAuthNav(this.base);
        this.base.setupNavigation(this.base);
        this.setupSidebar();
    }
}