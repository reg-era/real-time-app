import { GetData } from "../libs/post.js";
import { BASE } from "./_BASE.js";

export class Filter extends BASE {
    constructor(params) {
        super(params);
        this.setTitle("Filter");
    }

    async getHtml() {
        console.log(this.params);
        
        const html = `
        ${this.getHtmlBase()}
        <main>     
            ${this.getSideBar()}
            <section class="posts">
            </section>
        <main>
        `
        return html
    }
}