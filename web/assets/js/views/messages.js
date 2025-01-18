import { MessagesBase } from "./_MSG.js";

export class Messages extends MessagesBase {
    constructor(params) {
        super(params);
    }

    getMessages(){
        const conversation = `
        <div class="message user">
            <p>1111111111111111111111111?</p>
        </div>
        <div class="message sender">
            <p>I'm doing great, thank you! How about you?</p>
        </div>
        <div class="message user">
            <p>I'm doing well, just a little busy with work.</p>
        </div>
        <div class="message sender">
            <p>That sounds like a lot! Hopefully, you can take a break soon.</p>
        </div>        <div class="message user">
            <p>Hi! How are you today?</p>
        </div>
        <div class="message sender">
            <p>I'm doing great, thank you! How about you?</p>
        </div>
        <div class="message user">
            <p>I'm doing well, just a little busy with work.</p>
        </div>
        `
        return conversation
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
        return html;
    }
}