import Home from "./views/home.js";
import Login from "./views/login.js";
import Register from "./views/register.js";
import newPost from "./views/newPost.js";
import ErrorPage from "./views/error.js";

const getParams = (path, routePath) => {
    const pathParts = path.split('/').filter(part => part !== '');
    const routeParts = routePath.split('/').filter(part => part !== '');

    const params = {};
    // still made custumize for more type of params URL like query indexes ....
    routeParts.forEach((part, index) => {
        if (part.startsWith(':')) {
            params[part.slice(1)] = pathParts[index];
        }
    });

    return params;
};

const navigateTo = url => {
    history.pushState(null, null, url);
    router();
};

const router = async () => {
    const routes = [
        { path: "/", view: Home },
        { path: "/login", view: Login },
        { path: "/register", view: Register },
        { path: "/new-post", view: newPost },
        { path: "/error", view: ErrorPage },
    ];

    const path = location.pathname;

    let match = routes.find(route => {
        const pathParts = path.split('/').filter(part => part !== '');
        const routeParts = route.path.split('/').filter(part => part !== '');

        if (pathParts.length !== routeParts.length) {
            return false;
        }

        return routeParts.every((part, index) => part === pathParts[index]);
    });

    if (!match) { // make sure to pute the error page on the end of routes slice
        match = { route: routes[routes.length - 1], result: [path] };
    }

    const params = getParams(path, match.path);
    const view = new match.view(params);

    document.querySelector(".app").innerHTML = await view.getHtml();
};

window.addEventListener("popstate", router);

document.addEventListener("DOMContentLoaded", () => {
    document.body.addEventListener("click", e => {
        if (e.target.matches("[data-link]")) {
            e.preventDefault();
            navigateTo(e.target.href);
        }
    });

    router();
});
