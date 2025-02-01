export class popup {
    constructor(app) {
        this.base = app;
    }

    async getMessages(name) {
        const fieldmssg = document.createElement('div');
        fieldmssg.innerHTML = `
        <div class="messages-input">
            <input required placeholder="Type message ..." class="message-input"></input>
            <p class="error-comment"></p>
        </div>`;
        const popupmessages = document.createElement('div');
        popupmessages.classList.add('messages-section');
        try {
            const res = await fetch(`http://localhost:8080/api/messages?section=message&name=${name}`)
            const data = await res.json()

            for (let i = 0; i < data.length; i++) {
                const messageCompon = document.createElement('div');
                messageCompon.classList.add('message');
                messageCompon.id = name;
                console.log(data[i].sender_name);
                data[i].IsSender ? messageCompon.classList.add('receiver') : messageCompon.classList.add('sender');
                messageCompon.innerHTML = `<p>${data[i].Message}</p>`;
                popupmessages.appendChild(messageCompon);
            }

            popupmessages.append(fieldmssg);
            document.body.appendChild(popupmessages);
        } catch (error) {
            console.error(error);
        }
    }

    setupConversation(name) {
        const send = document.querySelector('.message-input');
        send.addEventListener("keydown", async (event) => {
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
                        messageCompon.innerHTML = `<p>${message}</p>`;
                        document.querySelector('.messages-section').appendChild(messageCompon);
                        send.value = '';

                    } catch (error) {
                        console.error(error);
                    }
                }
            }
        })
    }
}