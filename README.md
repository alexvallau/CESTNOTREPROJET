
# Générateur Automatique de Dashboards Grafana pour Capteurs IoT

## Table des Matières
- [Générateur Automatique de Dashboards Grafana pour Capteurs IoT](#générateur-automatique-de-dashboards-grafana-pour-capteurs-iot)
  - [Table des Matières](#table-des-matières)
  - [Introduction](#introduction)
  - [Prérequis Techniques](#prérequis-techniques)
  - [Déploiement et Exécution du Script](#déploiement-et-exécution-du-script)
    - [Étapes](#étapes)
  - [Fonctionnement du Code](#fonctionnement-du-code)
    - [Structure Principale](#structure-principale)
  - [Exemple d’Utilisation](#exemple-dutilisation)
    - [Étapes](#étapes-1)
    - [Requête cURL](#requête-curl)
  - [Conclusion](#conclusion)


## Introduction

Ce projet est un script écrit en Go qui permet de générer automatiquement des tableaux de bord Grafana pour le monitoring de capteurs IoT. Il s’inscrit dans un trio technologique comprenant **Node-RED**, **InfluxDB** et **Grafana**. 

L’idée générale est la suivante : 
- **Node-RED** collecte les données de capteurs et les stocke dans **InfluxDB**.
- **Grafana** visualise ces données via des dashboards.
- Ce script Go automatise la création des dashboards pour chaque nouveau capteur.

L’utilisateur indique le nombre de capteurs (devices) à superviser via un formulaire web, et le système génère dynamiquement le tableau de bord Grafana contenant les graphiques correspondants. Ce processus facilite le déploiement rapide de la supervision pour un nombre variable de capteurs, sans intervention manuelle dans Grafana.

---

## Prérequis Techniques

Avant de déployer ce projet, assurez-vous de disposer de l’environnement suivant :

- **Go** (version 1.XX ou ultérieure) – pour compiler et exécuter le script Go.
- **Grafana** (idéalement version 8 ou 9) – pour l’affichage des dashboards. Grafana doit être accessible via son API HTTP et disposer d’une datasource InfluxDB configurée.
- **InfluxDB** (v2 ou supérieur, avec API Flux) – base de données time-series où sont stockées les mesures des capteurs. Le bucket de données IoT doit être configuré (par exemple `iot-platform`) et relié à Grafana via une datasource.
- **Node-RED** – pour la collecte des données capteur et l’envoi vers InfluxDB. Node-RED n’interagit pas directement avec le script Go, mais il alimente InfluxDB en données (mesures) que les dashboards Grafana afficheront.
- **Token API Grafana** – une clé API Grafana avec les droits nécessaires (édition de dashboards) est requise pour que le script puisse créer/modifier des dashboards via l’API. Ce token doit être configuré dans le script (ou dans l’environnement) avant exécution.
- **Accès réseau** – le script doit pouvoir communiquer avec l’instance Grafana (URL et port de Grafana accessibles, par ex. `http://localhost:3000` ou une IP interne) et Grafana doit pouvoir accéder à InfluxDB.

---

## Déploiement et Exécution du Script

### Étapes

1. **Récupération du code**  
    Clonez le dépôt du projet sur le serveur web ou la machine où il sera exécuté. Assurez-vous d’avoir le fichier `main.go` du projet.

2. **Configuration**  
    - Mettez à jour les paramètres si nécessaire : dans le code `main.go`, vérifiez l’URL de Grafana et le token API Grafana.  
    - Par défaut, le code utilise l’URL `http://192.168.141.122:3000` et un token codé en dur – adaptez ces valeurs à votre environnement Grafana (par exemple `http://localhost:3000` et votre propre jeton d’API).  
    - Assurez-vous également que le nom de la datasource Grafana (ici `InfluxDB`) correspond bien à celle configurée dans Grafana pour InfluxDB.

3. **Compilation**  
    Compilez le programme Go. Vous pouvez utiliser la commande suivante :

    ```bash
    go build -o dashboard-generator .
    ```

    Cette commande produit un binaire nommé `dashboard-generator`.

4. **Lancement du serveur web**  
    Exécutez le binaire ou lancez le script :

    ```bash
    ./dashboard-generator
    ```

    Une fois lancé, le programme démarre un petit serveur web local. Par défaut, il écoute sur le port `8080`. Dans la console, vous devriez voir un message confirmant le démarrage, par exemple :  
    `Serveur démarré sur : http://localhost:8080`.

5. **Accès au formulaire web**  
    Ouvrez un navigateur web et rendez-vous à l’adresse du serveur (par exemple `http://localhost:8080`). Vous devriez voir apparaître une page web avec un formulaire de configuration.

6. **Utilisation du formulaire**  
    Sur la page, indiquez le **Nombre d’appareils** (nombre de capteurs) que vous souhaitez superviser, puis validez en cliquant sur le bouton **Générer**.

7. **Génération du dashboard**  
    À la soumission du formulaire, le script va automatiquement générer la configuration du (ou des) dashboard(s) Grafana et appeler l’API Grafana pour créer le tableau de bord correspondant. Si tout se passe bien, une page de confirmation s’affiche indiquant que le dashboard a été créé avec succès et mentionnant le nombre d’appareils configurés.

8. **Vérification dans Grafana**  
    Connectez-vous à votre interface Grafana. Un nouveau tableau de bord devrait être disponible (ou mis à jour) dans un dossier intitulé **Dashboards Dynamiques**. Le dashboard généré porte un titre par défaut, par exemple **Dashboard IoT Dynamique**, et contient un nombre de graphiques (panels) correspondant au nombre de capteurs saisi.

---

## Fonctionnement du Code

### Structure Principale

- **Serveur Web & Routes HTTP**  
  - La route `/` sert la page HTML du formulaire.  
  - La route `/generate` reçoit la requête POST et génère le dashboard.

- **Génération du YAML**  
  La fonction `generateYAML(numDevices)` construit la configuration du dashboard Grafana au format YAML, incluant un panel par capteur.

- **Déploiement via l’API Grafana**  
  La fonction `deployDashboard(yamlContent)` utilise la librairie Grabana pour créer ou mettre à jour le dashboard dans Grafana.

---

## Exemple d’Utilisation

### Étapes

1. Accédez au formulaire web et entrez le nombre de capteurs (par exemple, `3`).
2. Cliquez sur **Générer**.
3. Vérifiez dans Grafana : un tableau de bord nommé **Dashboard IoT Dynamique** contenant 3 graphiques sera disponible.

### Requête cURL

Vous pouvez également générer un dashboard directement via une requête cURL :

```bash
curl -X POST http://localhost:8080/generate -d "numDevices=3"
```

---

## Conclusion

Ce projet fournit une solution automatisée pour la création de tableaux de bord Grafana dans un contexte IoT. Grâce à un simple formulaire web, il élimine la nécessité de configurer manuellement chaque nouveau capteur dans Grafana. L’intégration transparente avec Node-RED et InfluxDB permet de rapidement passer du déploiement de capteurs à leur visualisation.

