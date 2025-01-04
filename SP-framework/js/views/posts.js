

import Base from "./_BASE.js";

export default class extends Base {
    constructor(params) {
        super(params);
        this.setTitle("Posts");
    }

    async getHtml() {
        return `
            <h1>welcom to post</h1>
        `;
    }
}