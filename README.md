# Vacuum-for-hire (over engineered)
- Récupère toutes les offres de job postées depuis le dernier run/chaque jour
- Site :
  - Google jobs
- Enregistrement dans une DB SQLITE
- Récupération des données depuis la db
- Mise en place d'une API hébergée temporairement en local (le temps de récupérer les données -> webhook/channel pour savoir quand arreter ?)
- Recuperation des données via l'API et enregistrement sur un fichier CSV
- Stack : Go, API, Scrapping Selenium/Goquery, Docker, SQlite3


#### améliorations:
- serveur kimsufi (debian + docker)
- hébergement docker (https://geekflare.com/fr/docker-hosting-platforms/)
- Mise en place d'un cron pour lancer une goroutine tous les X t - https://github.com/robfig/cron)
- pouvoir lancer manuellement avec un bouton (webhook)