import { BASE } from "./_BASE.js";

export class Error extends BASE {
    constructor(params) {
        super(params);

        this.statusError;
        this.statusMsg;
        this.errorMsg;

        this.setTitle("ERROR");
        this.setStyle("http://localhost:8080/api/css/error.css");
    }

    setAttribute() {
        console.log('try to set atribute', this.params);

        this.statusError = Number.parseInt(this.params['status']);
        switch (this.statusError) {
            case 404:
                this.statusMsg = "Page Not Found";
                this.errorMsg = "We couldn't find the page you're looking for. It might have been moved or deleted.";
                break;
            case 400:
                this.statusMsg = "Bad Request";
                this.errorMsg = "The request could not be understood by the server. Please check the URL or your input.";
                break;
            case 401:
                this.statusMsg = "Unauthorized";
                this.errorMsg = "You need to be logged in to access this page. Please log in to continue.";
                break;
            case 405:
                this.statusMsg = "Method Not Allowed";
                this.errorMsg = "Please Go to home.";
                break;
            case 500:
                this.statusMsg = "Internal Server Error";
                this.errorMsg = "Something went wrong on our end. We're working on it. Please try again later.";
                break;
        }
    }

    async getHtml() {
        this.setAttribute()
        const html = `
        ${this.getHtmlBase()}
        <main>
            <section class="container">
                <div class="error-message">
                    <h1>Oops! ${this.statusMsg} (${this.statusError})</h1>
                    <p>${this.errorMsg}</p>
                    <button class="err-button" onclick="window.location.href='/'">Go to Home</button>
                </div>
            </section>
        </main>
        `
        return html
    }
}