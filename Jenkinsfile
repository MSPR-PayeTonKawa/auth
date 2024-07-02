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
    parameters {
        string(name: 'PROJECT_PATH', defaultValue: '/home/gmn/apps/payetonkawa/auth', description: 'Path to the project directory')
    }
    stages {
        stage('Checkout SCM') {
            steps {
                checkout scm
            }
        }

        // stage('Build Docker Image') {
        //     steps {
        //         container('docker') {
        //             script {
        //                 def imageName = "registry.germainleignel.com/paye-ton-kawa/auth:${env.BUILD_NUMBER}"
        //                 sh "docker build -t ${imageName} ."
        //             }
        //         }
        //     }
        // }

        // stage('Push Docker Image') {
        //     steps {
        //         container('docker') {
        //             script {
        //                 def imageName = "registry.germainleignel.com/paye-ton-kawa/auth:${env.BUILD_NUMBER}"
        //                 withCredentials([usernamePassword(credentialsId: 'harbor-credentials', usernameVariable: 'HARBOR_USERNAME', passwordVariable: 'HARBOR_PASSWORD')]) {
        //                     sh "echo ${HARBOR_PASSWORD} | docker login registry.germainleignel.com -u ${HARBOR_USERNAME} --password-stdin"
        //                     sh "docker push ${imageName}"
        //                 }
        //             }
        //         }
        //     }
        // }
        stage('Deploy to K8s') {
            steps {
                sshagent(['6ff897ff-0cc3-4a47-86ca-a467266a6e4b']) {
                    sh """
                        ssh gmn@target-server "/home/gmn/scripts/deploy.sh ${params.PROJECT_PATH}"
                    """
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
