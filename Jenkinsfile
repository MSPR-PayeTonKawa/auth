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
              - name: busybox
                image: busybox
                command:
                - cat
                tty: true
            """
        }
    }

    stages {
        stage('Checkout') {
            steps {
                // Checkout the source code from the repository
                git branch: 'main', url: 'git@github.com:MSPR-PayeTonKawa/auth.git'
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
    }

    post {
        always {
            // Archive test results, logs, or any other artifacts if needed
            archiveArtifacts artifacts: '**/test-results/*.xml', allowEmptyArchive: true
        }
        success {
            echo 'Tests ran successfully.'
        }
        failure {
            echo 'Tests failed.'
        }
    }
}
