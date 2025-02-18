import { BASE } from "./views/_BASE.js"
import { Router } from "./rootes.js";

export const app = new BASE();

async function main() {
    app.router = new Router(app);

    window.addEventListener('popstate', () => {
        app.router.handleRoute()
    });

    try {
        await app.router.handleRoute();

        if (!app.connection) {
            await app.initializeWebSocket();
        }
    } catch (error) {
        console.error('Failed to handle route:', error);
    }
}

const validCookies = async () => {
    try {
        const res = await fetch('/api/me/check-in');
        if (res.status === 202) {
            const body = await res.json(); // or res.text() based on your response format
            return { valid: true, body };
        }
        return { valid: false };
    } catch (error) {
        console.error(error);
        return { valid: false };
    }
}


await main();

export { validCookies };