function categoriesListListener() {
    const categories_list = document.querySelector(".category-list");
    categories_list.addEventListener('click', (event) => {
        if (event.target.tagName === 'LI' || 'A' || 'SPAN') {
            const id = event.target.closest('LI').getAttribute("categoryId")
            SubmitForm(id)
        }
    })
}

function SubmitForm(category) {
    const params = new URLSearchParams({ category: category });
    window.location.href = `/categories?${params}`;
}

categoriesListListener()