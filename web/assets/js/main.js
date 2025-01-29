import { navigateTo, router } from "./routes.js";

window.addEventListener("popstate", router);

document.addEventListener("DOMContentLoaded", () => {
    document.body.addEventListener("click", e => {
        if (e.target.matches("[data-link]")) {
            e.preventDefault();
            navigateTo(e.target.href);
        }
    });

    router();
});

const websocket = new WebSocket('ws://localhost:8080/api/ws');

websocket.onopen = () => {
    console.log('Connected to the server');
};

websocket.onclose = () => {
    console.log('Disconnected from the server');
};

websocket.onmessage = (event) => {
    const data = JSON.parse(event.data);
    if (window.location.pathname === ('/messages/' + data.from)) {
        const container = document.querySelector('.messages-section');
        const messageCompon = document.createElement('div');
        messageCompon.classList.add('message', 'sender')
        messageCompon.innerHTML = `<p>${data.message}</p>`;
        container.appendChild(messageCompon)
    }
};

// window.scrollTo({
// top: document.body.scrollHeight,
// behavior: 'smooth'
// });