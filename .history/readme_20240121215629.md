Jeu du pendu
Le jeu du pendu est un simple jeu de console implémenté en Go. Il permet à 1 à 3 joueurs de deviner des lettres à tour de rôle pour révéler un mot caché. Les joueurs peuvent sélectionner le niveau de difficulté (facile, moyen ou difficile) et fixer une limite de temps pour chaque devinette. À la fin d'une partie, les scores sont affichés et les joueurs ont la possibilité de rejouer.

Comment jouer
Nombre de joueurs :

Saisissez le nombre de joueurs (1, 2 ou 3) lorsque vous y êtes invité.
Pour les parties multijoueurs, les joueurs devinent les lettres à tour de rôle.
Noms des joueurs :

Entrez les noms des joueurs.
Le score de chaque joueur est suivi tout au long du jeu.
Niveau de difficulté :

Sélectionnez le niveau de difficulté (facile, moyen ou difficile).
Le niveau de difficulté détermine la complexité du mot caché.
Deviner :

Les joueurs devinent les lettres à tour de rôle.
Saisissez une lettre lorsqu'on vous le demande ou tapez "indice" pour obtenir un indice.
Progression du jeu :

Le jeu affiche le joueur actuel, le mot deviné, les lettres devinées et les vies restantes.
Le jeu se termine si un joueur devine correctement le mot ou si les vies sont épuisées.
Score :

Les joueurs reçoivent des points pour les réponses correctes.
Les points sont déduits pour les devinettes incorrectes et les pénalités de temps.
Classement :

Après chaque partie, les scores finaux sont affichés.
Les 10 meilleurs scores sont affichés sur le tableau de classement.
Rejouer :

Permet de choisir si l'on veut jouer un autre tour.
Configuration
La liste des mots est lue à partir du fichier "words.txt", qui doit contenir des mots classés par ordre de difficulté.
Ajustez la limite de temps pour chaque devinette en définissant la variable guessTimeout en secondes.
Classement
Le classement est stocké dans le fichier "leaderboard.txt".
Il enregistre les noms des joueurs et leurs scores.
Les 10 meilleurs scores sont affichés après chaque partie.
Comment utiliser le jeu
Assurez-vous que Go est installé sur votre machine.

Clonez le dépôt et accédez au répertoire du projet.

Exécutez la commande suivante pour lancer le jeu du pendu :

go run main.go

Appréciez le jeu !
Amusez-vous à jouer au pendu avec vos amis ! N'hésitez pas à modifier le code ou à ajouter de nouvelles fonctionnalités pour améliorer le jeu. Si vous rencontrez des problèmes, vérifiez les messages d'erreur dans la console et assurez-vous que le fichier "words.txt" est correctement formaté.