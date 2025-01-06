export function showRegistrationModal() {
    const dialog = document.createElement('dialog');
    dialog.innerHTML = `
        <h2 id="dialogTitle">Access Restricted</h2>
        <p id="dialogMessage">You need to be logged in to react. Please register or log in to continue.</p>
        <button class="modal-button register-btn">Register Now</button>
        <button class="modal-button login-btn">Login</button>
        <button class="modal-button close-btn" aria-label="Close dialog">X</button>
    `;

    const registerButton = dialog.querySelector('.register-btn');
    const loginButton = dialog.querySelector('.login-btn');
    const closeButton = dialog.querySelector('.close-btn');

    registerButton.addEventListener('click', () => {
        window.location.href = '/register';
    });

    loginButton.addEventListener('click', () => {
        window.location.href = '/login';
    });

    closeButton.addEventListener('click', () => {
        dialog.close();
    });

    dialog.addEventListener('click', (event) => {
        if (event.target === dialog) {
            dialog.close();
        }
    });

    document.body.appendChild(dialog);
    dialog.showModal();
}


export function handleResize() {
    const menuButton = document.querySelector('.menu-button');
    const sideBar = document.querySelector('.sidebar');
    const postContainer = document.querySelector('.posts');
    const createPost = document.querySelector('.create-post');

    if (menuButton) {
        if (window.innerWidth <= 1200) {
            if (window.location.pathname === '/login' || window.location.pathname === '/register') {
                return
            }
            menuButton.style.display = 'block';
            if (sideBar) {
                sideBar.classList.add('hide');
            }
            if (postContainer) {
                postContainer.style.marginLeft = '0';
            }
            if (createPost) {
                createPost.style.marginLeft = '0';
            }
        } else {
            menuButton.style.display = 'none';
            if (sideBar) {
                sideBar.classList.remove('hide');
            }
            if (postContainer) {
                postContainer.style.marginLeft = '250px';
            }
            if (createPost) {
                createPost.style.marginLeft = '250px';
            }
        }
    }
}

export function debounce(func, delay) {
    let timer;
    return function (...args) {
        clearTimeout(timer);
        timer = setTimeout(() => func.apply(this, args), delay);
    };
}