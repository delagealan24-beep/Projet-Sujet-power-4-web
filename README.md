"# Projet-Sujet-power-4-web"

Les modifications apportées :
- J’ai refait tout mon main.go car il manquait énormément de choses.
- J’ai implanté toutes mes routes HTML dans mon code Go.
- J’ai refait mes templates en y apportant les modifications faites dans mon main.go.
- J’ai aussi ajouté un go.mod car je n’en avais pas, ce qui était une erreur de ma part.
- J’ai également ajouté un fichier de debug pour les potentiels problèmes.
Mais même avec ces modifications, je n’arrive pas à faire fonctionner le site web. Toutes mes routes sont mises, mais mon CSS ne s’affiche pas, aucun des boutons ne fonctionne et on ne peut pas changer de page.






Mathlouthi Chaima

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

Delage Alan

Dans ce projet, je me suis occupé de l’implémentation des routes ainsi que du code associé à chacune d’elles. J’ai également conçu l’ensemble de l’interface graphique du serveur.

L’objectif de cette partie était de créer une interface web intuitive et accessible, permettant de jouer au Puissance 4 à deux joueurs de manière fluide et agréable.
Tout d’abord, j’ai mis en place la route / (index), qui correspond à la page d’accueil du site. On y retrouve le logo du jeu, ainsi que plusieurs boutons utiles : les règles, le scoreboard, et bien sûr le bouton “Jouer”.

Ensuite, j’ai développé la route /game/init, qui permet d’initialiser une partie. Les joueurs peuvent y saisir leurs noms ou pseudos, choisir leurs couleurs respectives, puis lancer la partie via le bouton “Commencer la partie”.

La route /game/play correspond à la page principale du jeu. Elle affiche la grille de Puissance 4, les flèches interactives pour choisir une colonne, ainsi qu’un message indiquant clairement quel joueur doit jouer. On y retrouve également des boutons pour revenir à l’accueil, consulter les règles ou accéder au scoreboard à tout moment.

J’ai également créé la route /game/end, qui s’affiche à la fin d’une partie. Elle indique le nom du vainqueur, le nombre de tours joués, la date de la partie, les noms des deux joueurs et leurs couleurs respectives. En cas d’égalité, le message est adapté pour le signaler.

Enfin, la route /game/scoreboard permet de consulter l’historique des parties jouées. Elle affiche les noms des joueurs, la date de chaque partie, le vainqueur ou l’éventuelle égalité.

En résumé, mes code gère :

🏠 Route /index – Page d’accueil
- Affiche la page d’accueil avec le logo et les boutons :
- “Jouer”
- “Règles”
- “Scoreboard”

🧑‍🤝‍🧑 Route /game/init – Initialisation de la partie
- Affiche un formulaire pour :
- Entrer les noms ou pseudos des deux joueurs
- Choisir leurs couleurs (rouge, jaune, etc.)
- Lancer la partie via le bouton “Commencer la partie”

🧠 Route /game/init/traitement – Traitement du formulaire
- Récupère les données du formulaire via POST
- Initialise la structure GameState avec :
- Les noms des joueurs
- Les couleurs choisies
- Le tour initial (1)
- Une grille vide (6x7 cases)
- Redirige vers la route /game/play

🎮 Route /game/play – Page de jeu
- Affiche :
- La grille du Puissance 4 avec les jetons placés
- Le joueur actif et sa couleur
- Un message indiquant à qui c’est le tour
- Un formulaire pour choisir une colonne (0 à 6)
- Des boutons pour naviguer : accueil, règles, scoreboard

🟡 Route /game/play/move – Traitement d’un coup
- Récupère la colonne choisie via GET
- Vérifie si la colonne est valide et non pleine
- Place le jeton du joueur actif dans la première case libre de la colonne
- Incrémente le tour
- Redirige vers /game/play pour afficher la grille mise à jour

🏁 Route /game/end – Fin de partie (à ajouter)
- Affiche :
- Le nom du vainqueur ou un message d’égalité
- Le nombre de tours joués
- La date de la partie
- Les noms et couleurs des joueurs
- Un bouton “Rejouer” pour relancer une partie
- Un lien vers le scoreboard

📊 Route /game/scoreboard – Historique des parties (à ajouter)
- Affiche toutes les parties jouées :
- Noms des joueurs
- Date
- Vainqueur ou égalité
- Utilise une variable globale pour stocker l’historique




  MATHLOUTHI CHAIMA

  
   Partie 3 — Power'4 Web en Go


Dans cette troisième partie du projet, j’ai développé toute 
la **version web du jeu Puissance 4**, en utilisant Go pour la 
gestion du serveur, les routes et le moteur de templates.

Cette partie se concentre sur :
  • l’affichage du jeu dans le navigateur  
  • l’interaction (clics sur les colonnes)  
  • la logique serveur (tour, grille, victoire…)  
  • l’interface, les fichiers HTML et le style  


 OBJECTIFS DE MA PARTIE 3


Pour cette partie, j’ai réalisé :

* Un serveur HTTP en Go  
* Un système de templates (layout + pages dynamiques)  
* Un formulaire de création de partie (noms + symboles)  
* La gestion du tour par tour côté serveur  
* L’affichage de la grille dans le navigateur  
* La détection de victoire et d’égalité  
* Une page de fin de partie  
* Un scoreboard mémorisant les gagnants  
* Le style complet du site (CSS simple et propre)  
* La réinitialisation de la partie  


 CE QUE J’AI AJOUTÉ DANS LA PARTIE 3


 *Templates HTML*
  - layout.gohtml : template principal  
  - homepage.gohtml : page d’accueil  
  - game_init.gohtml : formulaire de création  
  - game_play.gohtml : page de jeu avec la grille  
  - game_end.gohtml : fin de partie  
  - scoreboard.gohtml : liste des victoires  

*Fichiers statiques*
  - style.css : mon design (couleurs, cartes, grille, boutons)

 *Routes Go*
  - `/`                    → accueil  
  - `/game/init`           → création d’une partie  
  - `/game/play`           → interface du jeu  
  - `/game/play/move?col=` → poser un jeton  
  - `/game/end`            → fin de partie  
  - `/game/scoreboard`     → historique des gagnants  
  - `/game/reset`          → remettre à zéro  

*Logique*
  - placeToken() : poser un jeton par colonne  
  - checkVictory() : détecter les 4 alignés  
  - isDraw() : match nul  
  - Gestion du tour (J1/J2)  
  - Stockage des gagnants en mémoire  
  - Protection avec Mutex  


✔ BILAN DE MA PARTIE 3


Dans cette partie 3, j’ai transformé le jeu en une vraie 
application web complète.  
J’ai appris à structurer des templates, gérer un serveur Go, 
transmettre des données aux pages, manipuler une grille en 
HTML, gérer les interactions et produire une interface 
entièrement fonctionnelle.
