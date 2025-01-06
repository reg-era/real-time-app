import { showRegistrationModal } from "./script.js";
import { escapeHTML } from "./post.js";

const commentSize = 3
const comentIndex = {}

export const initializeCommentSection = (postElement, post) => {
    const toggleCommentsButton = postElement.querySelector(".toggle-comments");
    const commentsSection = postElement.querySelector(".comments-section");
    const showMore = postElement.querySelector(".more-comment");
    const hidebotton = postElement.querySelector('.hide-comments');

    toggleCommentsButton.addEventListener("click", async () => {
        if (commentsSection.style.display === "none") {
            commentsSection.style.display = "block"
            const comment = commentsSection.querySelector(".comments");
            const index = comment.querySelectorAll('.comment');
            comentIndex[post.PostId] = index.length;
            if (index.length === 0) await loadComments(post.PostId, commentSize, commentsSection.querySelector(".comments"))
            toggleCommentsButton.style.display = "none";
        }
    });

    hidebotton.addEventListener("click", () => {
        commentsSection.style.display = "none";
        toggleCommentsButton.style.display = "block";
    })

    const commentInput = postElement.querySelector(".comment-input");
    commentInput.addEventListener("keydown", async (event) => {
        if (event.key === "Enter" && !event.shiftKey) {
            const commentInput = postElement.querySelector(".comment-input");
            if (commentInput.value.trim()) {
                await addComment(post.PostId, commentInput.value.trim(), commentsSection.querySelector(".comments"), commentsSection);
            }
        }
    })

    showMore.addEventListener('click', async () => {
        const comment = commentsSection.querySelector(".comments");
        const index = comment.querySelectorAll('.comment');
        comentIndex[post.PostId] = index.length;
        await loadComments(post.PostId, commentSize, commentsSection.querySelector(".comments"))
    })

    commentInput.addEventListener('input', () => {
        commentInput.style.height = 'auto'
        commentInput.style.height = commentInput.scrollHeight + "px"
    });
}

const loadComments = async (postId, limit, commentsContainer) => {
    try {
        const response = await fetch(`http://localhost:8080/comments?post=${postId}&limit=${limit}&from=${comentIndex[postId]}`)
        if (!response.ok) throw new Error("Failed to load comments.")

        const comments = await response.json()
        if (!comments || comments.length === 0) return

        let count = 0
        for (const comment of comments) {
            const commentSection = createCommentElement(comment)
            commentsContainer.appendChild(commentSection)
            count++
        }
        comentIndex[postId] += count

        if (count < limit) return
        await loadComments(postId, commentSize - count, commentsContainer)
    } catch (error) {
        console.error("Error loading comments:", error);
    }
}

const addComment = async (postId, content, commentsContainer, commentsection) => {
    try {
        const response = await fetch(`http://localhost:8080/comments`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
                post_id: postId,
                content: content
            }),
        })

        let commentInput = commentsection.querySelector(".comment-input");
        const error = commentsection.querySelector('.error-comment')
        const newComment = await response.json();

        switch (response.status) {
            case 400:
                error.textContent = newComment.error
                break
            case 401:
                showRegistrationModal()
                break
            case 201:
                error.textContent = ""
                const commentSection = createCommentElement(newComment)
                commentsContainer.prepend(commentSection)
                commentInput.style.height = "38px"
                commentInput.value = ""
                break
        }
    } catch (error) {
        console.error("Error adding comment:", error)
    }
}

const createCommentElement = (comment) => {
    const commentElement = document.createElement("div")
    commentElement.classList.add("comment")

    commentElement.innerHTML = `
    <p><strong>👤 ${comment.user_name}:</strong> ${escapeHTML(comment.content)}</p>
    `
    return commentElement
}
