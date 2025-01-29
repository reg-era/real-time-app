import { Home } from './views/home.js';
import { Login } from './views/login.js';
import { Register } from './views/register.js';
import { Messages } from './views/messages.js';
import { NewPost } from './views/newPost.js';
import { Error } from './views/error.js';
import { Messg } from './views/WsHub.js';
import { app } from './main.js';

export class Router {
    constructor() {
        this.routes = [
            { path: "/", view: Home, name: "home" },
            { path: "/register", view: Register, name: "register" },
            { path: "/login", view: Login, name: "login" },
            { path: "/messages", view: Messages, name: "messages" },
            { path: "/new-post", view: NewPost, name: "new-post" },
            { path: "/ws", view: Messg, name: ",msgs" }
        ];
        this.base = app;
        this.eventlistener = this.handleClick.bind(this); // Bind the listener once
        this.init();
        this.page = {};
    }

    init() {
        document.removeEventListener('click', this.eventlistener); // Ensure no duplicates
        document.addEventListener('click', this.eventlistener);
    }

    handleClick(e) {
        console.log('triggered');

        if (e.target.matches('[data-link]')) {
            e.preventDefault();
            const href = e.target.getAttribute('href');
            this.navigateTo(href);
        }
    }

    async handleRoute() {
        const path = window.location.pathname;
        const route = this.routes.find(r => r.path === path);
        const hasSession = document.cookie.includes('session_token');

        if (route) {
            const view = new route.view(this.getQueryParams());
            this.page = view;
            // Check for authentication
            if (!hasSession && (route.name !== 'login' && route.name !== "register")) {
                console.log("not authorized");
                if (this.base.connection) {
                    this.base.connection.close();
                }
                const view = new Login();
                const html = await view.renderHtml();
                const appElement = document.querySelector('.app');
                appElement.innerHTML = html;
                appElement.setAttribute('page', 'error');
                if (typeof view.afterRender === 'function') {
                    view.afterRender();
                }
                return;
            }

            const appElement = document.querySelector('.app');
            console.log(appElement.getAttribute('page'));

            // Render only if the page has changed
            if (appElement.getAttribute('page') !== route.name) {
                if (!this.base.connection) {
                    this.base.initializeWebSocket();
                }
                const html = await view.renderHtml();
                appElement.innerHTML = html;
                appElement.setAttribute('page', route.name);
                if (typeof view.afterRender === 'function') {
                    view.afterRender();
                }
            }
        } else {
            // Handle 404 case
            const errorView = new Error("404");
            const html = await errorView.renderHtml();
            const appElement = document.querySelector('.app');
            appElement.innerHTML = html;
            appElement.setAttribute('page', 'error');
            if (typeof errorView.afterRender === 'function') {
                errorView.afterRender();
            }
        }
    }

    getQueryParams() {
        const params = {};
        const queryParams = new URLSearchParams(window.location.search);
        queryParams.forEach((value, key) => params[key] = value);
        return params;
    }

    navigateTo(url) {
        // if (!window.location.pathname === url) {
        // }
        history.pushState(null, null, url);
        this.handleRoute();
    }
}
