pipeline {
    agent {
        label 'master'
    }

    stages {
        stage('Build') {
            steps {
                sh "docker pull golang"
                sh "docker run -t --rm -v ${WORKSPACE}:/go/src/cli -w /go/src/cli golang make || exit 0"
                s3Upload(file:'./bin', bucket:'smartling-connectors-releases', path:'cli/')
            }
        }
    }

    post {
        unstable {
            slackSend (
                    channel: "#emergency-connectors",
                    color: 'bad',
                    message: "Tests failed: <${env.RUN_DISPLAY_URL}|${env.JOB_NAME} #${env.BUILD_NUMBER}>"
            )
        }

        failure {
            slackSend (
                    channel: "#emergency-connectors",
                    color: 'bad',
                    message: "Build of <${env.RUN_DISPLAY_URL}|${env.JOB_NAME} #${env.BUILD_NUMBER}> is failed!"
            )
        }
    }
}
