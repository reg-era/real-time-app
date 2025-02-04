export class popup {
    constructor(app) {
        this.base = app;
    }

    async getMessages(name) {
        const popMessage = document.createElement('div')
        popMessage.classList.add('conversation');

        const inputMessage = document.createElement('div')
        inputMessage.classList.add('messages-input');
        inputMessage.innerHTML = `
            <input required placeholder="Type message ..." class="message-input"></input>
            <p class="error-comment"></p>`;

        try {
            const res = await fetch(`http://localhost:8080/api/messages?section=message&name=${name}`);
            // if (!res.ok) {
            //     const pageerror = new Error('500', this.base)
            //     const html = await errorView.renderHtml();
            //     const appElement = document.querySelector('.app');
            //     appElement.innerHTML = html;
            //     appElement.setAttribute('page', 'error');
            //     if (typeof errorView.afterRender === 'function') {
            //         errorView.afterRender();
            //     }
            //     return
            // }
            const data = await res.json()

            const allMessages = document.createElement('div')
            allMessages.classList.add('messages-section');
            if (data) {
                for (let i = 0; i < data.length; i++) {
                    const messageCompon = document.createElement('div');
                    messageCompon.classList.add('message');
                    messageCompon.id = name;
                    data[i].IsSender ? messageCompon.classList.add('receiver') : messageCompon.classList.add('sender');
                    messageCompon.innerHTML = `
                    <div class="message-header">
                        <span class="username-message">${data[i].IsSender ? data[i].sender_name : name}</span>
                        <span class="timestamp">${new Date(data[i].CreatedAt).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}</span>
                    </div>
                    <p>${data[i].Message}</p>`
                    allMessages.appendChild(messageCompon);
                }
            }

            popMessage.append(allMessages, inputMessage);

            const over = document.createElement('div')
            over.classList.add('over-layer');

            document.body.append(over, popMessage);
            document.body.classList.add('has-overlay');

            allMessages.scrollTop = allMessages.scrollHeight;
            over.addEventListener('click', (e) => {
                popMessage.remove()
                over.remove()
                document.body.classList.remove('has-overlay');

                const notification = document.querySelector(`#${name} .notification`)
                notification.classList.add('hide')
                const counter = notification.querySelector('.notification-counter')
                counter.textContent = 0
            })
        } catch (error) {
            console.error(error);
        }
    }

    setupConversation(name) {
        const allMessages = document.querySelector('.messages-section');
        const username = document.get
        const send = document.querySelector('.message-input');
        document.addEventListener("keydown", async (event) => {
            if (event.key === "Enter" && !event.shiftKey) {
                const message = send.value.trim();
                if (message) {
                    try {
                        this.base.connection.send(JSON.stringify({
                            ReceiverName: name,
                            Data: message,
                        }));
                        const messageCompon = document.createElement('div');
                        messageCompon.classList.add('message', 'receiver');
                        // handle the name of loged user !!!!
                        messageCompon.innerHTML = messageCompon.innerHTML = `
                        <div class="message-header">
                          <span class="username-message">${'to handel'}</span>
                          <span class="timestamp">${new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}</span>
                        </div>
                        <p>${message}</p>`
                        document.querySelector('.messages-section').appendChild(messageCompon);
                        send.value = '';
                        allMessages.scrollTop = allMessages.scrollHeight;
                    } catch (error) {
                        console.error(error);
                    }
                }
            }
        })
    }
}