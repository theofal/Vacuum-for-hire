# Vacuum-for-hire
- Récupère toutes les offres de job postées depuis le dernier run/chaque jour
- Pas de doublons si sur un même site - hash sur l’url ? 
- Site :
  - Google jobs
- Enregistrement quelque part + moyen facile d’y accéder (GoogleSheet ? /!\ passage aux plans payants à partir du 1er juillet 2022 https://9to5google.com/2022/01/19/g-suite-legacy-free-edition/)
- Stack : Go, API, Scrapping Selenium/Goquery, Docker, SQlite3
- Vérifier - implémenter un cache ?

  https://github.com/mattn/go-sqlite3

#### améliorations:
- serveur kimsufi (debian + docker)
- hébergement docker (https://geekflare.com/fr/docker-hosting-platforms/)
- Mise en place d'un cron pour lancer une goroutine tous les X t - https://github.com/robfig/cron) 
- pouvoir lancer manuellement avec un bouton