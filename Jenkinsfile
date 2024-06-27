pipeline {
    agent {
        kubernetes {
            label 'go-test-pod'
            yaml """
            apiVersion: v1
            kind: Pod
            spec:
              containers:
              - name: go
                image: golang:1.22
                command:
                - cat
                tty: true
              - name: docker
                image: docker:20.10.7
                command:
                - cat
                tty: true
                volumeMounts:
                - name: docker-sock
                  mountPath: /var/run/docker.sock
              - name: sonar-scanner
                image: sonarsource/sonar-scanner-cli
                command:
                - cat
                tty: true
              - name: kubectl
                image: bitnami/kubectl:latest
                command:
                - cat
                tty: true
                volumeMounts:
                - name: kubeconfig
                  mountPath: /root/.kube
              volumes:
              - name: docker-sock
                hostPath:
                  path: /var/run/docker.sock
              - name: kubeconfig
                secret:
                  secretName: kubeconfig
            """
        }
    }

    environment {
        SONAR_HOST_URL = credentials('e53ab44a-35fb-41ae-80ab-dd6836a9463c') // SONAR_HOST
        SONAR_TOKEN = credentials('9c1c3109-58e4-4890-b88f-2615d2221245') // SONAR_TOKEN
        HARBOR_USERNAME = credentials('db2c5c66-275f-440f-a0d5-73dce0f7355e') // HARBOR_USERNAME
        HARBOR_PASSWORD = credentials('a6c7d1c9-3c1b-4bdb-a0c5-4ca28f51c1f5') // HARBOR_PASSWORD
        KUBECONFIG = credentials('ec2c0a90-1f2e-461e-8851-5add47e2c7b2') // Kubeconfig secret
    }

    stages {
        stage('Test') {
            steps {
                container('go') {
                    sh 'go test ./... -v'
                }
            }
        }

        stage('Build Docker Image') {
            steps {
                container('docker') {
                    script {
                        sh 'docker build -t registry.germainleignel.com/paye-ton-kawa/auth:latest .'
                    }
                }
            }
        }

        stage('Push Docker Image') {
            steps {
                container('docker') {
                    script {
                        sh 'echo $HARBOR_PASSWORD | docker login registry.germainleignel.com --username $HARBOR_USERNAME --password-stdin'
                        sh 'docker push registry.germainleignel.com/paye-ton-kawa/auth:latest'
                    }
                }
            }
        }

        stage('SonarQube Analysis') {
            steps {
                container('sonar-scanner') {
                    sh 'sonar-scanner -Dsonar.projectKey=MSPR-PayeTonKawa_auth_7d40a8c4-4ff5-4034-acaf-0226d044b7c0 -Dsonar.sources=. -Dsonar.host.url=$SONAR_HOST_URL -Dsonar.login=$SONAR_TOKEN'
                }
            }
        }

        stage('Deploy to Kubernetes') {
            steps {
                container('kubectl') {
                    script {
                        // Save the KUBECONFIG file to the correct location
                        sh 'mkdir -p /root/.kube && echo "$KUBECONFIG" > /root/.kube/config'
                        // Apply the YAML files to the cluster
                        sh 'kubectl apply -f path/to/your/yaml/files --kubeconfig=/root/.kube/config'
                    }
                }
            }
        }
    }

    post {
        always {
            archiveArtifacts artifacts: '**/test-results/*.xml', allowEmptyArchive: true
        }
        success {
            echo 'Tests ran successfully, SonarQube analysis completed, and image was built, pushed, and deployed.'
        }
        failure {
            echo 'Tests, SonarQube analysis, Docker build/push, or deployment failed.'
        }
    }
}
