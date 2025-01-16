import { GetData } from "../libs/post.js";
import BASE from "./_BASE.js";

export default class extends BASE {
    constructor(params) {
        super(params);
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
            ${this.getSideBar()}
            <section class="posts">
            </section>
        <main>
        `

        setTimeout(this.getPosts, 0)
        return html
    }
}