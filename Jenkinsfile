pipeline {
    agent any

    environment {
        SONAR_HOST_URL = credentials('e53ab44a-35fb-41ae-80ab-dd6836a9463c') // SONAR_HOST
        SONAR_TOKEN = credentials('9c1c3109-58e4-4890-b88f-2615d2221245') // SONAR_TOKEN
        HARBOR_USERNAME = credentials('db2c5c66-275f-440f-a0d5-73dce0f7355e') // HARBOR_USERNAME
        HARBOR_PASSWORD = credentials('a6c7d1c9-3c1b-4bdb-a0c5-4ca28f51c1f5') // HARBOR_PASSWORD
    }

    stages {
        stage('Test Credentials') {
            steps {
                script {
                    echo "SONAR_HOST_URL: ${env.SONAR_HOST_URL}"
                    echo "SONAR_TOKEN: ${env.SONAR_TOKEN}"
                    echo "HARBOR_USERNAME: ${env.HARBOR_USERNAME}"
                    echo "HARBOR_PASSWORD: ${env.HARBOR_PASSWORD}"
                }
            }
        }
        
        stage('SonarQube Analysis') {
            steps {
                script {
                    withDockerContainer(args: '-u root:root', image: 'sonarsource/sonar-scanner-cli') { // Requires Docker Pipeline plugin
                        sh 'sonar-scanner -Dsonar.projectKey=MSPR-PayeTonKawa_auth_7d40a8c4-4ff5-4034-acaf-0226d044b7c0 -Dsonar.sources=. -Dsonar.host.url=$SONAR_HOST_URL -Dsonar.login=$SONAR_TOKEN'
                    }
                }
            }
        }

        stage('Go Test') {
            steps {
                script {
                    withDockerContainer(args: '-u root:root', image: 'golang') { // Requires Docker Pipeline plugin
                        sh 'go test ./... -v'
                    }
                }
            }
        }

        stage('Docker Build and Push') {
            steps {
                script {
                    withDockerContainer(args: '-u root:root', image: 'plugins/docker') { // Requires Docker Pipeline plugin
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
        pollSCM('* * * * *') // Poll SCM for changes every minute // Requires Git plugin
    }
}
