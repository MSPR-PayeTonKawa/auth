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
            steps {
                script {
                    docker.image('sonarsource/sonar-scanner-cli').inside('-u root:root') {
                        sh 'sonar-scanner -Dsonar.projectKey=MSPR-PayeTonKawa_auth_7d40a8c4-4ff5-4034-acaf-0226d044b7c0 -Dsonar.sources=. -Dsonar.host.url=$SONAR_HOST_URL -Dsonar.login=$SONAR_TOKEN'
                    }
                }
            }
        }

        stage('Go Test') {
            steps {
                script {
                    docker.image('golang').inside('-u root:root') {
                        sh 'go test ./... -v'
                    }
                }
            }
        }

        stage('Docker Build and Push') {
            steps {
                script {
                    docker.image('plugins/docker').inside('-u root:root') {
                        sh '''
                            docker login -u $HARBOR_USERNAME -p $HARBOR_PASSWORD registry.germainleignel.com
                            docker build -t registry.germainleignel.com/paye-ton-kawa/auth .
                            docker push registry.germainleignel.com/paye-ton-kawa/auth
                        '''
                    }
                }
            }
        }
    }

    triggers {
        pollSCM('* * * * *') // Poll SCM for changes every minute
    }
}
