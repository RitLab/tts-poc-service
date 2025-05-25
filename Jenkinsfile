pipeline {
    agent any

    environment {
        OPENSHIFT_PROJECT = 'ritlab'   // Change to your OpenShift namespace
        APP_NAME = 'tts'                // Change to your app name
        REGISTRY = 'image-registry.openshift-image-registry.svc:5000'
        GIT_REPO = 'https://github.com/RitLab/tts-poc-service'  // Change to your GitHub repo
        GIT_BRANCH = 'feature/old-cpu'
        SERVER = 'https://api.crc.testing:6443'
    }

    stages {
        stage('Checkout Code') {
            steps {
                git branch: "${GIT_BRANCH}", url: "${GIT_REPO}"
            }
        }

        stage('Build & Push Image') {
            steps {
                script {
                    withCredentials([string(credentialsId: "openshift-token", variable: 'OPENSHIFT_TOKEN')]) {
                        // Log in to OpenShift
                        sh "oc login --token=$OPENSHIFT_TOKEN --server=${SERVER} --insecure-skip-tls-verify"

                        // Switch to the project
                        sh "oc project ${OPENSHIFT_PROJECT}"

                        // Build image using OpenShift's BuildConfig
                        sh """
                        set -e
                        oc start-build ${APP_NAME} --from-dir=. --follow
                        """

                        // Tag and push the image to OpenShift registry
                        sh """
                        oc tag ${OPENSHIFT_PROJECT}/${APP_NAME}:dev ${OPENSHIFT_PROJECT}/${APP_NAME}:dev
                        """
                    }
                }
            }
        }

        stage('Cleanup Old Images') {
            steps {
                script {
                    withCredentials([string(credentialsId: "openshift-token", variable: 'OPENSHIFT_TOKEN')]) {
                        // Log in to OpenShift
                        sh "oc login --token=$OPENSHIFT_TOKEN --server=${SERVER} --insecure-skip-tls-verify"

                        // Prune unused images
                        sh """
                        oc adm prune images --confirm --keep-tag-revisions=2
                        """
                    }
                }
            }
        }
    }
}
