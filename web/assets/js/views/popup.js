import { escapeHTML } from "../libs/post.js";
import { validCookies } from "../main.js";

export class popup {
    constructor(app) {
        this.base = app;
    }

    async getMessages(name) {
        const res = await validCookies()
        if (!(res.valid)) {
            history.pushState(null, null, '/login');
            this.base.router.handleRoute()
            return
        }

        const popMessage = document.createElement('div')
        popMessage.classList.add('conversation');
        popMessage.setAttribute('name', name);


        const inputMessage = document.createElement('div')
        inputMessage.classList.add('messages-input');
        inputMessage.innerHTML = `
            <input required placeholder="Type message ..." class="message-input1"></input>
            <p class="error-message"></p>`;

        try {
            const res = await fetch(`/api/messages?section=message&name=${name}`);
            if (!res.ok) {
                await this.base.router.handleError('500');
                return
            }
            const data = await res.json()

            const allMessages = document.createElement('div')
            allMessages.classList.add('messages-section');
            if (data) {
                const BATCH_SIZE = 10
                let nextBatch = []

                const throttle = (func, delay) => {
                    let prev = 0

                    return () => {
                        const now = new Date().getTime()
                        if (now - prev >= delay) {
                            prev = now
                            func()
                        }
                    }
                }

                const loadMoreMessages = () => {
                    if (data.length > 0) {
                        nextBatch = data.splice(-BATCH_SIZE)
                        insertMessages(nextBatch, allMessages, name)
                    }
                }

                allMessages.addEventListener('scroll', throttle(() => {
                    if (allMessages.scrollTop < (allMessages.scrollHeight - allMessages.clientHeight) / 4) {
                        loadMoreMessages()
                    }
                }, 200))

                loadMoreMessages()
            }
            popMessage.append(allMessages, inputMessage);

            const over = document.createElement('div')
            over.classList.add('over-layer');

            document.body.append(over, popMessage);
            document.body.classList.add('has-overlay');

            allMessages.scrollTop = allMessages.scrollHeight;
            over.addEventListener('click', (e) => {
                popMessage.remove();
                over.remove();
                document.body.classList.remove('has-overlay');
                const notification = document.querySelector(`#${name} .notification`);
                notification.classList.add('hide');
                const counter = notification.querySelector('.notification-counter');
                counter.textContent = 0;
            });
        } catch (error) {
            console.error(error);
        }
    }

    setupConversation(name) {
        const allMessages = document.querySelector('.messages-section');
        const send = document.querySelector('.message-input1');
        const overlay = document.querySelector('.over-layer')

        const event = async (event) => {
            if (event.key === "Enter" && !event.shiftKey) {
                const err = document.querySelector('.error-message')
                err.textContent = ''

                const message = send.value.trim();
                if (message.length > 200 || message.length <= 0 || !message) {
                    const err = document.querySelector('.error-message')
                    err.textContent = 'invalid message'
                    return
                }

                const validCookie = await validCookies();
                if (validCookie.valid) {
                    try {
                        this.base.connection.send(JSON.stringify({
                            ReceiverName: name,
                            Data: message,
                        }));

                        const conversation = document.querySelector('.conversation');
                        const username = conversation.getAttribute('name');

                        const messageCompon = document.createElement('div');
                        messageCompon.classList.add('message', 'receiver');
                        messageCompon.innerHTML = `
                        <div class="message-header">
                            <span class="username-message">${username}</span>
                            <span class="timestamp-mssg">${new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}</span>
                        </div>
                        <p>${escapeHTML(message)}</p>`
                        document.querySelector('.messages-section').appendChild(messageCompon);
                        send.value = '';
                        allMessages.scrollTop = allMessages.scrollHeight;
                    } catch (error) {
                        console.error(error);
                    }
                } else {
                    const popMessage = document.querySelector('.conversation');
                    const over = document.querySelector('.over-layer');
                    popMessage.remove();
                    over.remove();
                    this.base.handleLogout();
                }
            }
        }

        document.addEventListener("keydown", event)
        overlay.addEventListener('click', (e) => { document.removeEventListener("keydown", event) })
    }
}

function insertMessages(messages, container, name) {
    const currentScroll = container.scrollTop;
    const prevHeight = container.scrollHeight;

    for (let i = messages.length - 1; i >= 0; i--) {
        const messageCompon = document.createElement('div');
        messageCompon.classList.add('message');
        messageCompon.id = name;
        messageCompon.classList.add(messages[i].IsSender ? 'receiver' : 'sender');
        messageCompon.innerHTML = `
        <div class="message-header">
            <span class="username-message">${messages[i].IsSender ? messages[i].sender_name : name}</span>
            <span class="timestamp-mssg">${new Date(messages[i].CreatedAt).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}</span>
        </div>
        <p>${escapeHTML(messages[i].Message)}</p>`;

        // Append to the end instead of inserting at the beginning
        // container.appendChild(messageCompon);
        container.insertAdjacentElement('afterbegin', messageCompon);
    }

    // Maintain scroll position
    const newHeight = container.scrollHeight;
    if (newHeight !== prevHeight) {
        container.scrollTop = currentScroll + (newHeight - prevHeight);
    }
}