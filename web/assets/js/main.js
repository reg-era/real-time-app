import { BASE } from "./views/_BASE.js"
import { Router } from "./rootes.js";


async function main() {
    const app = new BASE();
    app.router = new Router(app);
    if (!document.cookie) {
        deleteAllCookies();
        history.pushState(null, null, '/login');

        await app.router.handleRoute();
        console.log(app);


        try {
            //await app.initializeWebSocket();

        } catch (error) {
            console.error('Error during route handling:', error);
        }

    } else {
        await app.router.handleRoute();
        app.initializeWebSocket();


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

main();
//app.initializeWebSockety