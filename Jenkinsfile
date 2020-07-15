pipeline {
    agent any
    stages {
        stage('Lint') {
            steps {
                sh 'skipper make lint'
            }
        }
        stage('Unit Test') {
            steps {
                sh 'skipper make unit-test'
            }
        }
    }
}