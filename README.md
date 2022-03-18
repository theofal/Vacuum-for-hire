# Vacuum-for-hire

* Récupère toutes les offres de job postées depuis le dernier run/chaque jour (https://github.com/robfig/cron)
* Pas de doublons si sur un même site - hash sur l’url ?
* Site :
    * Google jobs (API : https://cloud.google.com/talent-solution/job-search/docs) -> récupérer le htidocid (https://cloud.google.com/talent-solution/job-search/docs)
* Enregistrement quelque part + moyen facile d’y accéder (GoogleSheet ?)
* Stack : Go, API, Scrapping (colly? TBD), Docker, (Redis/SQLite ? titre/employeur/url/hash si nécessaire)
* Programme qui tourne constamment lancer des goroutines toutes les X heures

améliorations: 
* serveur kimsufi (debian + docker) 
* hébergement docker (https://geekflare.com/fr/docker-hosting-platforms/) 
* pouvoir lancer manuellement avec un bouton
