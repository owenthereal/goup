SHELL=/bin/bash -o pipefail


.PHONY: vet
vet:
	shellcheck -s dash -- install.sh
