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
                        sh "oc login --token=$OPENSHIFT_TOKEN --server=${SERVER}"

                        // Switch to the project
                        sh "oc project ${OPENSHIFT_PROJECT}"

                        // Build image using OpenShift's BuildConfig
                        sh """
                        oc new-build ${GIT_REPO} --name=${APP_NAME} --branch=${GIT_BRANCH} --strategy=docker --context-dir=docker/app
                        oc start-build ${APP_NAME} --from-dir=. --follow
                        """

                        // Tag and push the image to OpenShift registry
                        sh """
                        oc tag ${OPENSHIFT_PROJECT}/${APP_NAME}:latest ${OPENSHIFT_PROJECT}/${APP_NAME}:stable
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
                        sh "oc login --token=$OPENSHIFT_TOKEN --server=${SERVER}"

                        // Delete older image tags (e.g., all but the latest and stable)
                        sh """
                        IMAGES=$(oc get istag -n ${OPENSHIFT_PROJECT} | grep ${APP_NAME} | awk '{print $1}' | head -n -2)
                        for IMG in $IMAGES; do oc delete istag $IMG; done
                        """
                    }
                }
            }
        }
    }
}
