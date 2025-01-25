export class Router {
    constructor(routes) {
        this.routes = routes;
        this.init();
    }

    init() {
        window.addEventListener('popstate', () => this.handleRoute());
        document.addEventListener('DOMContentLoaded', () => this.handleRoute());
        document.addEventListener('click', e => {
            if (e.target.matches('[data-link]')) {
                e.preventDefault();
                const href = e.target.getAttribute('href') || e.target.closest('[data-link]').getAttribute('href');
                this.navigateTo(href);
            }
        });
    }

    async handleRoute() {
        console.log("handled");
        const path = window.location.pathname;
        const route = this.routes.find(r => r.path === path) || this.routes.find(r => r.path === '/404');
        if (route) {
            const view = new route.view(this.getQueryParams());

            const html = await view.renderHtml();
            document.querySelector('.app').innerHTML = html;
            if (typeof view.afterRender === 'function') {
                view.afterRender();
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
        history.pushState(null, null, url);
        this.handleRoute();
    }


}
