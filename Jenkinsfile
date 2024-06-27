pipeline {
    agent any

    stages {
        stage('Go Test') {
            steps {
                script {
                    withDockerContainer(args: '-u root:root', image: 'golang') { // Requires Docker Pipeline plugin
                        sh 'go test ./... -v'
                    }
                }
            }
        }
    }

    triggers {
        pollSCM('* * * * *') // Poll SCM for changes every minute // Requires Git plugin
    }
}
