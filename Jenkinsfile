pipeline {
    agent any
    
    environment {
        SONAR_TOKEN = credentials('sonar-token')
        DOCKER_REGISTRY = 'your-docker-registry'
        IMAGE_NAME = 'ai-assistant'
        GO_VERSION = '1.24'
    }
    
    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }
        
        stage('Setup Go') {
            steps {
                sh '''
                    go version
                    go mod download
                    go mod tidy
                '''
            }
        }
        
        stage('Lint & Build Check') {
            steps {
                sh '''
                    go fmt ./...
                    go vet ./...
                    go build ./...
                '''
            }
        }
        
        stage('Test') {
            steps {
                sh '''
                    go test -v -coverprofile=coverage.out ./...
                    go tool cover -html=coverage.out -o coverage.html
                '''
            }
            post {
                always {
                    publishHTML([
                        allowMissing: false,
                        alwaysLinkToLastBuild: false,
                        keepAll: true,
                        reportDir: '.',
                        reportFiles: 'coverage.html',
                        reportName: 'Coverage Report'
                    ])
                }
            }
        }
        
        stage('SonarQube Analysis') {
            steps {
                withSonarQubeEnv('SonarQube') {
                    sh '''
                        sonar-scanner \
                        -Dsonar.projectKey=ai-assistant \
                        -Dsonar.sources=. \
                        -Dsonar.exclusions=**/*_test.go,**/vendor/**,**/testdata/**,**/*.pb.go \
                        -Dsonar.tests=. \
                        -Dsonar.test.inclusions=**/*_test.go \
                        -Dsonar.go.coverage.reportPaths=coverage.out
                    '''
                }
            }
        }
        
        stage('Quality Gate') {
            steps {
                timeout(time: 1, unit: 'HOURS') {
                    waitForQualityGate abortPipeline: true
                }
            }
        }
        
        stage('Security Scan - Trivy') {
            steps {
                script {
                    sh '''
                        trivy fs --format json --output trivy-fs-report.json .
                        trivy config --format json --output trivy-config-report.json .
                    '''
                }
            }
            post {
                always {
                    archiveArtifacts artifacts: 'trivy-*.json', fingerprint: true
                }
            }
        }
        
        stage('Build Application') {
            steps {
                sh '''
                    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/main ./cmd/api
                '''
            }
        }
        
        stage('Build Docker Image') {
            steps {
                script {
                    def image = docker.build("${IMAGE_NAME}:${BUILD_NUMBER}", "-f docker/Dockerfile .")
                    sh "trivy image --format json --output trivy-image-report.json ${IMAGE_NAME}:${BUILD_NUMBER}"
                }
            }
            post {
                always {
                    archiveArtifacts artifacts: 'trivy-image-report.json', fingerprint: true
                }
            }
        }
        
        stage('Deploy to Staging') {
            when {
                branch 'develop'
            }
            steps {
                sh '''
                    docker-compose -f docker/docker-compose.yml down
                    docker-compose -f docker/docker-compose.yml up -d
                '''
            }
        }
        
        stage('Deploy to Production') {
            when {
                branch 'main'
            }
            steps {
                input message: 'Deploy to production?', ok: 'Deploy'
                sh '''
                    docker tag ${IMAGE_NAME}:${BUILD_NUMBER} ${DOCKER_REGISTRY}/${IMAGE_NAME}:latest
                    docker push ${DOCKER_REGISTRY}/${IMAGE_NAME}:latest
                '''
            }
        }
    }
    
    post {
        always {
            cleanWs()
        }
        failure {
            emailext (
                subject: "Build Failed: ${env.JOB_NAME} - ${env.BUILD_NUMBER}",
                body: "Build failed. Please check the console output.",
                to: "${env.CHANGE_AUTHOR_EMAIL}"
            )
        }
        success {
            emailext (
                subject: "Build Success: ${env.JOB_NAME} - ${env.BUILD_NUMBER}",
                body: "Build completed successfully.",
                to: "${env.CHANGE_AUTHOR_EMAIL}"
            )
        }
    }
}