package data

type Book struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Example data, to be replaced with db
var books = []*Book{
	{
		ID:   1,
		Name: "First book",
	},
}

// Replace with DB logic
func AddBook(b *Book) {
	if len(books) == 0 {
		b.ID = 1
		books = append(books, b)
		return
	}

	maxID := books[len(books)-1].ID
	b.ID = maxID + 1
	books = append(books, b)

	return
}

func Books() []*Book {
	return books
}

func BookByID(id int) (*Book, error) {
	for _, b := range books {
		if b.ID == id {
			return b, nil
		}
	}

	return nil, ErrNotFound
}

func UpdateBook(id int, b *Book) error {
	for i, v := range books {
		if id == v.ID {
			b.ID = id
			books[i] = b
			return nil
		}
	}

	return ErrNotFound
}

func DeleteBook(id int) error {
	for i, b := range books {
		if b.ID == id {
			books = append(books[:i], books[i+1:]...)
			return nil
		}
	}

	return ErrNotFound
}
