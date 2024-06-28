pipeline {
    agent {
        kubernetes {
            label 'go-docker-agent'
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
                - name: kubeconfig
                  mountPath: /root/.kube/config
                  subPath: config
              serviceAccountName: jenkins-agent-sa
              volumes:
              - name: docker-sock
                hostPath:
                  path: /var/run/docker.sock
              - name: kubeconfig
                configMap:
                  name: kubeconfig
            """
        }
    }

    environment {
        SONAR_HOST_URL = credentials('e53ab44a-35fb-41ae-80ab-dd6836a9463c') // SONAR_HOST
        SONAR_TOKEN = credentials('9c1c3109-58e4-4890-b88f-2615d2221245') // SONAR_TOKEN
        HARBOR_USERNAME = credentials('db2c5c66-275f-440f-a0d5-73dce0f7355e') // HARBOR_USERNAME
        HARBOR_PASSWORD = credentials('a6c7d1c9-3c1b-4bdb-a0c5-4ca28f51c1f5') // HARBOR_PASSWORD
    }

    stages {
        stage('Setup RBAC') {
            steps {
                script {
                    sh '''
                    kubectl apply -f - <<EOF
                    apiVersion: v1
                    kind: ServiceAccount
                    metadata:
                      name: jenkins-agent-sa
                      namespace: default
                    ---
                    apiVersion: rbac.authorization.k8s.io/v1
                    kind: ClusterRole
                    metadata:
                      name: node-listing-role
                    rules:
                    - apiGroups: [""]
                      resources: ["nodes"]
                      verbs: ["get", "list", "watch"]
                    ---
                    apiVersion: rbac.authorization.k8s.io/v1
                    kind: ClusterRoleBinding
                    metadata:
                      name: jenkins-agent-sa-cluster-admin
                    subjects:
                    - kind: ServiceAccount
                      name: jenkins-agent-sa
                      namespace: default
                    roleRef:
                      kind: ClusterRole
                      name: cluster-admin
                      apiGroup: rbac.authorization.k8s.io
                    EOF
                    '''
                }
            }
        }

        stage('Test') {
            steps {
                container('go') {
                    // Running go test with verbosity
                    sh 'go test ./... -v'
                }
            }
        }

        stage('Build Docker Image') {
            steps {
                container('docker') {
                    script {
                        def imageName = "registry.germainleignel.com/paye-ton-kawa/auth:${env.BUILD_NUMBER}"
                        sh "docker build -t ${imageName} ."
                    }
                }
            }
        }

        stage('Push Docker Image') {
            steps {
                container('docker') {
                    script {
                        def imageName = "registry.germainleignel.com/paye-ton-kawa/auth:${env.BUILD_NUMBER}"
                        // Corrected to prevent Groovy string interpolation
                        sh 'echo $HARBOR_PASSWORD | docker login registry.germainleignel.com --username $HARBOR_USERNAME --password-stdin'
                        sh "docker push ${imageName}"
                    }
                }
            }
        }

        stage('Apply Kubernetes Manifests') {
            steps {
                container('docker') {
                    script {
                        sh 'kubectl apply -f your-manifest-file.yaml'
                    }
                }
            }
        }
    }

    post {
        always {
            // Archive test results, logs, or any other artifacts if needed
            archiveArtifacts artifacts: '**/test-results/*.xml', allowEmptyArchive: true
        }
        success {
            echo 'Tests ran successfully, image was built and pushed, and manifests were applied.'
        }
        failure {
            echo 'Tests, Docker build/push, or manifest application failed.'
        }
    }
}
