pipeline {
    agent {
        label 'master'
    }

    stages {
        stage('Build') {
            steps {
                sh "docker pull golang"
                sh "docker run -t --rm -v ${WORKSPACE}:/go/src/cli -w /go/src/cli golang make || exit 0"
                sh "aws-profile ${AWS_PROFILE} aws s3 cp ./bin s3://smartling-connectors-releases/cli/ --recursive"
            }
        }
    }

    post {
        unstable {
            slackSend (
                    channel: "${SLACK_CHANNEL}",
                    color: 'bad',
                    message: "Tests failed: <${env.RUN_DISPLAY_URL}|${env.JOB_NAME} #${env.BUILD_NUMBER}>"
            )
        }

        failure {
            slackSend (
                    channel: "${SLACK_CHANNEL}",
                    color: 'bad',
                    message: "Build of <${env.RUN_DISPLAY_URL}|${env.JOB_NAME} #${env.BUILD_NUMBER}> is failed!"
            )
        }
    }
}
