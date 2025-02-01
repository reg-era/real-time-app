import { MessagesBase } from "./_MSG.js";

export class Messages extends MessagesBase {
    constructor(params) {
        super(params);
        this.params = params
    }

    async getMessages() {
        const conversation = ""
        try {
            const res = await fetch(`http://localhost:8080/api/messages?section=message&name=${this.params.user}`)
            const data = await res.json()

            for (let i = 0; i < data.length; i++) {
                const messageCompon = document.createElement('div')
                messageCompon.classList.add('message')
                console.log(data[i].sender_name);
                data[i].is_sender ? messageCompon.classList.add('receiver') : messageCompon.classList.add('sender');
                messageCompon.innerHTML = `<p>${data[i].message}</p>`
                document.querySelector('.messages-section').appendChild(messageCompon)
            }
        } catch (error) {
            console.error(error);
        }
        return conversation
    }

    setupConversation() {
        const send = document.querySelector('.message-input')
        send.addEventListener("keydown", async (event) => {
            if (event.key === "Enter" && !event.shiftKey) {
                const message = send.value.trim();
                if (message) {
                    try {
                        const res = await fetch(`http://localhost:8080/api/messages`, {
                            method: 'POST',
                            headers: { "Content-Type": "application/json" },
                            body: JSON.stringify({
                                to: this.params.user,
                                message: message
                            })
                        })
                        if (!res.ok) {
                            window.location.href = '/error?status=500';
                        } else {
                            const messageCompon = document.createElement('div');
                            messageCompon.classList.add('message', 'receiver');
                            messageCompon.innerHTML = `<p>${message}</p>`;
                            document.querySelector('.messages-section').appendChild(messageCompon);
                            send.value = '';
                        }
                    } catch (error) {
                        console.error(error);
                    }
                }
            }
        })
    }

    async getHtml() {
        const base = await super.getHtml();
        const html = `
        ${base}
        <main>
            <section class="conversation">
                <div class="messages-section">
                </div>
                <div class="messages-input">
                    <input required placeholder="Type message ..." class="message-input"></input>
                    <p class="error-comment"></p>
                </div>
            </section>
        </main>`;

        setTimeout(this.setupConversation.bind(this), 0)
        setTimeout(this.getMessages.bind(this), 0)
        return html;
    }
}