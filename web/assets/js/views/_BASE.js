import { handleResize, debounce } from "../libs/script.js";
import { validCookies } from "../main.js";
import { popup } from "./popup.js";

export class BASE {
    constructor(params) {
        this.router = null;
        this.params = params;
        this.loged = false;
        this.styleUrls = [
            '/api/css/base.css',
            '/api/css/posts.css',
            '/api/css/messages.css',
        ];
        this.users = {
            Friends: [],
        };
        this.mssglistener = null;
        this.navlistener = null;
        this.connection = null;
        this.onlineusers = null;
        this.initializeStyles();
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

    async initializeWebSocket() {
        if (this.connection &&
            (
                this.connection.readyState === WebSocket.CONNECTING ||
                this.connection.readyState === WebSocket.OPEN
            )
        ) {
            return;
        }

        const sessionToken = document.cookie
            .split('; ')
            .find(row => row.startsWith('session_token='))
            ?.split('=')[1];

        const wsUrl = new URL(`ws://${window.location.hostname}:8080/ws`);
        if (sessionToken) {
            wsUrl.searchParams.append('token', sessionToken);
        }

        this.connection = new WebSocket(wsUrl.toString());
        console.log('Initializing WebSocket connection...');

        this.connection.onopen = async () => {
            console.log('WebSocket connection established');
            this.setupConnReader();
        };

        this.connection.onclose = (event) => {
            console.log('WebSocket closed:', event.code, event.reason);
        };

    }

    async setupConnReader() {
        this.connection.onmessage = async (event) => {


            try {
                const data = JSON.parse(event.data);
                if (!data.Type) {
                    console.error('Received message without type:', data);
                    return;
                }
                switch (data.Type) {
                    case 'message':
                        await this.handleWebSocketMessage(data);
                        break;
                    case 'onlineusers':
                        if (data.users) {
                            this.users = data.users;
                            this.renderSidebar();
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

    async handleWebSocketMessage(message) {
        const allMessages = document.querySelector('.messages-section');
        const conversation = document.querySelector('.conversation');

        if (allMessages && conversation.getAttribute('name') === message.Message.sender_name) {
            const msg = document.createElement('div');
            msg.classList.add('message', 'sender')
            msg.innerHTML = `
            <div class="message-header">
                <span class="username-message">${message.Message.sender_name}</span>
                <span class="timestamp-mssg">${new Date(message.Message.CreatedAt).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}</span>
            </div>
            <p>${message.Message.Message}</p>`
            allMessages.insertAdjacentElement("beforeend", msg);
            allMessages.scrollTop = allMessages.scrollHeight;
        } else {
            const notification = document.querySelector(`#${message.Message.sender_name} .notification`);
            notification.classList.remove('hide');
            const counter = notification.querySelector('.notification-counter');
            counter.textContent = parseInt(counter.textContent) + 1;
            showNotification(message.Message);
        }
    }

    async handleLogout() {
        try {
            const response = await fetch('/api/logout', {
                method: 'POST',
                credentials: 'include'
            });

            if (response.ok) {
                const authNav = document.getElementById('auth-nav');
                if (authNav) {
                    authNav.innerHTML = `
                        <span href="/login" class="active" data-link>Login</span>
                        <span href="/register" data-link>Signup</span>
                    `;
                }
                if (this.connection) {
                    this.connection.close();
                }
                history.pushState(null, null, '/login');
                await this.router.handleRoute();
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
        handleResize()
        let debouncedHandleResize = debounce(handleResize, 100)
        window.addEventListener('resize', debouncedHandleResize)

        const menuButton = document.querySelector('.menu-button');
        const sideBar = document.querySelector('.sidebar-for-min');
        if (menuButton && sideBar) {
            menuButton.addEventListener('click', () => {
                sideBar.classList.toggle('hide');
            });
        }
    }

    getSidebar() {
        return `
        <aside class="sidebar-for-min">
        <span href="/new-post" class="creat-post" data-link>Create Post</span>
        <nav class="sidebar-nav-inside">
        </nav>
        </aside>
        `;
    }

    getNavBar() {
        return `
        <header>
            <button class="menu-button">☰</button>
            <span href="/new-post" class="creat-post" data-link>Create Post</span>
            <span href="/" class="nav__link" data-link>
                <div class="logo" href="/" data-link>
                    <img src="/api/icons/logo.png" alt="Logo" href="/" data-link>
                </div>
            </span>
            <buttom
            <nav class="top-bar" id="auth-nav"></nav>
        </header>
        `;
    }

    getOnlineBar() {
        return `
        <aside class="onligne-bar">
            <nav class="sidebar-nav">
            </nav>
        </aside>
        `;
    }

    getHtmlBase() {
        return this.getNavBar();
    }

    async setupmssglistner(app) {
        if (app.mssglistener !== null) return;
        app.mssglistener = document.addEventListener('click', async (event) => {
            const linkElement = event.target.closest('[data-mssg-link]');
            if (linkElement) {
                const pop = new popup(app);
                await pop.getMessages(linkElement.getAttribute('id'));
                pop.setupConversation(linkElement.getAttribute('id'));
            }
        });
    }

    async setupNavigation(app) {
        if (app.navlistener !== null) return;
        app.navlistener = document.addEventListener('click', async (event) => {
            const linkElement = event.target.closest('[data-link]');
            if (linkElement) {
                event.preventDefault();
                const href = linkElement.getAttribute('href');
                if (href) {
                    window.history.pushState(null, null, href);
                    app.router.handleRoute();
                }
            }
        });
    }

    renderSidebar() {
        try {
            //this for online bar 
            let onlinebar = document.querySelector('.onligne-bar');
            if (onlinebar) {
                const nav = onlinebar.querySelector('.sidebar-nav');
                nav.innerHTML = '';
                if (Array.isArray(this.users.Friends)) {
                    this.users.Friends.forEach(user => {
                        if (user.Seen === 0 && !user.IsSender) {
                            showNotification({
                                Message: user.LastMessage,
                                sender_name: user.Name,
                            });
                        }
                        nav.appendChild(makeBar(user.Online, user));
                    });
                }
            }
            //this for side bar hiden
            const sidebar = document.querySelector('.sidebar-for-min');
            if (sidebar) {
                const navbar = sidebar.querySelector('.sidebar-nav-inside');
                navbar.innerHTML = '';
                if (Array.isArray(this.users.Friends)) {
                    this.users.Friends.forEach(user => navbar.appendChild(makeBar(user.Online, user)));
                }
            }
            handleResize();
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

function showNotification(message) {
    const notificationContainer = document.getElementById('notification-container');
    const oldnotif = notificationContainer.querySelector(`#${message.sender_name}`);
    if (oldnotif) oldnotif.remove();

    const notification = document.createElement('div');
    notification.className = 'notification';

    notification.innerHTML = `
        ${message.Message}
      <span>  ${message.sender_name}</span>
        <span class="close-btn">&times;</span>
    `;

    notification.setAttribute('data-mssg-link', null);
    notification.id = message.sender_name;

    notificationContainer.appendChild(notification);

    notification.style.display = 'block';

    const closeBtn = notification.querySelector('.close-btn');
    closeBtn.addEventListener('click', (e) => {
        e.stopPropagation();
        notification.remove();
    });

    setTimeout(() => {
        notification.remove();
    }, 5000);
}

function makeBar(online, user) {
    const bar = document.createElement('div')
    bar.setAttribute('data-mssg-link', null)
    bar.id = user.Name
    bar.classList.add('status-bar')
    bar.innerHTML = `
    <div class="status-info">
        <span id="online-status" class="username">${online ? "🟢" : "🔴"} ${user.Name}</span>
        <span class="lastmessage"> ${user.LastMessage} </span>
        <span class="timestamp">${user.LastMessage ? new Date(user.Time).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }) : ''}</span>
    </div>
    <div class="notification hide"><span class="notification-counter">0</span></div>`
    return bar
}