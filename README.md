# Vacuum-for-hire
- Récupère toutes les offres de job postées depuis le dernier run/chaque jour (https://github.com/robfig/cron)
- Pas de doublons si sur un même site - hash sur l’url ? 
- Site :
  - Google jobs (API : https://cloud.google.com/talent-solution/job-search/docs) -> récupérer le htidocid (https://cloud.google.com/talent-solution/job-search/docs)
  - Indeed API : https://developer.indeed.com/docs/authorization/3-legged-oauth
- Enregistrement quelque part + moyen facile d’y accéder (GoogleSheet ?)
- Stack : Go, API, Scrapping Selenium/Goquery, Docker, (Redis/SQLite ? titre/employeur/url/hash si nécessaire)
- Vérifier comment implémenter un cache
- Programme qui tourne constamment lancer des goroutines toutes les X heures

https://github.com/mattn/go-sqlite3

#### améliorations:
- serveur kimsufi (debian + docker)
- hébergement docker (https://geekflare.com/fr/docker-hosting-platforms/)
- pouvoir lancer manuellement avec un bouton