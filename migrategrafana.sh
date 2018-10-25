#!/bin/bash

filereplace_github() {
	# Replace go path in file
	sed -i 's/github.com\/grafana\/grafana\//github.com\/miprotek\/grafana-de\//g' ${1}
}

filereplace_cli() {
	# Replace grafana-cli with grafana-de-cli
	sed -i 's/github.com\/miprotek\/grafana-de\/pkg\/cmd\/grafana-cli/github.com\/miprotek\/grafana-de\/pkg\/cmd\/grafana-de-cli/g' ${1}
}

foldersearch() {
	# Search files in folders recursiv
	echo migrate folder ${1}
	for file in ${1}/*; do
		if [ -d ${file} ]; then
			foldersearch ${file} ${2}
		else
			case "${2}" in
				github) filereplace_github ${file} ;;
				cli)    filereplace_cli ${file} ;;
				*) echo "type not supportet"
				   exit 1
				   ;;
			esac
		fi
	done
}

# Replace github path
foldersearch pkg github

# Replace binary names
foldersearch pkg/cmd/grafana-de-cli cli
