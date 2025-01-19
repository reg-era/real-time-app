import { initializeCommentSection } from "./comment.js";

export const GetData = async (postIds = false) => {
    if (postIds == null) {
        return;
    }

    try {
        if (postIds === false) {
            postIds = [];
            let response = await fetch('api/posts');
            if (!response.ok) throw new Error("Network response was not ok");
            let lastPostId = await response.json();
            for (let postId = 1; postId <= lastPostId; postId++) {
                postIds.push(postId);
            }
        }

        const body = await renderPage(postIds);
        return body
        // const debouncedRenderPage = debounce(renderPage, 1000)
        // 
        // window.addEventListener('scroll', () => {
        // const scrollPosition = window.scrollY;
        // const documentHeight = document.documentElement.scrollHeight;
        // const windowHeight = window.innerHeight;
        // if (scrollPosition + windowHeight >= documentHeight - 10) {
        // debouncedRenderPage(postIds, postsContainer)
        // }
        // });
    } catch (err) {
        console.error(err);
    }
};

export async function renderPage(postIds) {
    try {
        let targets = [];
        let i = 0
        while (postIds.length > 0 && i < 10) {            
            const link = `http://localhost:8080/api/posts?post_id=${postIds.pop()}`;
            const postResponse = await fetch(link);
            if (postResponse.ok) {
                const post = await postResponse.json();
                targets.push(post);
            } else {
                if (postResponse.status !== 404) {
                    throw new Error("Response not ok");
                }
            }
            i++
        }
        const data = await renderPosts(targets);
        return data;
    } catch (error) {
        console.error(error);
    }
}

export async function renderPosts(posts) {
    const res = []
    for (const post of posts) {
        const postElement = document.createElement("div");
        postElement.classList.add("post");
        try {
            postElement.innerHTML = generatePostHTML(post/*, reactInfo*/);
            res.push(postElement)
            initializeCommentSection(postElement, post);
        } catch (error) {
            console.error("Error rendering post:", error);
        }
    }
    return res
}


function generatePostHTML(post) {
    return `
    <div class="post-container">
        <div class="post-header">
            <div class="post-meta">
                <span class="author">👤  ${post.UserName}</span>
                <span class="date">${new Date(post.CreatedAt).toLocaleString()}</span>
                <br>
                <span class="categories">${post.Categories || "Not categorized"}</span>
            </div>
            <h1 class="post-title">${escapeHTML(post.Title)}</h1>
        </div>

        <div class="post-body">
            <p class="content">${escapeHTML(post.Content)}</p>
        </div>

        <button class="toggle-comments">💬 Show Comments</button>

        <div class="comments-section" style="display: none;">
            <div class="comments">
            </div>
            <div calss="comment-controllers">
                <button class="more-comment">Show more</button>
                <button class="hide-comments">hide comments</button>
            </div>
            <div class="comment-input-wrapper">
                <textarea required maxlength="2000" placeholder="Add a comment..." class="comment-input"></textarea>
            </div>
            <p class="error-comment"></p>
        </div>
    </div>
    `;
}

export function escapeHTML(input) {
    const map = {
        '&': '&amp;',
        '<': '&lt;',
        '>': '&gt;',
        '"': '&quot;',
        "'": '&#39;',
    };
    return input.replace(/[&<>"']/g, (char) => map[char]);
}