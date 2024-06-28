pipeline {
    agent {
        kubernetes {
            label 'go-docker-agent'
            defaultContainer 'jnlp'
            yaml '''
            apiVersion: v1
            kind: Pod
            spec:
              containers:
              - name: go
                image: golang:1.22
                command:
                - cat
                tty: true
                volumeMounts:
                - name: workspace-volume
                  mountPath: /home/jenkins/agent
              - name: docker
                image: docker:20.10.7
                command:
                - cat
                tty: true
                volumeMounts:
                - name: docker-sock
                  mountPath: /var/run/docker.sock
              - name: kubectl
                image: bitnami/kubectl:latest
                command:
                - cat
                tty: true
                volumeMounts:
                - name: kubeconfig
                  mountPath: /root/.kube/config
                  subPath: config
                - name: workspace-volume
                  mountPath: /home/jenkins/agent
              volumes:
              - name: docker-sock
                hostPath:
                  path: /var/run/docker.sock
              - name: workspace-volume
                emptyDir: {}
              - name: kubeconfig
                configMap:
                  name: kubeconfig
            '''
        }
    }
    stages {
        stage('Checkout SCM') {
            steps {
                checkout scm
            }
        }
        stage('Deploy App to Kubernetes') {
            steps {
                container('kubectl') {
                    withCredentials([file(credentialsId: 'ec2c0a90-1f2e-461e-8851-5add47e2c7b2', variable: 'KUBECONFIG')]) {
                        sh 'kubectl apply -f ./k8s/deployment.yaml'
                        sh 'kubectl apply -f ./k8s/middleware.yaml'
                        sh 'kubectl apply -f ./k8s/postgres-deployment.yaml'
                        sh 'kubectl apply -f ./k8s/postgres-service.yaml'
                        sh 'kubectl apply -f ./k8s/secrets.yaml'
                        sh 'kubectl apply -f ./k8s/service.yaml'
                    }
                }
            }
        }
    }
    post {
        always {
            archiveArtifacts artifacts: '**/*.yaml', allowEmptyArchive: true
            echo 'Tests, Docker build/push, or manifest application failed.'
        }
    }
}
