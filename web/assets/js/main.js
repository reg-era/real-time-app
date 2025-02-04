import { BASE } from "./views/_BASE.js"
import { Router } from "./rootes.js";

async function main() {
    const app = new BASE();
    app.router = new Router(app);
    await app.router.handleRoute();
    if (!app.connection) {
        await app.initializeWebSocket();
    }
}

function deleteAllCookies() {
    const cookies = document.cookie.split(";");
    console.log(cookies);

    for (let cookie of cookies) {
        const name = cookie.split("=")[0].trim();
        document.cookie = name + "=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
    }
}

const validCookies = async () => {
    try {
        const res = await fetch('http://localhost:8080/api/me/check-in');
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