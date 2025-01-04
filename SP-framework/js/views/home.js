import BASE from "./_BASE.js";

export default class extends BASE {
    constructor(params) {
        super(params);
        this.setTitle("Home");
    }

    async getHtml() {
        return `
            <h1>hello word</h1>
        `;
    }
}