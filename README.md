"# Projet-Sujet-power-4-web"

Dans ce projet, j’ai développé l’ensemble de la logique du jeu Puissance 4 en langage Go. L’objectif était de créer un système capable de gérer une partie complète entre deux joueurs, depuis l’initialisation du plateau jusqu’à la détection d’un vainqueur ou d’une égalité.

Pour cela, j’ai mis en place une structure de données représentant la grille du jeu, les joueurs, les coups réalisés ainsi que l’état général de la partie. Le programme initialise la grille vide, attribue les informations aux joueurs et démarre la partie en enregistrant le moment de début.

J’ai implémenté un mécanisme permettant à chaque joueur de choisir une colonne pour y déposer un pion. Le programme vérifie que la colonne est valide et non remplie, puis place automatiquement le pion dans la case la plus basse disponible, comme dans les règles officielles du Puissance 4.

Après chaque coup, le programme analyse si ce mouvement entraîne une victoire. La détection se fait dans les quatre directions possibles : horizontalement, verticalement et dans les deux diagonales. Si un alignement de quatre pions est trouvé, la partie se termine et le vainqueur est enregistré. Dans le cas où la grille est remplie entièrement sans qu’aucun joueur n’aligne quatre jetons, le programme considère la partie comme une égalité.

J’ai également prévu un suivi des coups joués, avec des informations comme le joueur concerné, la colonne jouée, la ligne atteinte par le pion et l’heure précise du mouvement. La partie enregistre aussi son heure de fin dès qu’un résultat est déterminé.

Enfin, j’ai ajouté une fonctionnalité permettant d’afficher la grille à tout moment pendant la partie. L’affichage représente les cases vides et les pions des deux joueurs, ce qui permet de suivre visuellement l’état du jeu au fur et à mesure.

En résumé, mon code gère :

* la création et l’initialisation complète d’une partie ;

* le placement des pions selon les règles du jeu ;

* la gestion des tours des joueurs ;

* la détection automatique des victoires ;

* la gestion des égalités ;

* l’enregistrement des coups joués ;

* l’affichage de la grille du jeu ;

* l’état global de la partie (en cours, victoire, égalité).

" Tout cela permet de jouer une partie de Puissance 4 entièrement fonctionnelle via un programme écrit en Go "