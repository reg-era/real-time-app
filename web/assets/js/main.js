import { BASE } from "./views/_BASE.js"
import { Router } from "./rootes.js";

export const app = new BASE();
app.router = new Router();
await app.router.handleRoute();

//app.initializeWebSocket();
