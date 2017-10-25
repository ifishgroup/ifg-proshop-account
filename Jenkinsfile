#!groovy

def version  = "0.1.${env.BUILD_NUMBER}"
def image    = "ifg-proshop-account-service"
def imageDb  = "ifg-proshop-account-service-db"

try {
    node('docker') {

        // notifySlack('STARTED')

        withEnv([
                "COMPOSE_FILE=common-services.yml"
        ]) {

            def dbName = "postgres-${convertBranchName(env.BRANCH_NAME)}-${version}"

            stage('checkout') {
                checkout scm
            }

            stage('unit/integration test') {
                sh "docker build -t $imageDb:$version db/"
                sh "docker run --rm --name ${dbName} -p 5432:5432 -d $imageDb:$version"
                sh "docker run --rm --link ${dbName} -v ${env.WORKSPACE}:/go/src/github.com/ifishgroup/ifg-proshop-account -w /go/src/github.com/ifishgroup/ifg-proshop-account golang:1.8.3 bash -c \"go get -d -v -t && go test --cover -v -tags=integration ./...\""
                sh "docker stop ${dbName}"
            }

            if (env.BRANCH_NAME =~ /(?i)^pr-/ || env.BRANCH_NAME == "master") {
                stage('docker build') {
                    sh "docker build -t $image:$version ."
                }

                try {
                    stage('deploy to staging') {
                        echo """
                            provision aws environment
                            deploy container
                            publishStagedInfo
                        """
                    }

                    stage('functional test') {
                        echo "run functional tests"
                    }

                    stage('load test') {
                        echo "production readiness tests"
                    }

                    stage('publish reports') {
                        echo "publish load testing reports"
                    }
                } catch (Exception e) {
                    throw e
                } finally {
                    stage('staging teardown') {
                        echo "teardown staged environment"
                        // notifyGithub("Staged build @ $ip was removed")
                        // slackSend(color: 'good', message: "Staged build @ $ip was removed")
                    }
                }
            }

            if (env.BRANCH_NAME == "master") {
                stage('publish') {
                    sh "docker tag $image:$version dlish27/$image:latest"
                    sh "docker tag $image:$version dlish27/$image:$version"
                    sh "docker tag $imageDb:$version dlish27/$imageDb:latest"
                    sh "docker tag $imageDb:$version dlish27/$imageDb:$version"

                    sh "docker push dlish27/$image:$version"
                    sh "docker push dlish27/$image:latest"
                    sh "docker push dlish27/$imageDb:$version"
                    sh "docker push dlish27/$imageDb:latest"

                    sh "docker rmi $image:$version"
                    sh "docker rmi $imageDb:$version"
                }

                stage('deploy to production') {
                    echo "deploy service to production"
                }
            }
        }
    }

    currentBuild.result = "SUCCESS"

} catch (Exception e) {
    error "Failed: ${e}"
    currentBuild.result = "FAILED"
} finally {
    // notifySlack(currentBuild.result)
}

def publishStagedInfo(String ip) {
    notifyGithub("${env.JOB_NAME}, build [#${env.BUILD_NUMBER}](${env.BUILD_URL}) - Staged deployment can be viewed at: [https://$ip](https://$ip). Staged builds require UAT, click on Jenkins link when finished with UAT to mark the build as 'pass' or 'failed'")
    slackSend(color: 'good',
            message: "${env.JOB_NAME}, build #${env.BUILD_NUMBER} ${env.BUILD_URL} - Staged deployment can be viewed at: https://$ip. Staged builds require UAT, click on Jenkins link when finished with testing to mark the build as 'pass' or 'failed'")
}

def notifyGithub(String comment) {
    def pr  = env.BRANCH_NAME.split("-")[1].trim()
    def pat = readFile('/root/.pat').trim()
    sh "curl -H \"Content-Type: application/json\" -u ifg-bot:$pat -X POST -d '{\"body\": \"$comment\"}' https://api.github.com/repos/ifishgroup/ifg-proshop/issues/$pr/comments"
}

def convertBranchName(String name) {
    return name.replaceAll('/', '_')
}

def notifySlack(String buildStatus) {
    if (env.BRANCH_NAME =~ /(?i)^pr-/ || env.BRANCH_NAME == "master") {
        echo "currentBuild.result=$buildStatus"

        if (buildStatus == null || buildStatus == "") {
            buildStatus = 'FAILED'
        }

        def subject = "${buildStatus}: Job '${env.JOB_NAME}, build #${env.BUILD_NUMBER}'"
        def summary = "${subject} (${env.BUILD_URL})"

        if (buildStatus == 'STARTED') {
            color = 'warning'
        } else if (buildStatus == 'SUCCESS') {
            color = 'good'
        } else {
            color = 'danger'
        }

        slackSend(color: color, message: summary)
    }
}
