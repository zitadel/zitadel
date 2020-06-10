#!/bin/bash
#debugger
set -x

source ./.github/scripts/variables.env

############################
function setup_git {
############################
  echo "###############"
  echo "set git config"
  echo "###############"

  git config --global user.email "$GIT_USER_MAIL"
  git config --global user.name "$GIT_USER_NAME"
}

############################
function checkout_project {
############################
  echo "###############"
  echo "clone repository $GIT_URL"
  echo "###############"

  # clone opsrepo
  git clone $GIT_URL $LOCAL_TMP_DIR/$GIT_OPSREPO
}

############################
function change_image_version {
############################
  echo "###############"
  echo "checkout master"
  echo "###############"

  cd $LOCAL_TMP_DIR/$GIT_OPSREPO/$GIT_OPSREPO_APPFOLDER/$GIT_OPSREPO_APPLICATION_NAME/overlay/$TARGET_ENVIRONMENT
  git checkout master
  git pull
  echo "###############"
  echo "change image version and commit"
  echo "###############"
  sed -i "s#image: $REGISTRY_IMAGE:.*#image: $REGISTRY_IMAGE:$CAOS_NEXT_VERSION#g" $GIT_OPSREPO_IMAGEFILE
  git add $GIT_OPSREPO_IMAGEFILE
  git commit --message "Github Workflow: $GITHUB_WORKFLOW"
}

############################
function upload_files {
############################
  echo "###############"
  echo "git push"
  echo "###############"
  git push --quiet --set-upstream origin
}
