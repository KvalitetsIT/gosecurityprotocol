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

	}
	post {
		always {

			dir('testenv') {
				sh 'docker-compose rm -f'
			}
		}
	}
}
