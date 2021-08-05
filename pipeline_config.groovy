libraries {
    appimage {
        source = "gombokey"
        destination = 'gombokey-${VERSION}.appimage'
    }
    git {
        repo_url = "https://github.com/livelace/gombokey.git"
    }
    go
    k8s_build {
        image = "harbor-core.k8s-2.livelace.ru/dev/gobuild:latest"
        privileged = true
    }
    mattermost
    nexus {
        source = 'gombokey-${VERSION}.appimage'
        destination = 'dists-internal/gombokey/gombokey-${VERSION}.appimage'
    }
    version
}
