pipeline {
  agent any
  stages {
    stage('clear deployment') {
      steps {
        sh 'export PATH=$PATH:/usr/local/go/bin;'
        sh 'make clear-deployment'
      }
    }

    stage('Deploy') {
      steps {
        sh '''export OBJEXP=quay.io/ocpmetal/s3-object-expirer:latest; make deploy-test
'''
        sleep 60
        sh '''# Dump pod statuses
kubectl  get pods -A'''
      }
    }

    stage('test') {
      steps {
        sh 'make subsystem-run'
      }
    }

  }
}