pipeline {
   agent any
   
   environment {
      CGO_ENABLED = 0 
      GOPATH = "${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}"
      DOCKER_HOST = 'unix:///var/run/docker.sock'
   }

   stages {
      stage('Verify Branch') {
         steps {
            echo "$GIT_BRANCH"
         }
      }
      stage('Docker Build') {
         steps {
            sh(script: 'docker compose build')
         }
      }
      stage('Start App') {
         steps {
            sh(script: 'docker compose up -d')
         }
      }
      stage("Unit Test") {
            agent {
                docker {
                    image 'golang:1.23'
                    args '-v /go/pkg/mod:/go/pkg/mod'
                }
            }
            steps {
                echo "UNIT TEST EXECUTION STARTED in $WORKSPACE/app"
                 sh 'go test ./... -v'
            }
      }
      stage('Docker Push') {
         steps {
            echo "Runnning in $WORKSPACE"
            dir("$WORKSPACE/jenkins") {
               script {
                  docker.withRegistry('', 'dockerhub') {
                     def image = docker.build('gators/jenkins-rundoo:2024')
                     image.push()
                  }
               }
            }
         }
      }
   }
   post {
      always {
         sh(script: 'docker compose down')
      }
   }
}