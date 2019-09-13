node {
	def app
	def scmInfo

	stage('Clone repository') {
		scmInfo = checkout scm
	}

	stage('Startup the testenvironment used by the integration tests') {
		dir('testenv') {
			sh 'docker-compose up -d'
		}
	}

	stage('Build Docker image') {
            app = docker.build("kvalitetsit/loginproxy-siemens-documentconsumer:${scmInfo.GIT_COMMIT}", "--network testenv_gosecurityprotocol -f Dockerfile .")
	}

	always {

		stage('Stop and remove the testenvironment used by the integration tests') {
		
			dir('testenv') {
				sh 'docker-compose rm -s -f'
			}
		}
	}
}
