import { handleResize, debounce } from "../libs/script.js";
// import { Error } from "./error.js";
export class BASE {
    constructor(params) {
        this.params = params;
        this.styleUrls = [
            'http://localhost:8080/api/css/base.css',
            'http://localhost:8080/api/css/posts.css'
        ];
        this.initializeStyles();
    }

    initializeStyles() {
        this.styleUrls.forEach(url => this.setStyle(url));
    }

    setTitle(title) {
        document.title = title;
    }

    setStyle(url) {
        const links = Array.from(document.head.getElementsByTagName('link'));
        if (!links.some(link => link.href === url)) {
            const linkElement = document.createElement('link');
            linkElement.rel = 'stylesheet';
            linkElement.href = url;
            document.head.appendChild(linkElement);
        }
    }

    async handleLogout() {
        try {
            const response = await fetch('http://localhost:8080/api/logout', {
                method: 'POST',
                credentials: 'include'
            });

            if (response.ok) {
                window.location.href = '/';
            } else {
                throw new Error('Logout failed');
            }
        } catch (error) {
            console.error('Logout error:', error);
        }
    }

    setupAuthNav() {
        const authNav = document.getElementById('auth-nav');
        const hasSession = document.cookie.includes('session_token');

        authNav.innerHTML = hasSession
            ? '<span class="active" data-link>Logout</span>'
            : `
                <span href="/login" class="active" data-link>Login</span>
                <span href="/register" data-link>Signup</span>
              `;

        if (hasSession) {
            authNav.querySelector('.active').addEventListener('click', () => this.handleLogout());
        }
    }

    setupSidebar() {
        // const { handleResize, debounce } = require('../libs/script.js');
        const debouncedResize = debounce(handleResize, 100);
        window.addEventListener('resize', debouncedResize);

        const menuButton = document.querySelector('.menu-button');
        const sidebar = document.querySelector('.sidebar');

        if (menuButton && sidebar) {
            menuButton.addEventListener('click', () => sidebar.classList.toggle('hide'));
        }
    }

    getSidebar() {
        return `
        <aside class="sidebar">
            <nav class="sidebar-nav">
                <span href="/new-post" class="nav__link" data-link>Create Post</span>
                <span href="/filter/created-posts" class="nav__link" data-link>Created Posts</span>
                <span href="/filter/liked-posts" class="nav__link" data-link>Liked Posts</span>
                <span href="/filter/categories" class="nav__link" data-link>Categories</span>
                <span href="/ws" class="nav__link" data-link>Messages</span>
            </nav>
        </aside>
        `;
    }

    getNavBar() {
        return `
        <header>
            <button class="menu-button">☰</button>
            <span href="/" class="nav__link" data-link>
                <div class="logo" data-link>
                    <img src="http://localhost:8080/api/icons/logo.png" alt="Logo" data-link>
                </div>
            </span>
            <nav class="top-bar" id="auth-nav"></nav>
        </header>
        `;
    }

    getHtmlBase() {
        return this.getNavBar();
    }

    setupNavigation() {
        document.querySelectorAll('[data-link]').forEach(link => {
            link.addEventListener('click', (event) => {
                event.preventDefault();
                const href = link.getAttribute('href');
                if (href) window.history.pushState(null, null, href);

            });
        });
    }

    afterRender() {
        this.setupAuthNav();
        this.setupSidebar();
        this.setupNavigation();
    }
}