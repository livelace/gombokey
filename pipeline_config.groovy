def APP_VERSION = "gombokey-${env.VERSION}-${env.GIT_COMMIT_SHORT}"

libraries {
    appimage {
        source = "gombokey"
        destination = "${APP_VERSION}.appimage"
    }
    dependency_check
    dependency_track {
        project = "gombokey"
        version = "${env.VERSION}"
    }
    git {
        repo_url = "https://github.com/livelace/gombokey.git"
        repo_branch = "${env.VERSION}"
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
        source = "gombokey-${env.VERSION}.appimage"
        destination = "dists-internal/gombokey/${APP_VERSION}.appimage"
    }
    sonarqube
}
