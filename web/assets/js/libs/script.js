export function handleResize() {
    const menuButton = document.querySelector('.menu-button');
    const sideBar = document.querySelector('.sidebar-for-min');
    const postContainer = document.querySelector('.posts');
    const createPost = document.querySelector('.create-post');
    const onlinebar = document.querySelector('.onligne-bar');
    const bottonCreat = document.querySelector('.creat-post');

    if (window.location.pathname === '/login' || window.location.pathname === '/register') {
        bottonCreat.style.display = 'none';
        // the function handlesize is not called in login page and register
        return
    }
    if (menuButton) {
        if (window.innerWidth <= 900) {
            if (window.location.pathname === '/login' || window.location.pathname === '/register') {
                bottonCreat.style.display = 'none';
                // the function handlesize is not called in login page and register
                return
            }
            
            menuButton.style.display = 'block';
            onlinebar.style.display = 'none';
            bottonCreat.style.display = 'none';
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
            if (window.location.pathname === '/') {
                bottonCreat.style.display = 'block';
            }
            menuButton.style.display = 'none';
            onlinebar.style.display = 'block';

            if (sideBar) {
                if (!sideBar.classList.contains('hide')) {
                    sideBar.classList.add('hide');
                }
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