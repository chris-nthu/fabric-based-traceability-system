push_to_github() {
    git add -A ;
    git commit -m $2 ;
    git push ;
}

collect_files() {
    sudo cp -r $HOME/go/src/github.com/hyperledger/fabric-samples $HOME/Github/fabric-based-traceability-system/ ;
    sudo cp -r $HOME/go/src/github.com/hyperledger/fabric $HOME/Github/fabric-based-traceability-system/ ;
    sudo cp -r $HOME/go/src/github.com/hyperledger/agriculture_webapp $HOME/Github/fabric-based-traceability-system/ ;

    sudo rm -rf fabric/.git* ;
    sudo rm -rf fabric-samples/.git* ;
}

help_options() {
    printf "%s\n" "Allowed options:
    [-c]       Collect the files I want to backup to this folder.
    [-p]       Use git to push this project to github.
    "
}

main() {
    case $1 in
        -c) collect_files ;;
        -p) push_to_github ;;
        *) help_options ;;
    esac
}

main "$@" ;
