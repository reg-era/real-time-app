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