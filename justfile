# vim: set tabstop=4 shiftwidth=4 expandtab:

[private]
default:
    just --list

[private]
fmt:
    just --fmt --check --unstable

[private]
dagger:
    @dagger version > /dev/null \
    || echo "{{ _REDB }}{{ _WHITE }} Please install dagger CLI {{ _RESET }} https://docs.dagger.io/install/"

[private]
flyctl:
    @flyctl version > /dev/null \
    || echo "{{ _REDB }}{{ _WHITE }} Please install flyctl CLI {{ _RESET }} https://fly.io/docs/flyctl/install/"
    @flyctl auth whoami > /dev/null

_DAGRR_APP_NAME := "testing-dagrr-" + datetime("%Y-%m-%d")

# Run through daggr commands and make sure they work as expected
dagrr: dagger flyctl
    @echo "{{ _MAGENTA }}\n🗡️ Get auto-generated app name...{{ _RESET }}"
    dagger call -m dagrr \
        get-app

    @echo "{{ _MAGENTA }}\n🗡️ Get custom app name...{{ _RESET }}"
    dagger call -m dagrr --app={{ _DAGRR_APP_NAME }} \
        get-app

    @echo "{{ _MAGENTA }}\n🗡️ Generate Dagger Engine on Fly.io manifest...{{ _RESET }}"
    FLY_API_TOKEN=$(flyctl auth token) dagger call -m dagrr --app={{ _DAGRR_APP_NAME }} \
        on-flyio --token=env:FLY_API_TOKEN \
            manifest \
                file --path=fly.toml \
                    export --path=fly.toml
    cat fly.toml

    @echo "{{ _MAGENTA }}\n🗡️ Deploy Dagger Engine on Fly.io with generated manifest...{{ _RESET }}"
    FLY_API_TOKEN=$(flyctl auth token) dagger call -m dagrr --app={{ _DAGRR_APP_NAME }} \
        on-flyio --token=env:FLY_API_TOKEN \
            deploy --dir=.

    @echo "{{ _MAGENTA }}\n🗡️ Delete Dagger Engine on Fly.io...{{ _RESET }}"
    FLY_API_TOKEN=$(flyctl auth token) dagger call -m dagrr --app={{ _DAGRR_APP_NAME }} \
        on-flyio --token=env:FLY_API_TOKEN \
            destroy
    rm fly.toml

    @echo "{{ _MAGENTA }}\n🗡️ Generate Dagger Engine on Fly.io manifest with GPU support and all the other customizations...{{ _RESET }}"
    FLY_API_TOKEN=$(flyctl auth token) dagger call -m dagrr --version=0.13.7 --app={{ _DAGRR_APP_NAME }}-gpu \
        on-flyio --token=env:FLY_API_TOKEN \
            manifest --primary-region=ord --size=performance-1x --memory=4GB --gpu-kind=a10 --disk=50GB --environment="FOO = 'foo'" --environment="BAR = 'bar'" \
                file --path=fly.toml \
                    export --path=fly.toml
    cat fly.toml

    @echo "{{ _MAGENTA }}\n🗡️ Deploy Dagger Engine on Fly.io with new manifest...{{ _RESET }}"
    FLY_API_TOKEN=$(flyctl auth token) dagger call -m dagrr --app={{ _DAGRR_APP_NAME }}-gpu \
        on-flyio --token=env:FLY_API_TOKEN --org=dagger \
            deploy --dir=.

    @echo "{{ _MAGENTA }}\n🗡️ Delete Dagger Engine on Fly.io...{{ _RESET }}"
    FLY_API_TOKEN=$(flyctl auth token) dagger call -m dagrr --app={{ _DAGRR_APP_NAME }}-gpu \
        on-flyio --token=env:FLY_API_TOKEN --org=dagger \
            destroy
    rm fly.toml

# https://linux.101hacks.com/ps1-examples/prompt-color-using-tput/

_BOLD := "$(tput bold)"
_RESET := "$(tput sgr0)"
_BLACK := "$(tput bold)$(tput setaf 0)"
_RED := "$(tput bold)$(tput setaf 1)"
_GREEN := "$(tput bold)$(tput setaf 2)"
_YELLOW := "$(tput bold)$(tput setaf 3)"
_BLUE := "$(tput bold)$(tput setaf 4)"
_MAGENTA := "$(tput bold)$(tput setaf 5)"
_CYAN := "$(tput bold)$(tput setaf 6)"
_WHITE := "$(tput bold)$(tput setaf 7)"
_BLACKB := "$(tput bold)$(tput setab 0)"
_REDB := "$(tput setab 1)$(tput setaf 0)"
_GREENB := "$(tput setab 2)$(tput setaf 0)"
_YELLOWB := "$(tput setab 3)$(tput setaf 0)"
_BLUEB := "$(tput setab 4)$(tput setaf 0)"
_MAGENTAB := "$(tput setab 5)$(tput setaf 0)"
_CYANB := "$(tput setab 6)$(tput setaf 0)"
_WHITEB := "$(tput setab 7)$(tput setaf 0)"
