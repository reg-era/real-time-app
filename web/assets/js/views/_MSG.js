import { BASE } from "./_BASE.js";

export class MessagesBase extends BASE {
    constructor(params) {
        super(params);
        this.setStyle("http://localhost:8080/api/css/messages.css");
    }


    setupSearch() {
        const input = document.querySelector('.users-input')
        const search = document.querySelector('.search-user')

        search.addEventListener('click', async (e) => {
            try {
                if (input.value.length == 0) return
                const allNavUsers = document.querySelectorAll('.sidebar-nav .nav__link')
                for (let i = 0; i < allNavUsers.length; i++) {
                    if (allNavUsers[i]?.id == input.value) {
                        return
                    }
                }
                const res = await fetch(`http://localhost:8080/api/messages?section=user&name=${input.value}`)
                if (res.ok) {
                    const sidbare = document.querySelector('.sidebar-nav')
                    const newConv = document.createElement('a');
                    newConv.href = `/messages/${input.value}`;
                    newConv.classList.add('nav__link');
                    newConv.setAttribute('data-link', '');
                    newConv.setAttribute('id', input.value)
                    newConv.innerHTML = `👤  ${input.value}`;
                    sidbare.prepend(newConv);
                } else if (res.status == 404) {
                    console.log("user not found");
                }
            } catch (error) {
                console.error(error);
            }
        })
    }

    async getPrevConversation() {
        try {
            const res = await fetch(`http://localhost:8080/api/messages?section=user`)
            const data = await res.json()
            return data.friends
        } catch (error) {
            console.error(error);
        }
    }

    async getSideBar() {
        let conversation = `<input type="text" class="users-input" placeholder="type users name...">
        <button class="search-user">search</button>`

        const prevConvers = await this.getPrevConversation();
        prevConvers.forEach(user => {
            conversation += `<a href="/messages/${user}" class="nav__link" id="${user}" data-link >👤  ${user}</a>`
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
        ${await this.getSideBar()}
        `
        setTimeout(this.setupSearch, 0)
        setTimeout(super.setListners, 0)
        return html
    }
}