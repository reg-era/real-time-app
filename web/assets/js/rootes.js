import { Home } from "./views/home.js";
import { Login } from "./views/login.js";
import { Register } from "./views/register.js";
import { NewPost } from "./views/newPost.js";
import { Error } from "./views/error.js";

const getParams = (path, routePath) => {
    const pathParts = path.split('/').filter(part => part !== '');
    const routeParts = routePath.split('/').filter(part => part !== '');

    // only get params from querys and indexes like ids for displaying only one post page inchaaaalah hhhhhh
    const params = {};

    routeParts.forEach((part, index) => {
        if (part.startsWith(':')) {
            const paramName = part.slice(1);
            params[paramName] = pathParts[index];
        }
    });

    const queryParams = new URLSearchParams(window.location.search);
    queryParams.forEach((value, key) => {
        params[key] = value;
    });

    return params;
};


const router = async () => {
    const routes = [
        { path: "/", view: Home },
        { path: "/login", view: Login },
        { path: "/register", view: Register },
        { path: "/new-post", view: NewPost },
        // { path: "/posts/:id", view: PostDetails }, khliw had twichiya tal mn ba3d
        { path: "/error", view: Error }
    ];

    const path = location.pathname;

    let match = routes.find(route => {
        const pathParts = path.split('/').filter(part => part !== '');
        const routeParts = route.path.split('/').filter(part => part !== '');

        if (pathParts.length !== routeParts.length) {
            return false;
        }

        return routeParts.every((part, index) => {
            return part.startsWith(':') || part === pathParts[index];
        });
    });

    if (!match) {
        window.location.href = '/error?status=404'
    }

    const params = getParams(path, match.path);

    const view = new match.view(params);

    document.querySelector(".app").innerHTML = await view.getHtml();
};

const navigateTo = url => {
    history.pushState(null, null, url);
    router();
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
