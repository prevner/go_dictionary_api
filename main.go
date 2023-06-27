package main

import (
	"bufio"
	"fmt"
	"go_dictionary_api/dictionary"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/arriqaaq/flashdb"
	"github.com/gin-gonic/gin"
)

func main() {
	// Creation du dictionnaire
	dico := dictionary.New()

	// Ajout de la base de donné
	config := &flashdb.Config{Path: "./db", EvictionInterval: 10}
	db, err := flashdb.New(config)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// Initialisation du Router
	router := gin.Default()

	// methode cli pour le CRUD
	go cli(dico)

	// List
	// Route pour obtenir la liste des mots et de leurs définitions
	router.GET("/list", func(c *gin.Context) {

		err := db.View(func(tx *flashdb.Tx) error {
			val, err := tx.Get("java")
			if err != nil {
				return err
			}
			c.JSON(http.StatusOK, gin.H{
				"mot":        "java",
				"définition": val,
			})
			return nil
		})

		if err != nil {
			log.Println("Erreur lors de la récuperation du mot dans flashDB:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récuperation du mot"})
			return
		}
	})

	// Defintion
	router.GET("word/:name", func(c *gin.Context) {
		/* word := c.Params.ByName("name")
		entrie, _ := dico.Get(word) */

		err := db.View(func(tx *flashdb.Tx) error {
			word := c.Params.ByName("name")
			val, err := tx.Get(word)
			if err != nil {
				return err
			}
			c.JSON(http.StatusOK, gin.H{
				"mot":        word,
				"définition": val,
			})
			return nil
		})

		if err != nil {
			log.Println("Erreur lors de la recuperation du mot dans flashDB:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de l'ajout du mot"})
			return
		}
	})

	// ajout de mot et défintion
	router.POST("word", func(c *gin.Context) {
		mot := c.PostForm("mot")
		definition := c.PostForm("definition")

		// Ajouter le mot et la definition dans la map
		err := db.Update(func(tx *flashdb.Tx) error {
			err := tx.Set(mot, definition)
			return err
		})

		if err != nil {
			log.Println("Erreur lors de l'ajout du mot dans flashDB:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de l'ajout du mot"})
			return
		}

		// Ajouter le mot et la définition dans le dictionnaire
		dico.Add(mot, definition)

		c.JSON(http.StatusOK, gin.H{"message": "Mot enregistré avec succès"})
	})

	// Supprimer un mot
	router.DELETE("delete/:name", func(c *gin.Context) {
		word := c.Params.ByName("name")
		// Supprimer le mot
		//dico.Remove(word)

		err := db.Update(func(tx *flashdb.Tx) error {
			word := c.Params.ByName(word)
			tx.Delete(word)
			if err != nil {
				return err
			}
			c.JSON(http.StatusOK, gin.H{
				"msg": "Mot supprimé ",
			})
			return nil
		})

		if err != nil {
			log.Println("Erreur lors de la suppression du mot dans flashDB:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression du mot"})
			return
		}

	})

	// Mise a jour un mot
	router.POST("update/:name", func(c *gin.Context) {
		word := c.Params.ByName("name")
		definition := c.PostForm("definition")

		err := db.Update(func(tx *flashdb.Tx) error {
			tx.Set(word, definition)
			if err != nil {
				return err
			}
			c.JSON(http.StatusOK, gin.H{
				"msg": "Mot mis à jour ",
			})
			return nil
		})

		if err != nil {
			log.Println("Erreur lors de la mise a jour du mot dans flashDB:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise a jour du mot"})
			return
		}

	})

	router.Run()
}

func cli(dico *dictionary.Dictionary) {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter ADD, DEF, REMOVE, LIST, EXIT : ")
		arg, _ := reader.ReadString('\n')

		word := strings.TrimSpace(arg)

		switch word {
		case "ADD":
			actionAdd(dico, reader)
		case "DEF":
			actionDefine(dico, reader)
		case "REMOVE":
			actionRemove(dico, reader)
		case "LIST":
			actionList(dico)
		case "EXIT":
			return

		}
	}
}

func actionAdd(d *dictionary.Dictionary, reader *bufio.Reader) {
	// Inserer le mot
	fmt.Print("Entrer le mot: ")
	mot, _ := reader.ReadString('\n')
	mot = strings.TrimSpace(mot)
	// Inserer la definition
	fmt.Print("Entrer la définition: ")
	definition, _ := reader.ReadString('\n')
	definition = strings.TrimSpace(definition)
	// Ajouter le mot et la definition dans la map
	d.Add(mot, definition)
}

func actionDefine(d *dictionary.Dictionary, reader *bufio.Reader) {
	fmt.Print("Entrer le mot: ")
	mot, _ := reader.ReadString('\n')
	mot = strings.TrimSpace(mot)
	entrie, _ := d.Get(mot)
	fmt.Println("  Définition:", entrie)
}

func actionRemove(d *dictionary.Dictionary, reader *bufio.Reader) {
	// Inserer le mot
	fmt.Print("Entrer le mot: ")
	mot, _ := reader.ReadString('\n')
	mot = strings.TrimSpace(mot)
	// Supprimer le mot
	d.Remove(mot)
}

func actionList(d *dictionary.Dictionary) {

	words, entries := d.List()

	fmt.Println("Mots dans le dictionnaire:")

	for _, word := range words {

		entry := entries[word]

		fmt.Println("- Mot:", word)

		fmt.Println("- Définition:", entry)

	}

}
