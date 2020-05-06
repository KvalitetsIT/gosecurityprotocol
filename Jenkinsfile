podTemplate(
        containers: [containerTemplate(image: 'kvalitetsit/docker-compose:dev', name: 'docker', command: 'cat', ttyEnabled: true)],
        volumes: [hostPathVolume(hostPath: '/var/run/docker.sock', mountPath: '/var/run/docker.sock')],
) {
    node(POD_LABEL) {

		stage('Clone repository') {
				checkout scm
		}

		stage('Startup the testenvironment used by the integration tests') {
            container('docker') {
                dir('testenv') {
                    sh 'docker-compose up -d'
                }
            }
		}
		stage('Build Docker image') {
            container('docker') {
			  docker.build("kvalitetsit/gosecurityprotocol", "--network testenv_gosecurityprotocol -f Dockerfile .")
			}
		}
	}
	post {
		always {
            container('docker') {
                dir('testenv') {
                    sh 'docker-compose stop'
                    sh 'docker-compose rm -f'
                }
            }
		}
	}
}
