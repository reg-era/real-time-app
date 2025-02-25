import { initializeCommentSection } from "./comment.js";
import { reactToggle, getReactInfo } from "./likes.js";

export const GetData = async () => {
    let postIds = [];
    try {
        let response = await fetch('api/posts');
        if (!response.ok) throw new Error("Network response was not ok");
        let lastPostId = await response.json();
        for (let postId = 1; postId <= lastPostId; postId++) {
            postIds.push(postId);
        }
        return postIds
    } catch (err) {
        console.error(err);
    }
};

export async function renderPage(postIds) {
    try {
        let targets = [];
        let i = 0
        while (postIds.length > 0 && i < 10) {
            const link = `/api/posts?post_id=${postIds.pop()}`;
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
        let data = null;
        try {
            data = await renderPosts(targets);
        } catch (err) {
            return err
        }
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
            const reactInfo = await getReactInfo({
                target_type: "post",
                target_id: post.PostId,
            }, "GET");

            postElement.innerHTML = generatePostHTML(post, reactInfo);
            res.push(postElement)
            reactToggle(postElement, post.PostId, "post");
            initializeCommentSection(postElement, post);
        } catch (error) {
            console.error("Error rendering post:", error);
            return error
        }
    }
    return res
}


function generatePostHTML(post, reactInfo) {
    let liked = false;
    let disliked = false;

    let likeCount = reactInfo.data.liked_by ? reactInfo.data.liked_by.length : 0;
    let disLikeCount = reactInfo.data.disliked_by ? reactInfo.data.disliked_by.length : 0;

    if (!!reactInfo.data.user_reaction) {
        liked = reactInfo.data.user_reaction === "like"
        disliked = !liked
    } else {
        liked = false
        disliked = false;
    }

    return `
    <div class="post-container">
        <div class="post-header">
            <div class="post-meta">
            <img class="profile" src="https://avatar.iran.liara.run/public/${post.Gender === "male" ? "girl" : "boy"}?username${post.UserName}" loading="lazy">
                <span class="author">${post.UserName}</span>
                <span class="date">${new Date(post.CreatedAt).toLocaleString()}</span>
                <br>
                <span class="categories">${post.Categories || "Not categorized"}</span>
            </div>
        <h1 class="post-title">${escapeHTML(post.Title)}</h1>
        </div>
    
        <div class="post-body">
            <p class="content">${escapeHTML(post.Content)}</p>
        </div>
    
        <div class="reaction-section">
            <div class="reaction-buttons">
            <button class="like like-button ${liked ? "clicked" : ""}" data-clicked=${liked}>
                <span class="like-emoji"> </span> (<span class="count">${likeCount}</span>)
            </button>
            <button class="dislike dislike-button ${disliked ? "clicked" : ""}" data-clicked=${disliked}>
                <span class="dislike-emoji"> </span> (<span class="count">${disLikeCount}</span>)
            </button>
            </div>
            <button class="toggle-comments">💬</button>
        </div>
    
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