import BASE from "./_BASE.js";

export default class extends BASE {
    constructor(params) {
        super(params);
        this.setTitle("Home");
        this.setStyle("http://localhost:8080/assets/css/new_post.css")
    }

    async getCategories() {
        //fetch to get categories
    }

    async getHtml() {
        // call the categories
        return `
        ${this.getHtmlBase()}
        ${this.getSideBar()}
        <section class="create-post">
            <h2>Create a New Post</h2>
            <form id="createPostForm" action="/new_post" method="post">
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
        `;
    }
}