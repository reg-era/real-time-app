import { BASE } from "./_BASE.js";

export class MessagesBase extends BASE {
    constructor(params) {
        super(params);
        this.setStyle("http://localhost:8080/api/css/messages.css");
    }

    getConversation() {
        return [
            { user: 'ilyass' },
            { user: 'hasssan' },
            { user: 'lmolabi' },
            { user: '3daysa' },
        ]
    }

    getSideBar() {
        let conversation = ''
        this.getConversation().forEach(user => {
            conversation += `<a href="/messages/${user.user}" class="nav__link" data-link >👤  ${user.user}</a>`
        })
        return `
        <aside class="sidebar">
            <nav class="sidebar-nav">
            ${conversation}
            </nav>
        </aside>
        `
    }

    async getHtml() {
        const html = `
        ${this.getHtmlBase()}
        ${this.getSideBar()}
        `
        return html
    }
}