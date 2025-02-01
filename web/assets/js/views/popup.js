export class popup {
    constructor(app) {
        this.base = app;
    }

    async getMessages(name) {
        const conversation = ""
        const popupmessages = document.createElement('div');
        popupmessages.classList.add('messages-section');
        try {
            const res = await fetch(`http://localhost:8080/api/messages?section=message&name=${name}`)
            const data = await res.json()
            console.log(data);

            for (let i = 0; i < data.length; i++) {
                const messageCompon = document.createElement('div');
                messageCompon.classList.add('message');
                console.log(data[i].sender_name);
                data[i].IsSender ? messageCompon.classList.add('receiver') : messageCompon.classList.add('sender');
                messageCompon.innerHTML = `<p>${data[i].Message}</p>`;
                popupmessages.appendChild(messageCompon);
                document.body.appendChild(popupmessages);
            }
        } catch (error) {
            console.error(error);
        }
        return conversation
    }
}