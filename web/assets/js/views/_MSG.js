import { BASE } from "./_BASE.js";

export class MessagesBase extends BASE {
    constructor(params) {
        super(params);
        this.setStyle("http://localhost:8080/api/css/messages.css");
    }

    async getSideBar() {
        try {
            const res = await fetch(`http://localhost:8080/api/messages?section=user`)
            const data = await res.json()
            if (!res.ok) {
                window.location.href = `/error?status=${res.status}`;
            }
            let conversation = ""

            data.forEach(user => {
                conversation += `<a href="/messages/${user}" class="nav__link" id="${user}" data-link >
                👤  ${user} 
                <span class="status-tag ${user ? 'online' : 'offline'}"> ${user ? 'Online' : 'Offline'}</span>
                </a>`
            })

            return `
            <aside class="sidebar">
                <nav class="sidebar-nav">
                ${conversation}
                </nav>
            </aside>
            `
        } catch (error) {
            console.error(error);
        }

    }

    async getHtml() {
        const html = `
        ${this.getHtmlBase()}
        ${await this.getSideBar()}
        `
        setTimeout(this.setListners, 0)
        return html
    }
}