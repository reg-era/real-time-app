import { renderPosts } from "./posts.js";

export const GetData = async (postIds = false) => {
    if (postIds == null) {
        return;
    }

    const postsContainer = document.querySelector(".posts");
    postsContainer.innerHTML = "";

    try {
        if (postIds === false) {
            postIds = [];
            let response = await fetch('http://localhost:8080/posts');
            if (!response.ok) throw new Error("Network response was not ok");
            let lastPostId = await response.json();
            for (let postId = 1; postId <= lastPostId; postId++) {
                postIds.push(postId);
            }
        }

        renderPage(postIds, postsContainer);
        const debouncedRenderPage = debounce(renderPage, 1000)

        window.addEventListener('scroll', () => {
            const scrollPosition = window.scrollY;
            const documentHeight = document.documentElement.scrollHeight;
            const windowHeight = window.innerHeight;
            if (scrollPosition + windowHeight >= documentHeight - 10) {
                debouncedRenderPage(postIds, postsContainer)
            }
        });
    } catch (err) {
        console.error(err);
    }
};

function debounce(func, delay) {
    let timer;
    return function (...args) {
        clearTimeout(timer);
        timer = setTimeout(() => func.apply(this, args), delay);
    };
}


async function renderPage(postIds, postsContainer) {
    let target = [];
    let i = 0
    while (postIds.length > 0 && i < 10) {
        let link = `http://localhost:8080/posts?post_id=${postIds.pop()}`;
        let postResponse = await fetch(link);
        if (postResponse.ok) {
            let post = await postResponse.json();
            target.push(post);
        } else {
            if (postResponse.status !== 404) {
                throw new Error("Response not ok");
            }
        }
        i++
    }
    await renderPosts(postsContainer, target);
}

// Logout event
export const logoutEvent = (log) => {
    log.addEventListener('click', async (event) => {
        event.preventDefault()
        try {
            const response = await fetch('http://localhost:8080/logout', {
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
};


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

const authNav = document.getElementById('auth-nav');
const hasSession = document.cookie.includes('session_token');

if (hasSession) {
    authNav.innerHTML = `
        <a href="/" class="active">Logout</a>
    `;
    const logoutLink = authNav.querySelector('a');
    logoutEvent(logoutLink);
} else {
    authNav.innerHTML = `
        <a href="/login" class="active">Login</a>
        <a href="/register">Signup</a>
    `;
}

function handleResize() {
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

handleResize();

let debouncedHandleResize = debounce(handleResize, 100);
window.addEventListener('resize', debouncedHandleResize);

const menuButton = document.querySelector('.menu-button');
const sideBar = document.querySelector('.sidebar');
if (menuButton && sideBar) {
    menuButton.addEventListener('click', () => {
        sideBar.classList.toggle('hide');
    });
}
