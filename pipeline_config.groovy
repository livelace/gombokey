def APP_REPO = "https://github.com/livelace/gombokey.git"
def APP_NAME = "gombokey"
def APP_VERSION = env.VERSION + '-${GIT_COMMIT_SHORT}'

libraries {
    appimage {
        source = "${APP_NAME}"
        destination = "gombokey-${APP_VERSION}.appimage"
    }
    dependency_check
    dependency_track {
        project = "${APP_NAME}"
        version = env.VERSION
    }
    git {
        repo_url = "${APP_REPO}"
        repo_branch = env.VERSION
    }
    go {
        options = "github.com/livelace/gombokey/cmd/gombokey"
    }
    k8s_build {
        image = "harbor-core.k8s-2.livelace.ru/dev/gobuild:latest"
        privileged = true
    }
    mattermost
    nexus {
        source = "${APP_NAME}-${APP_VERSION}.appimage"
        destination = "dists-internal/gombokey/${APP_NAME}-${APP_VERSION}.appimage"
    }
    sonarqube
}
