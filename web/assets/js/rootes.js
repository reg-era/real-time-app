import { Home } from './views/home.js';
import { Login } from './views/login.js';
import { Register } from './views/register.js';
// import { Messages } from './views/messages.js';
import { NewPost } from './views/newPost.js';
import { Error } from './views/error.js';

export class Router {
    constructor(app) {
        this.routes = [
            { path: "/", view: Home, name: "home" },
            { path: "/register", view: Register, name: "register" },
            { path: "/login", view: Login, name: "login" },
            // { path: "/messages", view: Messages, name: "messages" },
            { path: "/new-post", view: NewPost, name: "new-post" },
        ];
        this.base = app;
    }

    async handleRoute() {
        const path = window.location.pathname;
        const route = this.routes.find(r => r.path === path);
        const hasSession = document.cookie.includes('session_token');


        if (route) {
            const view = new route.view(this.base);
            this.page = view;
            // Check for authentication
            if ((!hasSession && (route.name !== 'login' && route.name !== "register")) || (route.name === 'login' || route.name === "register")) {
                console.log(document.cookie);

                if (this.base.connection) {
                    this.base.connection.close();
                }
                history.pushState(null, null, '/login');
                const view = new Login(this.base);
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

            // Render only if the page has changed
            if (appElement.getAttribute('page') !== route.name && hasSession) {
                const html = await view.renderHtml();
                appElement.innerHTML = html;
                appElement.setAttribute('page', route.name);
                if (typeof view.afterRender === 'function') {
                    view.afterRender();
                }
            }
        } else {
            // Handle 404 case
            const errorView = new Error("404", this.base);
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

    async navigateTo(url) {
        if (window.location.pathname !== url) {
            history.pushState(null, null, url);
            await this.handleRoute();
        }
    }
}
