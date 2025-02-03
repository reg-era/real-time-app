import { handleResize, debounce } from "../libs/script.js";
import { popup } from "./popup.js";

export class BASE {
    constructor(params) {
        this.router = null;
        this.params = params;
        this.loged = false;
        this.styleUrls = [
            'http://localhost:8080/api/css/base.css',
            'http://localhost:8080/api/css/posts.css',
            'http://localhost:8080/api/css/messages.css',
        ];
        this.users = {
            online: new Set([]),
            offline: new Set([])
        };
        this.connection = null;
        this.initializeStyles();
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
            // this.handleLogout();
        };
    }

    setupConnReader() {
        this.connection.onmessage = async (event) => {
            try {
                const data = JSON.parse(event.data);
                if (!data.Type) {
                    console.error('Received message without type:', data);
                    return;
                }

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

    handleWebSocketMessage(message) {
        const allMessages = document.querySelector('.messages-section');
        const conversation = document.querySelector('.messages-section');

        if (conversation) {
            const msg = document.createElement('div');
            msg.classList.add('message', 'sender')
            msg.innerHTML = `<p>${message.Message.Message}</p>`
            conversation.insertAdjacentElement("beforeend", msg);
            allMessages.scrollTop = allMessages.scrollHeight;
        } else {
            const notification = document.querySelector(`#${message.Message.sender_name} .notification`)
            notification.classList.remove('hide')
            const counter = notification.querySelector('.notification-counter')
            counter.textContent = parseInt(counter.textContent) + 1
        }
    }
    setTitle(title) {
        document.title = title;
    }

    async handleLogout() {
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

    async setupmssglistner(app) {
        document.addEventListener('click', async (event) => {
            const linkElement = event.target.closest('[data-mssg-link]');
            if (linkElement) {
                // event.preventDefault();
                const pop = new popup(app);
                await pop.getMessages(linkElement.getAttribute('id'));
                pop.setupConversation(linkElement.getAttribute('id'));
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
        const makeBar = (online, user) => {
            const bar = document.createElement('div')
            bar.setAttribute('data-mssg-link', null)
            bar.id = user
            bar.classList.add('status-bar')
            bar.innerHTML = `
            <div class="status-info">
                <span id="online-status" class="status ${online ? "online" : "offline"}">${online ? "🟢" : "🔴"}</span>
                <span id="username" class="username">${user} </span>
            </div>
            <div class="notification hide"><span class="notification-counter">0</span></div>`
            return bar
        }

        try {
            const sidebar = document.querySelector('.onligne-bar .sidebar-nav');
            if (!sidebar) {
                console.error('Sidebar element not found');
                return;
            }

            sidebar.innerHTML = '';

            if (Array.isArray(this.users.online)) {
                this.users.online.forEach(user => sidebar.appendChild(makeBar(true, user)));
            }

            if (Array.isArray(this.users.offline)) {
                this.users.offline.forEach(user => sidebar.appendChild(makeBar(false, user)));
            }
        } catch (error) {
            console.error('Error rendering sidebar:', error);
            const sidebar = document.querySelector('.onligne-bar .sidebar-nav');
            if (sidebar) {
                sidebar.innerHTML = '<div class="error-message">Error loading users</div>';
            }
        }
    }

    afterRender() {
        this.renderSidebar();
        this.setupAuthNav(this);
        this.setupSidebar();
        this.setupNavigation(this);
    }
}