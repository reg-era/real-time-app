package database

type Category struct {
	Id          int    `json:"Id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var CategoriesTable = `
    CREATE TABLE IF NOT EXISTS categories (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL UNIQUE,
        description TEXT NOT NULL
    );
    INSERT OR IGNORE INTO categories (name, description) VALUES
        ('Music', 'Discuss everything related to music, including genres, artists, and concerts'), 
        ('Sports', 'Talk about all types of sports, games, and tournaments'), 
        ('Movies & TV Shows', 'Share recommendations and discuss your favorite films and series'), 
        ('Technology', 'Discuss the latest trends in tech, gadgets, and software'), 
        ('Gaming', 'A place for gamers to discuss games, consoles, and tips'), 
        ('Books & Literature', 'Share and discover books, authors, and literary genres'), 
        ('Travel', 'Exchange travel tips, favorite destinations, and experiences'), 
        ('Food & Cooking', 'Discuss recipes, restaurants, and all things culinary');
`

var PostCategoriesTable = `
    CREATE TABLE IF NOT EXISTS post_categories (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        post_id INTEGER NOT NULL,
        category_id INTEGER NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
        FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
        FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
    );
`
