import BASE from "./_BASE.js";

export default class extends BASE {
    constructor(params) {
        super(params);
        this.setTitle("Home");
        this.setStyle("http://localhost:8080/assets/css/login.css")
    }

    async getHtml() {
        return `
        ${this.getHtmlBase()}
        ${this.getNavigation()}
        `;
    }
}