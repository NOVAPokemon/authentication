pipeline {
    agent any
    stages {
            stage('Build executable') {
                steps {
                      container('golang') {
                        sh 'build --race -v -o executable .'
                    }
                }
            }

            stage('Docker Build') {
                steps {
                    container('docker') {
                        sh "docker build -t authentication:race ."
                    }
                }
            }
        }
}
