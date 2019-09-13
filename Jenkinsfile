pipeline {
	agent any

	stages {

		stage('Clone repository') {
			checkout scm
		}

		stage('Startup the testenvironment used by the integration tests') {
			dir('testenv') {
				sh 'docker-compose up -d'
			}
		}

		stage('Build Docker image') {
			docker.build("kvalitetsit/loginproxy-siemens-documentconsumer", "--network testenv_gosecurityprotocol -f Dockerfile .")
		}
	}
	post {
		always {

			stage('Stop and remove the testenvironment used by the integration tests') {
		
				dir('testenv') {
					sh 'docker-compose rm -s -f'
				}
			}
		}
	}
}
