import { BASE } from "./views/_BASE.js"
import { Router } from "./rootes.js";

const app = new BASE();
app.router = new Router();
console.log(app.router);

await app.router.handleRoute();
 app.initializeWebSocket();
