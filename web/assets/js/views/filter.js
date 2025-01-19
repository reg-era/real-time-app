import { renderPage } from "../libs/post.js";
import { BASE } from "./_BASE.js";

export class Filter extends BASE {
    constructor(params) {
        super(params);
        this.setTitle("Filter");
    }

    async getFiltredData() {
        const { by } = this.params;

        if (by === 'categories') {
            this.setStyle("http://localhost:8080/api/css/categories.css");
            document.querySelector('.posts')?.remove();
            try {
                const res = await fetch('http://localhost:8080/api/new_post');
                if (res.ok) {
                    const categories = await res.json();
                    const html = categories.map(obj => `
                        <li categoryId="${obj.Id}">
                            <a href="/filter/${obj.Id}" class="nav__link" data-link>
                                <span class="category-name">${obj.name}:</span>
                                <span class="category-description">${obj.description}</span>
                            </a>
                        </li>
                    `).join('');

                    document.querySelector('main').innerHTML = `
                        <section class="categories">
                            <h2>Select a Category</h2>
                            <ul class="category-list">${html}</ul>
                        </section>`;
                } else {
                    window.location.href = `/error?status=${res.status}`;
                }
            } catch (err) {
                console.error(err);
            }
            return;
        }

        const fetchPosts = async (url, postIds) => {
            try {
                const response = await fetch(url);
                if (response.ok) {
                    const data = await response.json();
                    const posts = await renderPage(postIds || data.post_ids);
                    const container = document.body.querySelector('.posts');
                    posts.forEach(post => container.appendChild(post));
                } else {
                    window.location.href = `/error?status=${response.status}`;
                }
            } catch (err) {
                console.error(err);
            }
        };

        if (by === "created-posts" || by === "liked-posts") {
            await fetchPosts(`http://localhost:8080/api/me/${by}`);
        }
        else if (Number.parseInt(by)) {
            await fetchPosts(`http://localhost:8080/api/categories?category=${Number.parseInt(by)}`);
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