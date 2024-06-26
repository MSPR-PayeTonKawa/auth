pipeline {
    agent any

    environment {
        SONAR_HOST_URL = credentials('sonar_host')
        SONAR_TOKEN = credentials('sonar_token')
        HARBOR_USERNAME = credentials('harbor_username')
        HARBOR_PASSWORD = credentials('harbor_password')
    }

    stages {
        stage('SonarQube Analysis') {
            agent {
                docker {
                    image 'sonarsource/sonar-scanner-cli'
                    args '-u root:root'
                }
            }
            steps {
                sh 'sonar-scanner -Dsonar.projectKey=MSPR-PayeTonKawa_auth_7d40a8c4-4ff5-4034-acaf-0226d044b7c0 -Dsonar.sources=. -Dsonar.host.url=$SONAR_HOST_URL -Dsonar.login=$SONAR_TOKEN'
            }
        }

        stage('Go Test') {
            agent {
                docker {
                    image 'golang'
                    args '-u root:root'
                }
            }
            steps {
                sh 'go test ./... -v'
            }
        }

        stage('Docker Build and Push') {
            agent {
                docker {
                    image 'plugins/docker'
                    args '-u root:root'
                }
            }
            steps {
                sh '''
                    docker login -u $HARBOR_USERNAME -p $HARBOR_PASSWORD registry.germainleignel.com
                    docker build -t registry.germainleignel.com/paye-ton-kawa/auth .
                    docker push registry.germainleignel.com/paye-ton-kawa/auth
                '''
            }
        }
    }
    triggers {
        githubPush(branch: 'main')
    }
}
