libraries {
    appimage {
        source = "gombokey"
        destination = "gombokey-${env.VERSION}.appimage"
    }
    dependency_check
    dependency_track {
        project = "gombokey"
        version = "${env.VERSION}"
    }
    git {
        repo_url = "https://github.com/livelace/gombokey.git"
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
        destination = 'dists-internal/gombokey/gombokey-${env.VERSION}.appimage'
    }
    sonarqube
    version
}
