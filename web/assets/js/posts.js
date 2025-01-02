import { initializeCommentSection } from "./comments.js";
import { getReactInfo, reactToggle } from "./likes.js";

export async function renderPosts(postsContainer, posts) {
  for (const post of posts) {
    const postElement = document.createElement("div");
    postElement.classList.add("post");

    try {
      const reactInfo = await getReactInfo({
        target_type: "post",
        target_id: post.PostId,
      }, "GET");

      postElement.innerHTML = generatePostHTML(post, reactInfo);
      postsContainer.appendChild(postElement);

      reactToggle(postElement, post.PostId, "post");
      initializeCommentSection(postElement, post);
    } catch (error) {
      console.error("Error rendering post:", error);
    }
  }
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

    <div class="reaction-section">
      <div class="reaction-buttons">
        <button class="like like-button ${liked ? "clicked" : ""}" data-clicked=${liked}>
          <span class="emoji">👍</span> Like (<span class="count">${likeCount}</span>)
        </button>
        <button class="dislike dislike-button ${disliked ? "clicked" : ""}" data-clicked=${disliked}>
          <span class="emoji">👎</span> Dislike (<span class="count">${disLikeCount}</span>)
        </button>
      </div>
      <button class="toggle-comments">💬 Show Comments</button>
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

export const isValidPost = (event) => {
  event.preventDefault()
  const checkbox = document.querySelectorAll('[name="category"]')
  let test = false
  checkbox.forEach((box) => {
    if (box.checked) test = true
  })
  if (!test) {
    document.getElementById('responseMessage').textContent = 'Oops! It looks like every post needs to have at least one category.'
    return
  }
  document.querySelector('#submition-button').disabled = true;
  event.target.submit()
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