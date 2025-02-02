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
// }

function deleteAllCookies() {
    const cookies = document.cookie.split(";");
    console.log(cookies);

    for (let cookie of cookies) {
        const name = cookie.split("=")[0].trim();
        document.cookie = name + "=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
    }
}


main();
//app.initializeWebSockety