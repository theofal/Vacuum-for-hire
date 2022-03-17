# Vacuum-for-hire

* Récupère toutes les offres de job postées depuis le dernier run/chaque jour
* Pas de doublons si sur un même site - hash sur l’url ?
* Site :
    * Google jobs ? (API : https://cloud.google.com/talent-solution/job-search/docs) -> récupérer le htidocid
* Enregistrement quelque part + moyen facile d’y accéder (GoogleSheet ?)
* Stack : Go, API, Scrapping (colly? TBD), Docker, (Redis/SQLite ? titre/employeur/url/hash)

Serveur web qui tourne constamment 
lancer des goroutines toutes les X heures


ameliorations:
serveur kimsufi (debian + docker)
hebergement docker (https://geekflare.com/fr/docker-hosting-platforms/)
pouvoir lancer manuellement avec un bouton
