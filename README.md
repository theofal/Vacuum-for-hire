# Vacuum-for-hire

* Récupère toutes les offres de job postées depuis le dernier run/chaque jour
* Pas de doublons si sur un même site - hash sur l’url ?
* Différents sites :
    * Welcome to the jungle (scrap)
    * Linkedin (scrap)
    * Google jobs ? (API : https://cloud.google.com/talent-solution/job-search/docs)
    * Smartrecruiters (API : https://developers.smartrecruiters.com/docs/the-smartrecruiters-platform)
    * Indeed (API : https://developer.indeed.com/)
    * ...
* Enregistrement quelque part + moyen facile d’y accéder
* Stack : Go, API, Scrapping (colly? TBD), Docker, (Redis/db ? titre/url/site/hash)
