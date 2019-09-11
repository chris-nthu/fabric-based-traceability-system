push_to_github() {
    git add -A ;
    git commit -m "Backup" ;
    git push ;
}

collect_files() {
    sudo rm -r fabric/ ;
    sudo rm -r fabric-samples/ ;
    sudo rm -r agriculture_webapp/ ;

    sudo cp -r $HOME/go/src/github.com/hyperledger/fabric-samples $HOME/Github/fabric-based-traceability-system/ ;
    sudo cp -r $HOME/go/src/github.com/hyperledger/fabric $HOME/Github/fabric-based-traceability-system/ ;
    sudo cp -r $HOME/go/src/github.com/hyperledger/agriculture_webapp $HOME/Github/fabric-based-traceability-system/ ;

    sudo rm -rf fabric/.git* ;
    sudo rm -rf fabric-samples/.git* ;
}

upload_to_google_drive() {
    ls -al
}

help_options() {
    printf "%s\n" "Allowed options:
    [-g]       Compress the traceability project and upload to google drive.
    [-c]       Collect the files I want to backup to this folder.
    [-p]       Use git to push this project to github.
    "
}

main() {
    case $1 in
        -g) upload_to_google_drive ;;
        -c) collect_files ;;
        -p) push_to_github ;;
        *) help_options ;;
    esac
}

main "$@" ;
