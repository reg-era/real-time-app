export default class {
    constructor(params) {
        this.params = params;
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

    getNavigation() {
        return `
        <aside class="sidebar">
            <nav class="sidebar-nav">
                <a href="/login" class="nav__link" data-link >login</a>
                <a href="/register" class="nav__link" data-link >register</a>
                <a href="/posts" class="nav__link" data-link >posts</a>
                <a href="/new-post" class="nav__link" data-link >newpost</a>
            </nav>
        </aside>
        `
    }

    getHtmlBase() {
        return `
        <header>
            <button class="menu-button">☰</button>
            <a href="/SP-framework/index.html" class="nav__link" data-link >
                <div class="logo">
                    <img src="http://localhost:8080/assets/icons/logo.png" alt="Logo">
                </div>
            </a>
            <nav class="top-bar" id="auth-nav">
            </nav>
        </header>
        `
    }
}