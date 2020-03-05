pipeline {
	agent any

	stages {

		stage('Clone repository') {
			steps {
				checkout scm
			}
		}

		stage('Startup the testenvironment used by the integration tests') {
			steps {
				dir('testenv') {
					sh 'docker-compose up -d'
				}
			}
		}
		stage('Build Docker image') {
			steps {
				script {
					docker.build("kvalitetsit/gosecurityprotocol", "--network testenv_gosecurityprotocol -f Dockerfile .")
				}
			}
		}

	}
	post {
		always {

			dir('testenv') {
				sh 'docker-compose stop'
				sh 'docker-compose rm -f'
			}
		}
	}
}
