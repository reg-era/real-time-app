import { handleResize, debounce } from "../libs/script.js";
import { popup } from "./popup.js";
// import { app } from "../main.js";
// import { Router } from "../rootes.js";

export class BASE {
    constructor(params) {
        this.router = null;
        this.params = params;
        this.loged = false;
        // this.currentPage = this.router.page;
        this.styleUrls = [
            'http://localhost:8080/api/css/base.css',
            'http://localhost:8080/api/css/posts.css'
        ];
        this.users = {
            online: new Set([]),
            offline: new Set([])
        };
        this.connection = null;
        this.initializeStyles();
        //   this.initializeWebSocket();
    }

    async initializeWebSocket() {
        // Don't create a new connection if one exists and is healthy
        if (this.connection &&
            (this.connection.readyState === WebSocket.CONNECTING ||
                this.connection.readyState === WebSocket.OPEN)) {
            return;
        }

        const sessionToken = document.cookie
            .split('; ')
            .find(row => row.startsWith('session_token='))
            ?.split('=')[1];

        const wsUrl = new URL('ws://localhost:8080/ws');
        if (sessionToken) {
            wsUrl.searchParams.append('token', sessionToken);
        }

        this.connection = new WebSocket(wsUrl.toString());
        console.log('Initializing WebSocket connection...');

        this.setupWebSocket();
    }

    setupWebSocket() {
        this.connection.onopen = async () => {
            console.log('WebSocket connection established');
            this.setupConnReader();
        };

        this.connection.onerror = (error) => {
            console.error('WebSocket error:', error);
        };

        this.connection.onclose = (event) => {
            console.log('WebSocket closed:', event.code, event.reason);
            this.handleLogout();
        };
    }

    setupConnReader() {
        this.connection.onmessage = async (event) => {
            console.log(event, "received");

            try {
                const data = JSON.parse(event.data);
                if (!data.Type) {
                    console.error('Received message without type:', data);
                    return;
                }

                //  console.log(data.users);
                switch (data.Type) {
                    case 'message':
                        this.handleWebSocketMessage(data);
                        break;
                    case 'onlineusers':
                        if (data.users) {
                            this.users = data.users;
                            await this.renderSidebar();
                        }
                        break;
                    default:
                        console.warn('Unknown message type:', data.type);
                }
            } catch (error) {
                console.error('Error processing WebSocket message:', error);
            }
        };
    }

    handleWebSocketMessage(data) {
        console.log('Received message:', data);
        // Implement specific message handling logic here
    }

    initializeStyles() {
        this.styleUrls.forEach(url => this.setStyle(url));
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

    setTitle(title) {
        document.title = title;
    }

    async handleLogout() {
        console.log(this.router);

        try {
            const response = await fetch('http://localhost:8080/api/logout', {
                method: 'POST',
                credentials: 'include'
            });

            if (response.ok) {
                const authNav = document.getElementById('auth-nav');
                authNav.innerHTML = `
                    <span href="/login" class="active" data-link>Login</span>
                    <span href="/register" data-link>Signup</span>
                `;
                if (this.connection) {
                    this.connection.close();
                }
                history.pushState(null, null, '/login');
                this.router.handleRoute()
            } else {
                throw new Error('Logout failed');
            }
        } catch (error) {
            console.error('Logout error:', error);
        }
    }

    setupAuthNav(app) {
        const authNav = document.getElementById('auth-nav');
        const hasSession = document.cookie.includes('session_token');

        authNav.innerHTML = hasSession
            ? '<span class="active" data-link>Logout</span>'
            : `
                <span href="/login" class="active" data-link>Login</span>
                <span href="/register" data-link>Signup</span>
              `;

        if (hasSession) {
            authNav.querySelector('.active').addEventListener('click', () => app.handleLogout());
        }
    }

    setupSidebar() {
        handleResize();
        const debouncedHandleResize = debounce(handleResize, 100);
        window.addEventListener('resize', debouncedHandleResize);

        const menuButton = document.querySelector('.menu-button');
        const sideBar = document.querySelector('.sidebar');
        if (menuButton && sideBar) {
            menuButton.addEventListener('click', () => {
                sideBar.classList.toggle('hide');
            });
        }
    }

    async getSidebar() {
        return `
        <aside class="sidebar">
            <nav class="sidebar-nav">
                <span href="/new-post" class="nav__link" data-link>Create Post</span>
                <span href="/ws" class="nav__link" data-link>Messages</span>
            </nav>
        </aside>
        ${await this.getOnlineBar()}
        `;
    }

    getNavBar() {
        return `
        <header>
            <button class="menu-button">☰</button>
            <span href="/" class="nav__link" data-link>
                <div class="logo" href="/" data-link>
                    <img src="http://localhost:8080/api/icons/logo.png" alt="Logo" href="/" data-link>
                </div>
            </span>
            <nav class="top-bar" id="auth-nav"></nav>
        </header>
        `;
    }

    async getOnlineBar() {
        return `
        <aside class="onligne-bar">
            <nav class="sidebar-nav">
                <div class="loading-indicator">Loading users...</div>
                ${await this.renderSidebar()}
            </nav>
        </aside>
        `;
    }

    getHtmlBase() {
        return this.getNavBar();
    }

    setupmssglistner(app) {
        document.addEventListener('click', (event) => {
            const linkElement = event.target.closest('[data-mssg-link]');
            if (linkElement) {
                event.preventDefault();
                const pop = new popup(app);
                console.log(linkElement.getAttribute('id'));
                pop.getMessages(linkElement.getAttribute('id'));
            }
        });
    }
    setupNavigation(app) {
        document.querySelectorAll('[data-link]').forEach(link => {
            link.addEventListener('click', (event) => {
                event.preventDefault();
                const href = link.getAttribute('href');
                if (href) {
                    window.history.pushState(null, null, href);
                    app.router.handleRoute();
                }
            });
        });
    }

    async renderSidebar() {
        try {
            const sidebar = document.querySelector('.onligne-bar .sidebar-nav');
            if (!sidebar) {
                console.error('Sidebar element not found');
                return;
            }
            sidebar.innerHTML = '';
            let html = '';

            // Render online users
            if (Array.isArray(this.users.online)) {
                html += this.users.online.map(user => `
                    <a class="nav__link" id="${user}" data-mssg-link>
                        👤 ${user} 
                        <span class="status-tag online">
                            🟢Online
                        </span>
                    </a>
                `).join('');
            }

            // Render offline users
            if (Array.isArray(this.users.offline)) {
                html += this.users.offline.map(user => `
                    <a class="nav__link" id="${user}" data-mssg-link>
                        👤 ${user} 
                        <span class="status-tag offline">
                            🔴Offline
                        </span>
                    </a>
                `).join('');
            }

            sidebar.innerHTML = html || '<div>No users available</div>';
            //  this.setupNavigation();
        } catch (error) {
            console.error('Error rendering sidebar:', error);
            const sidebar = document.querySelector('.onligne-bar .sidebar-nav');
            if (sidebar) {
                sidebar.innerHTML = '<div class="error-message">Error loading users</div>';
            }
        }
    }

    cleanup() {
        if (this.connection) {
            this.connection.close();
        }
        window.removeEventListener('resize', this.debouncedHandleResize);
    }

    afterRender() {
        this.renderSidebar();
        this.setupAuthNav(this);
        this.setupSidebar();
        this.setupNavigation(this);
        this.setupmssglistner(this);
    }
}