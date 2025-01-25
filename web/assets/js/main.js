import { Router } from './rootes.js';
import { Home } from './views/home.js';
import { Login } from './views/login.js';
import { Register } from './views/register.js';
import { Messages } from './views/messages.js';
// import { Error404 } from './views/error.js';

export const routes = [
    { path: "/", view: Home },
    { path: "/register", view: Register },
    { path: "/login", view: Login },
    { path: "/messages", view: Messages },
    // { path: "/404", view: Error404 }
];

const router = new Router(routes);