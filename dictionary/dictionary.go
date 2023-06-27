package dictionary

import (
	"errors"
	"fmt"
	"time"
)

// Défini la structure de Entry
type Entry struct {
	Definition string    `json:"definition"`
	Date       time.Time `json:date`
}

// CEtte fonction prend en parametre un element de type Entry et retour un string
func (e Entry) String() string {
	return fmt.Sprintf("Definition: %s\nDate: %s", e.Definition, e.Date.Format("2006-01-02 15:04:05"))
}

// Défini la structure du dictionnaire (struct)
// avec la clé entries et la valeur une map
// la map prend comme clé un string et une valeur Entry
type Dictionary struct {
	entries map[string]Entry
}

// Crée un dictionnaire de pointeur et retourne l'adresse du dictionnaire
func New() *Dictionary {
	return &Dictionary{
		entries: make(map[string]Entry),
	}
}

func (d *Dictionary) Add(word string, definition string) {
	entrie := Entry{
		Definition: definition,
		Date:       time.Now(),
	}
	d.entries[word] = entrie
}

func (d *Dictionary) Get(word string) (Entry, error) {

	entry, found := d.entries[word]
	if !found {
		return Entry{}, errors.New("Mot non trouvé dans le dictionnaire")
	}
	return entry, nil
}

func (d *Dictionary) Remove(word string) {
	delete(d.entries, word)
}

func (d *Dictionary) List() ([]string, map[string]Entry) {
	wordList := make([]string, 0, len(d.entries))
	for word := range d.entries {
		wordList = append(wordList, word)
	}

	return wordList, d.entries
}
