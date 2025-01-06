import { GetData } from "../libs/post.js";
import BASE from "./_BASE.js";

export default class extends BASE {
    constructor(params) {
        super(params);
        this.setTitle("Home");
        this.setStyle("http://localhost:8080/assets/css/base.css")
        this.setStyle("http://localhost:8080/assets/css/posts.css")

        this.allPosts = []
    }

    async getPosts() {
        const posts = await GetData()
        const container = document.body.querySelector('.posts')
        posts.forEach(post => container.appendChild(post))
    }

    setupEvents() {
        console.log(document.body.querySelectorAll('.post-container'));
    }

    async getHtml() {
        const html = `
        ${this.getHtmlBase()}
        <main>     
            ${this.getNavigation()}
            <section class="posts">
            </section>
        <main>
        `

        setTimeout(this.getPosts, 0)
        return html
    }
}