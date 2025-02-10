import { validCookies } from "../main.js";
import { debounce } from "../libs/script.js";

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
        popMessage.setAttribute('name', res.body);


        const inputMessage = document.createElement('div')
        inputMessage.classList.add('messages-input');
        inputMessage.innerHTML = `
            <input required placeholder="Type message ..." class="message-input1"></input>
            <p class="error-comment"></p>`;

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
                // Store original data to avoid mutation
                let messageQueue = [...data];
                const BATCH_SIZE = 10;
                const loadMoreMessages = debounce((container) => {
                    // Remove existing "See More" button
                    const existingButton = container.querySelector('.seemore');
                    if (existingButton) {
                        existingButton.remove();
                    }

                    // Get next batch of messages
                    const nextBatch = messageQueue.splice(-BATCH_SIZE);

                    if (nextBatch.length > 0) {
                        insertMessages(nextBatch, container);

                        // Only add "See More" if there are more messages
                        if (messageQueue.length > 0) {
                            const seeMoreButton = document.createElement('div');
                            seeMoreButton.className = 'seemore';
                            seeMoreButton.textContent = `See More (${messageQueue.length} remaining)`;
                            container.insertBefore(seeMoreButton, container.firstChild);

                            // Add click handler to new button
                            seeMoreButton.addEventListener('click', () => loadMoreMessages(container));
                        }
                    }
                }, 500);

                // Initial load
                loadMoreMessages(allMessages);
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
        const send = document.querySelector('.message-input1');
        document.addEventListener("keydown", async (event) => {
            if (event.key === "Enter" && !event.shiftKey) {
                const message = send.value.trim();
                if (message) {
                    try {
                        this.base.connection.send(JSON.stringify({
                            ReceiverName: name,
                            Data: message,
                        }));
                        const conversation = document.querySelector('.conversation');
                        const username = conversation.getAttribute('name');

                        const messageCompon = document.createElement('div');
                        messageCompon.classList.add('message', 'receiver');
                        // handle the name of loged user !!!!
                        messageCompon.innerHTML = messageCompon.innerHTML = `
                        <div class="message-header">
                          <span class="username-message">${username}</span>
                          <span class="timestamp-mssg">${new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}</span>
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

function insertMessages(messages, container) {
    // Store current scroll position
    const currentScroll = container.scrollTop;
    const prevHeight = container.scrollHeight;

    messages.forEach(mssg => {
        const messageCompon = document.createElement('div');
        messageCompon.classList.add('message');
        messageCompon.id = name;
        messageCompon.classList.add(mssg.IsSender ? 'receiver' : 'sender');
        messageCompon.innerHTML = `
      <div class="message-header">
        <span class="username-message">${mssg.IsSender ? mssg.sender_name : name}</span>
        <span class="timestamp-mssg">${new Date(mssg.CreatedAt).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}</span>
      </div>
      <p>${mssg.Message}</p>`;
        container.insertBefore(messageCompon, container.firstChild);
    });

    // Maintain scroll position
    const newHeight = container.scrollHeight;
    if (newHeight !== prevHeight) {
        container.scrollTop = currentScroll + (newHeight - prevHeight);
    }
};