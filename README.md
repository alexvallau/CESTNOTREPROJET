

Générateur Automatique de Dashboards Grafana pour Capteurs IoT

Introduction

Ce projet est un script écrit en Go qui permet de générer automatiquement des tableaux de bord Grafana pour le monitoring de capteurs IoT. Il s’inscrit dans un trio technologique comprenant Node-RED, InfluxDB et Grafana. L’idée générale est la suivante : Node-RED collecte les données de capteurs et les stocke dans InfluxDB, puis Grafana visualise ces données via des dashboards. Plutôt que de créer manuellement les dashboards pour chaque nouveau capteur, ce script Go automatise leur création. L’utilisateur indique le nombre de capteurs (devices) à superviser via un formulaire web, et le système génère dynamiquement le tableau de bord Grafana contenant les graphiques correspondants. Ce processus facilite le déploiement rapide de la supervision pour un nombre variable de capteurs, sans intervention manuelle dans Grafana.

Prérequis Techniques

Avant de déployer ce projet, assurez-vous de disposer de l’environnement suivant :
	•	Go (version 1.XX ou ultérieure) – pour compiler et exécuter le script Go.
	•	Grafana (idéalement version 8 ou 9) – pour l’affichage des dashboards. Grafana doit être accessible via son API HTTP et disposer d’une datasource InfluxDB configurée.
	•	InfluxDB (v2 ou supérieur, avec API Flux) – base de données time-series où sont stockées les mesures des capteurs. Le bucket de données IoT doit être configuré (par exemple iot-platform) et relié à Grafana via une datasource.
	•	Node-RED – pour la collecte des données capteur et l’envoi vers InfluxDB. Node-RED n’interagit pas directement avec le script Go, mais il alimente InfluxDB en données (mesures) que les dashboards Grafana afficheront.
	•	Token API Grafana – une clé API Grafana avec les droits nécessaires (édition de dashboards) est requise pour que le script puisse créer/modifier des dashboards via l’API. Ce token doit être configuré dans le script (ou dans l’environnement) avant exécution.
	•	Accès réseau – le script doit pouvoir communiquer avec l’instance Grafana (URL et port de Grafana accessibles, par ex. http://localhost:3000 ou une IP interne) et Grafana doit pouvoir accéder à InfluxDB.

Déploiement et Exécution du Script

Voici les étapes pour déployer et exécuter le générateur de dashboards dans le cadre du site web :
	1.	Récupération du code – Clonez le dépôt du projet sur le serveur web ou la machine où il sera exécuté. Assurez-vous d’avoir le fichier main.go du projet.
	2.	Configuration – Mettez à jour les paramètres si nécessaire : dans le code main.go, vérifiez l’URL de Grafana et le token API Grafana. Par défaut, le code utilise l’URL http://192.168.141.122:3000 et un token codé en dur – adaptez ces valeurs à votre environnement Grafana (par exemple http://localhost:3000 et votre propre jeton d’API). Assurez-vous également que le nom de la datasource Grafana (ici InfluxDB) correspond bien à celle configurée dans Grafana pour InfluxDB.
	3.	Compilation – Compilez le programme Go. Vous pouvez utiliser la commande go build pour obtenir un exécutable, ou go run main.go pour exécuter directement le script. Par exemple :

go build -o dashboard-generator .

Cette commande produit un binaire nommé dashboard-generator.

	4.	Lancement du serveur web – Exécutez le binaire ou lancez le script. Par exemple :

./dashboard-generator

Une fois lancé, le programme démarre un petit serveur web local. Par défaut, il écoute sur le port 8080. Dans la console, vous devriez voir un message confirmant le démarrage, par exemple : Serveur démarré sur : http://localhost:8080.

	5.	Accès au formulaire web – Ouvrez un navigateur web et rendez-vous à l’adresse du serveur (par exemple http://localhost:8080). Vous devriez voir apparaître une page web avec un formulaire de configuration.
	6.	Utilisation du formulaire – Sur la page, indiquez le Nombre d’appareils (nombre de capteurs) que vous souhaitez superviser, puis validez en cliquant sur le bouton Générer. (Voir la capture d’écran ci-dessous pour un exemple du formulaire.)
	7.	Génération du dashboard – À la soumission du formulaire, le script va automatiquement générer la configuration du (ou des) dashboard(s) Grafana et appeler l’API Grafana pour créer le tableau de bord correspondant. Si tout se passe bien, une page de confirmation s’affiche indiquant que le dashboard a été créé avec succès et mentionnant le nombre d’appareils configurés.
	8.	Vérification dans Grafana – Connectez-vous à votre interface Grafana. Un nouveau tableau de bord devrait être disponible (ou mis à jour) dans un dossier intitulé “Dashboards Dynamiques”. Le dashboard généré porte un titre par défaut, par exemple “Dashboard IoT Dynamique”, et contient un nombre de graphiques (panels) correspondant au nombre de capteurs saisi. Vous pouvez alors ouvrir ce tableau de bord pour visualiser les données de vos capteurs en temps réel.

Remarque : La première exécution crée le tableau de bord dans Grafana. Les exécutions suivantes mettent à jour ce tableau de bord (plutôt que d’en créer un nouveau à chaque fois), car le script utilise la même title de dashboard. Ainsi, changer le nombre de capteurs via le formulaire ajoutera ou retirera des panels dans le dashboard existant. Le rafraîchissement automatique du dashboard est configuré (5 secondes par défaut) afin d’actualiser régulièrement les mesures depuis InfluxDB.

Structure du Code (main.go)

Le cœur du projet réside dans le fichier main.go écrit en Go. Voici les principaux éléments de sa structure et leur rôle :
	•	Serveur Web & Routes HTTP – La fonction main() initialise un petit serveur web HTTP (via http.ListenAndServe sur le port 8080) et définit deux routes principales:
	•	La route “/” (racine) est gérée par la fonction homeHandler. Cette route sert la page HTML du formulaire permettant de configurer le nombre de capteurs. Le HTML est construit dans le code (utilisation de html/template) et comprend un formulaire avec un champ numérique pour le nombre d’appareils et un bouton de soumission.
	•	La route “/generate” est gérée par la fonction generateHandler. C’est cette route qui reçoit la requête POST lorsque l’utilisateur soumet le formulaire. Elle va alors traiter la génération du dashboard.
	•	Handler homeHandler – Cette fonction construit dynamiquement la page HTML du formulaire de configuration. Le formulaire contient un champ numDevices (de type nombre) pour saisir le nombre de capteurs souhaités. Une fois rendu, ce formulaire est présenté à l’utilisateur pour input.
	•	Handler generateHandler – Cette fonction est appelée lors du POST du formulaire. Son rôle est :
	1.	Lecture de l’entrée utilisateur – Récupérer la valeur numDevices envoyée par le formulaire (en tant que string) et la convertir en entier. Si la méthode HTTP n’est pas POST ou si la valeur est invalide (ex : inférieure à 1), la fonction renvoie une erreur ou redirige l’utilisateur vers la page du formulaire.
	2.	Génération du YAML de dashboard – Appeler la fonction generateYAML(numDevices) pour construire la configuration du dashboard Grafana au format YAML, en fonction du nombre de capteurs. Cette configuration inclura un panel par capteur.
	3.	Déploiement du dashboard sur Grafana – Appeler la fonction deployDashboard(yamlContent) qui utilise l’API de Grafana (via la librairie Go Grabana) pour créer ou mettre à jour le dashboard à partir du YAML généré. Si une erreur survient à ce stade (par exemple, problème de connexion à Grafana ou erreur dans le YAML), un message d’erreur est renvoyé à l’utilisateur.
	4.	Réponse à l’utilisateur – Si le déploiement réussit, la fonction renvoie une page HTML simple confirmant la réussite de la génération du dashboard et indiquant le nombre d’appareils configurés. Un lien de retour vers la page d’accueil est également proposé.
	•	Fonction generateYAML(numDevices int) string – Cette fonction génère une chaîne de caractères contenant la définition du dashboard en format YAML (conforme au schéma attendu par Grafana/Grabana). Les principales parties de ce YAML sont :
	•	Un titre de dashboard (par ex. "Dashboard IoT Dynamique") et quelques paramètres globaux (le dashboard est éditable, un tag auto-generated est ajouté, et un auto-refresh de 5s est défini).
	•	Une section de ligne/row intitulée “Mesures” qui regroupe les panneaux.
	•	À l’intérieur, la fonction boucle de 1 jusqu’au nombre de devices demandé (numDevices) et pour chaque capteur, ajoute un panel de type Timeseries. Chaque panel est configuré avec :
	•	Un titre correspondant au nom du device (par ex. "device_1", "device_2", etc.).
	•	Une datasource pointant vers InfluxDB (il est présupposé que Grafana a une datasource nommée “InfluxDB” connectée au bucket Influx approprié).
	•	Une requête Flux spécifique à InfluxDB pour ce capteur. La requête filtre les données par nom de mesure (_measurement) égal au nom du device et par champ (_field) correspondant aux métriques d’intérêt (par exemple consigne, setpoint ou temperature). Ensuite, la requête effectue une agrégation (mean) sur une fenêtre de temps adaptée à l’intervalle d’affichage, et renvoie la moyenne.
	•	Une hauteur de panel prédéfinie (par exemple 400px) pour un affichage harmonieux.
	•	La fonction construit le YAML via un strings.Builder. Pour diagnostic, elle écrit également le YAML généré dans un fichier local debug_dashboard.yaml (pratique pour vérifier la configuration générée). Enfin, elle retourne la chaîne YAML complète.
	•	Fonction deployDashboard(yamlContent string) error – Cette fonction utilise la librairie Grabana pour intéragir avec Grafana :
	•	Elle crée un client Grafana (grabana.NewClient) avec l’URL de Grafana et le token API (fournis en paramètres dans le code).
	•	Elle convertit le YAML du dashboard en objet interne (decoder.UnmarshalYAML) compréhensible par Grafana/Grabana.
	•	Elle cherche ou crée un dossier Grafana nommé “Dashboards Dynamiques” via l’API (cela permet de regrouper les dashboards générés automatiquement dans un même dossier).
	•	Puis elle envoie le dashboard à Grafana avec client.UpsertDashboard. La méthode Upsert crée le dashboard s’il n’existe pas, ou le met à jour s’il existe déjà (basé sur le titre) dans le dossier spécifié.
	•	Si tout se passe bien, Grafana enregistre le nouveau tableau de bord (ou les modifications du tableau existant). En cas d’échec, l’erreur est retransmise pour être gérée par le handler (affichage d’une erreur à l’utilisateur).

En résumé, le code est organisé pour offrir une interface web simple à l’utilisateur, convertir la demande de celui-ci en configuration de dashboard, puis appeler Grafana pour matérialiser ce dashboard sans intervention manuelle.

Comment ça marche (du formulaire à la création du dashboard)

Ce chapitre décrit le flot d’exécution complet, du moment où l’utilisateur utilise le formulaire web jusqu’à la visualisation du tableau de bord Grafana généré :
	1.	Saisie via le formulaire web – L’utilisateur accède à la page de configuration (route /) et voit un formulaire intitulé “Configuration du Dashboard”. Sur ce formulaire, il indique le nombre de capteurs à superviser en renseignant le champ “Nombre d’appareils”. Par exemple, s’il souhaite suivre 3 capteurs, il entre 3 dans le champ numérique. Une fois le nombre saisi, il clique sur le bouton Générer pour valider.
(Illustration : voir ci-dessous une capture d’écran du formulaire web de configuration.)
	2.	Envoi de la requête et traitement – En cliquant sur “Générer”, le navigateur envoie une requête HTTP POST à l’URL /generate du serveur Go, avec le champ numDevices=3 (dans notre exemple) encodé dans le corps de la requête. Le handler generateHandler reçoit cette requête sur le serveur :
	•	Il vérifie que la méthode est bien POST et extrait la valeur numDevices.
	•	Il convertit cette valeur en entier (dans l’exemple, 3) et vérifie qu’elle est valide (>= 1).
	•	Puis il appelle la fonction de génération de configuration YAML en passant ce nombre.
	3.	Génération de la configuration du dashboard – La fonction generateYAML(3) construit le contenu YAML d’un tableau de bord Grafana adapté pour 3 capteurs. Concrètement, le YAML contiendra 3 panneaux timeseries. Chacun aura pour titre device_1, device_2 et device_3 respectivement, et inclura la requête vers InfluxDB filtrant les données du capteur correspondant. À ce stade, un fichier debug_dashboard.yaml peut être créé en local avec cette configuration (utile pour débugger ou visualiser le résultat brut). Le script détient maintenant en mémoire la définition du dashboard à créer.
	4.	Appel à l’API Grafana – Le handler continue en appelant deployDashboard(yamlContent) avec le YAML généré. Cette fonction établit une connexion à Grafana (via l’API HTTP) en utilisant le client Grabana et le token API configuré. Elle s’assure que le dossier “Dashboards Dynamiques” existe, puis utilise l’API pour créer le dashboard. Grafana reçoit la requête de création (ou mise à jour) et enregistre le tableau de bord. Si Grafana retourne une erreur (par exemple, si le token est invalide ou la datasource InfluxDB n’est pas trouvée), le script renvoie un message d’erreur HTTP au navigateur, indiquant à l’utilisateur qu’il y a eu un problème. Sinon, la création se déroule correctement.
	5.	Confirmation côté utilisateur – Si la génération est un succès, le serveur renvoie une page HTML de confirmation. L’utilisateur voit apparaître un message du type “Dashboard généré avec succès ! 3 appareils configurés.” (selon le nombre entré), avec un lien pour retourner au formulaire initial si besoin. Cela confirme que Grafana a bien pris en compte la création du dashboard.
	6.	Visualisation du dashboard Grafana – L’utilisateur peut maintenant se rendre sur Grafana (via l’interface web de Grafana) pour voir le résultat. Dans le dossier Dashboards Dynamiques, il trouvera un nouveau tableau de bord nommé “Dashboard IoT Dynamique”. En l’ouvrant, il découvre 3 graphiques alignés (si 3 capteurs ont été configurés dans l’exemple), chacun affichant les mesures du capteur correspondant. Chaque graphique interroge InfluxDB en temps réel pour afficher, par exemple, la température, la consigne ou le setpoint du capteur, et se mettra à jour automatiquement toutes les 5 secondes (grâce à l’auto-refresh configuré).
(Illustration : voir ci-dessous un exemple de tableau de bord Grafana généré, contenant plusieurs panels de données de capteurs.)
	7.	Flux des données – En arrière-plan, Node-RED continue à acheminer les données de capteurs vers InfluxDB. Grâce à la convention de nommage des mesures (chaque capteur ayant un identifiant unique utilisé comme nom de “measurement” dans InfluxDB, par ex. device_1), les requêtes configurées dans chaque panel Grafana récupèrent les bonnes données. Ainsi, le dashboard affiche en direct les valeurs envoyées par les capteurs via Node-RED. Si de nouveaux capteurs doivent être ajoutés plus tard, il suffira de retourner au formulaire, d’indiquer le nouveau nombre total et de générer à nouveau le dashboard : Grafana mettra alors à jour le tableau de bord en conséquence, sans avoir besoin de recréer manuellement des visualisations.

Aperçu visuel du système

Pour mieux illustrer le fonctionnement, voici quelques captures d’écran du formulaire de configuration et du dashboard Grafana généré :

Formulaire Web de configuration

Figure 1 : Formulaire web où l’utilisateur spécifie le nombre de capteurs à monitorer. Après validation, le système génère automatiquement le dashboard Grafana correspondant.

Dashboard Grafana généré automatiquement

Figure 2 : Exemple de tableau de bord Grafana généré par le script. Ici, trois capteurs (“device_1” à “device_3”) sont supervisés, chacun ayant son graphique de mesures (température, consigne, etc.). Le tableau de bord est actualisé régulièrement via l’auto-refresh.

(Remarque : Les captures ci-dessus sont indicatives. Le style et l’apparence peuvent varier selon la version de Grafana et la personnalisation effectuée.)

Exemple d’utilisation et configuration

Pour illustrer concrètement l’utilisation, prenons un exemple simple :

Supposons que l’on veuille superviser 3 capteurs IoT. L’utilisateur accède au formulaire web et entre 3 comme nombre d’appareils, puis clique sur Générer. Le script crée alors (ou met à jour) le tableau de bord Grafana avec 3 panels. Chaque panel est connecté aux données InfluxDB du capteur correspondant (par exemple, un capteur de température device_1, un capteur d’humidité device_2, etc., selon la façon dont Node-RED nomme et stocke les mesures). Quelques instants plus tard, en ouvrant Grafana, l’utilisateur voit le tableau de bord contenant ces 3 graphiques mis à jour en temps réel.

D’un point de vue technique, cette action équivaut à appeler directement l’API HTTP du script. Par exemple, on pourrait obtenir le même résultat en envoyant manuellement une requête POST avec un outil comme cURL :

curl -X POST http://localhost:8080/generate -d "numDevices=3"

La commande ci-dessus enverrait la requête de génération pour 3 devices sans passer par le formulaire web (ce qui est exactement ce que fait le formulaire en arrière-plan). Le script répondrait avec le HTML de confirmation, et Grafana aurait, suite à cet appel, le dashboard mis à jour.

En termes de configuration, assurez-vous que Grafana dispose bien d’une datasource InfluxDB nommée InfluxDB pointant vers votre instance InfluxDB et le bucket adéquat (par exemple iot-platform). De plus, adaptez dans le code (ou via des variables d’environnement si vous modifiez le script) l’URL de Grafana et le token API. Par défaut, ces valeurs sont fixées en dur dans main.go à des fins de démonstration, mais pour un déploiement réel il est recommandé de les externaliser pour plus de sécurité et de flexibilité.

Conclusion

Ce projet fournit une solution automatisée pour la création de tableaux de bord Grafana dans un contexte IoT. Grâce à un simple formulaire web, il élimine la nécessité de configurer manuellement chaque nouveau capteur dans Grafana. L’intégration transparente avec Node-RED et InfluxDB permet de rapidement passer du déploiement de capteurs à leur visualisation. Ce README a présenté le fonctionnement général, la mise en place et l’utilisation du script, la structure du code Go, ainsi qu’un exemple concret. Avec cette base en place, il est possible d’étendre ou d’adapter le système – par exemple, personnaliser les panels Grafana, ajouter d’autres types de métriques, ou intégrer une interface utilisateur plus élaborée – tout en conservant le principe d’une génération automatique des dashboards pour simplifier la supervision IoT.



