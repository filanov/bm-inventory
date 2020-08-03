pipeline {
  agent any
  stages {
    stage('clear deployment') {
      steps {
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
      	echo "love"
	}
    }
stage('Deploy to prod') {
  when {
    branch 'master'
  }
  agent any
  steps {
	sh '''docker login quay.io -u oscohen -p nata2411'''
	sh '''docker tag quay.io/ocpmetal/bm-inventory:test quay.io/ocpmetal/bm-inventory-push-test'''
	sh '''docker push quay.io/ocpmetal/bm-inventory'''  
}
}

  }
}
