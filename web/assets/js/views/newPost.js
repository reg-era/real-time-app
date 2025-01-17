import { BASE } from "./_BASE.js";

export class NewPost extends BASE {
    constructor(params) {
        super(params);
        this.setTitle("New");
        this.setStyle("http://localhost:8080/api/css/base.css");
        this.setStyle("http://localhost:8080/api/css/newPost.css");
        this.listenerSet = false;
    }

    setAttribute() {
    }

    setListners() {
        if (this.listenerSet) return; // just to be sure the lisner i call one time
        this.listenerSet = true;

        // setup textarea resizing
        const textarea = document.querySelector('.post-content');
        textarea.addEventListener('input', (event) => {
            textarea.style.height = 'auto';
            textarea.style.height = textarea.scrollHeight + "px";
        });

        // setup post form submission listener
        const button = document.querySelector('#submition-button');
        button.disabled = true;

        // validation posts and submit them to backend
        document.getElementById('createPostForm').addEventListener('submit', async (e) => {
            e.preventDefault();

            const checkbox = document.querySelectorAll('[name="category"]');
            let test = false;
            checkbox.forEach((box) => {
                if (box.checked) test = true;
            });

            if (!test) {
                document.getElementById('responseMessage').textContent = 'Oops! It looks like every post needs to have at least one category.';
                return;
            }

            document.querySelector('#submition-button').disabled = true;
            const formData = new FormData(createPostForm);

            if (!createPostForm.checkValidity()) {
                responseMessage.textContent = 'Please fill out all required fields.';
                return;
            }

            const res = await fetch(createPostForm.action, {
                method: 'POST',
                body: formData,
            });

            if (!res.ok) {
                responseMessage.textContent = 'An unexpected error occurred.';
            }

            window.location.href = '/';
        });
    }

    async getCategories() {
        try {
            const res = await fetch('http://localhost:8080/api/new_post');
            if (res.ok) {
                const categories = await res.json();
                let html = '';
                categories.forEach(obj => {
                    html += `<label><input type="checkbox" name="category" value="${obj.Id}"> ${obj.name}</label>`;
                });
                document.querySelector('.categories').innerHTML = html;
            } else {
                window.location.href = `/error?status=${res.status}`;
            }
        } catch (err) {
            console.error(err);
        }
    }

    async getHtml() {
        const html = `
        ${this.getHtmlBase()}
        <main>
            <section class="create-post">
                <h2>Create a New Post</h2>
                <form id="createPostForm" action="/api/new_post" method="post">
                    <input name="title" type="text" placeholder="Post Title" class="post-title" minlength="3" maxlength="60" required>
                    <textarea name="content" placeholder="Post Content" class="post-content" required minlength="10" maxlength="10000"></textarea>
                    <div class="form-group">
                        <label>Choose Categories:</label>
                        <div class="categories">
                        </div>
                    </div>
                    <button id="submition-button" type="submit">Create Post</button>
                    <p id="responseMessage"></p>
                </form>
            </section>
        </main>
        `;

        this.getCategories();
        setTimeout(() => this.setListners(), 0);

        return html;
    }
}
