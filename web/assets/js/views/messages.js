import { MessagesBase } from "./_MSG.js";

export class Messages extends MessagesBase {
    constructor(params) {
        super(params);
        this.params = params
    }

    getMessages() {
        const conversation = ""
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
                        }
                    } catch (error) {
                        console.error(error);
                    }
                }
            }
        })
    }

    async getHtml() {
        console.log(this.params);
        const base = await super.getHtml();
        const html = `
        ${base}
        <main>
            <section class="conversation">
                <div class"messages-section">
                    ${this.getMessages()}
                </div>
                <div class="messages-input">
                    <input required placeholder="Type message ..." class="message-input"></input>
                    <p class="error-comment"></p>
                </div>
            </section>
        </main>
        `;
        setTimeout(this.setupConversation.bind(this), 0)
        return html;
    }
}