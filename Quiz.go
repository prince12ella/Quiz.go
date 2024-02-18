package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)


// Question représente une question et ses réponses possibles
type Question struct {
	Text       string   // Texte de la question
	Choices    []string // Choix de réponses possibles
	CorrectAns string   // Réponse correcte
}

// quiz contient toutes les questions du jeu
var quiz []Question
var incorrectQuestion string
var correctAnswer string

func main() {
	// Charger les questions depuis un fichier texte
	if err := loadQuestions("questions.txt"); err != nil {
		log.Fatalf("Erreur lors du chargement des questions : %v", err)
	}

	// Démarrer le serveur web
	http.HandleFunc("/", quizHandler)
	http.HandleFunc("/correction", correctionHandler)
	fmt.Println("Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Fonction de gestionnaire pour la page principale du quiz
func quizHandler(w http.ResponseWriter, r *http.Request) {
	if len(quiz) == 0 {
		// Si aucune question n'est disponible, afficher un message d'erreur
		http.Error(w, "Aucune question disponible pour le moment. Veuillez réessayer plus tard.", http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodPost {
		// Obtenir la réponse de l'utilisateur depuis le formulaire
		userAns := r.FormValue("choice")

		// Sélectionner une question aléatoire
		rand.Seed(time.Now().UnixNano())
		question := quiz[rand.Intn(len(quiz))]

		// Vérifier si la réponse de l'utilisateur correspond à la réponse correcte
		if question.CorrectAns == userAns {
			// Rediriger vers la page principale pour la prochaine question
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		} else {
			// Stocker la question incorrecte et la réponse correcte
			incorrectQuestion = question.Text
			correctAnswer = question.CorrectAns
			// Afficher la correction
			http.Redirect(w, r, "/correction", http.StatusSeeOther)
			return
		}
	}

	// Sélectionner une question aléatoire
	rand.Seed(time.Now().UnixNano())
	question := quiz[rand.Intn(len(quiz))]

	// Afficher la question et les choix possibles
	// Génération du contenu HTML pour afficher la question et les choix possibles
	fmt.Fprintf(w, "<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n<meta charset=\"UTF-8\">\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n<title>bienvenue sur Quiz Game</title>\n<link rel=\"stylesheet\" href=\"/style.css\">\n</head>\n<body>\n<h1>bienvenue sur Quiz Game</h1>\n<div id=\"question-container\">\n<p id=\"question-text\">%s</p>\n<form id=\"choices-form\" method=\"post\">\n<div id=\"choices\">\n", question.Text)
	for _, choice := range question.Choices {
		// Ajout des choix de réponse sous forme de boutons radio dans le formulaire
		fmt.Fprintf(w, "<input type=\"radio\" name=\"choice\" value=\"%s\">%s<br>\n", choice, choice)
	}
	// Fermeture des balises HTML
	fmt.Fprint(w, "</div>\n<button type=\"submit\">Submit</button>\n</form>\n<p id=\"response\"></p>\n</div>\n\n</body>\n</html>")
}

// Fonction de gestionnaire pour la page de correction
func correctionHandler(w http.ResponseWriter, r *http.Request) {
	// Afficher la correction
	fmt.Fprintf(w, "<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n<meta charset=\"UTF-8\">\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n<title>Correction</title>\n<link rel=\"stylesheet\" href=\"/style.css\">\n</head>\n<body>\n<h1>Correction</h1>\n<div id=\"correction-container\">\n<p>Question incorrecte : %s</p>\n<p>Réponse correcte : %s</p>\n</div>\n</body>\n</html>", incorrectQuestion, correctAnswer)

	// Redirection vers la page principale après un court délai
	time.Sleep(3 * time.Second) // Attendre 3 secondes avant la redirection
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Fonction pour charger les questions à partir d'un fichier texte
func loadQuestions(filename string) error {
	// Ouvrir le fichier
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var questions []Question
	var q Question

	// Lire les questions ligne par ligne
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Q:") { // Texte de la question
			if q.Text != "" { // Si une question précédente est en cours de traitement, l'ajouter à la liste
				questions = append(questions, q)
				q = Question{} // Réinitialiser la question pour la prochaine itération
			}
			q.Text = strings.TrimSpace(strings.TrimPrefix(line, "Q:"))
		} else if strings.HasPrefix(line, "A:") { // Réponse correcte
			q.CorrectAns = strings.TrimSpace(strings.TrimPrefix(line, "A:"))
		} else if strings.HasPrefix(line, "C:") { // Choix de réponses
			choice := strings.TrimSpace(strings.TrimPrefix(line, "C:"))
			q.Choices = append(q.Choices, choice)
		}
	}

	// Ajouter la dernière question à la liste
	if q.Text != "" {
		questions = append(questions, q)
	}

	if len(questions) == 0 {
		return fmt.Errorf("aucune question trouvée dans le fichier")
	}

	// Assigner les questions chargées à la variable globale quiz
	quiz = questions

	return nil
}
