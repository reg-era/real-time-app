import { handleResize, debounce } from "../libs/script.js";

export class BASE {
    constructor(params) {
        this.params = params;
        this.setStyle("http://localhost:8080/assets/css/base.css")
        this.setStyle("http://localhost:8080/assets/css/posts.css")
    }

    setTitle(title) {
        document.title = title;
    }

    setStyle(link) {
        const existingLink = Array.from(document.head.getElementsByTagName('link'))
            .some(el => el.href === link);

        if (!existingLink) {
            const linkElement = document.createElement('link');
            linkElement.rel = 'stylesheet';
            linkElement.href = link;
            document.head.appendChild(linkElement);
        }
    }

    setListners() {
        const authNav = document.getElementById('auth-nav');
        const hasSession = document.cookie.includes('session_token');

        if (hasSession) {
            authNav.innerHTML = `
                <a class="active" data-link>Logout</a>
            `;
            const logoutLink = authNav.querySelector('.active');
            logoutLink.addEventListener('click', async (event) => {
                event.preventDefault()
                try {
                    const response = await fetch('http://localhost:8080/api/logout', {
                        method: 'POST',
                        credentials: 'include'
                    });

                    if (response.ok) {
                        window.location.href = "/"
                    } else {
                        console.error('Logout failed');
                    }
                } catch (error) {
                    console.error('Error logging out:', error);
                }
            });
        } else {
            authNav.innerHTML = `
                <a href="/login" class="active" data-link>Login</a>
                <a href="/register" data-link>Signup</a>
            `;
        }

        handleResize()
        let debouncedHandleResize = debounce(handleResize, 100)
        window.addEventListener('resize', debouncedHandleResize)

        const menuButton = document.querySelector('.menu-button');
        const sideBar = document.querySelector('.sidebar');
        if (menuButton && sideBar) {
            menuButton.addEventListener('click', () => {
                sideBar.classList.toggle('hide');
            });
        }
    }

    getSideBar() {
        return `
        <aside class="sidebar">
            <nav class="sidebar-nav">
                <a href="/new-post" class="nav__link" data-link >Create Post</a>
                <a href="/Created-post" class="nav__link" data-link >Created Posts</a>
                <a href="/Liked-post" class="nav__link" data-link >Liked Posts</a>
                <a href="/Categories" class="nav__link" data-link >Categories</a>
            </nav>
        </aside>
        `
    }

    getNavBar() {
        return `
        <header>
            <button class="menu-button">☰</button>
            <a href="/" class="nav__link" data-link >
                <div class="logo">
                    <img src="http://localhost:8080/assets/icons/logo.png" alt="Logo">
                </div>
            </a>
            <nav class="top-bar" id="auth-nav">
            </nav>
        </header>
        `
    }

    getHtmlBase() {
        const html = `
        ${this.getNavBar()}
        `
        setTimeout(this.setListners, 0);
        return html
    }
}